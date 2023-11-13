package globalreducer

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.uber.org/atomic"
	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/globalreduce/v1"
)

// globalReduceTask represents a global reduce task for a global reduce operation.
type globalReduceTask struct {
	combinedKey     string
	partition       *v1.Partition
	reduceStreamer  GlobalReducer
	inputCh         chan Datum
	outputCh        chan Messages
	latestWatermark *atomic.Time
	Done            chan struct{}
}

// buildGlobalReduceResponse builds the global reduce response from the messages.
func (rt *globalReduceTask) buildGlobalReduceResponse(messages Messages) *v1.GlobalReduceResponse {
	results := make([]*v1.GlobalReduceResponse_Result, 0, len(messages))
	for _, message := range messages {
		results = append(results, &v1.GlobalReduceResponse_Result{
			Keys:  message.Keys(),
			Value: message.Value(),
			Tags:  message.Tags(),
		})
	}

	// the event time is the latest watermark
	// since for global window, partition start and end time will be -1
	response := &v1.GlobalReduceResponse{
		Results:     results,
		Partition:   rt.partition,
		CombinedKey: rt.combinedKey,
		EventTime:   timestamppb.New(rt.latestWatermark.Load()),
	}

	return response
}

// uniqueKey returns the unique key for the reduce task to be used in the task manager to identify the task.
func (rt *globalReduceTask) uniqueKey() string {
	return fmt.Sprintf("%d:%d:%s",
		rt.partition.GetStart().AsTime().UnixMilli(),
		rt.partition.GetEnd().AsTime().UnixMilli(),
		rt.combinedKey)
}

// globalReduceTaskManager manages the reduce tasks for a global reduce operation.
type globalReduceTaskManager struct {
	Tasks  map[string]*globalReduceTask
	Output chan *v1.GlobalReduceResponse
	Mutex  sync.RWMutex
}

func newReduceTaskManager() *globalReduceTaskManager {
	return &globalReduceTaskManager{
		Tasks:  make(map[string]*globalReduceTask),
		Output: make(chan *v1.GlobalReduceResponse),
	}
}

// CreateTask creates a new reduce task and starts the global reduce operation.
func (rtm *globalReduceTaskManager) CreateTask(ctx context.Context, request *v1.GlobalReduceRequest, globalReducer GlobalReducer) error {
	rtm.Mutex.Lock()

	if len(request.Operation.Partitions) != 1 {
		return fmt.Errorf("invalid number of partitions")
	}

	task := &globalReduceTask{
		combinedKey:     strings.Join(request.GetPayload().GetKeys(), delimiter),
		partition:       request.Operation.Partitions[0],
		reduceStreamer:  globalReducer,
		inputCh:         make(chan Datum),
		outputCh:        make(chan Messages),
		Done:            make(chan struct{}),
		latestWatermark: atomic.NewTime(time.Time{}),
	}

	key := task.uniqueKey()
	rtm.Tasks[key] = task

	rtm.Mutex.Unlock()

	go func() {
		defer close(task.Done)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
		readLoop:
			for {
				select {
				case <-ctx.Done():
					return
				case message, ok := <-task.outputCh:
					if !ok {
						break readLoop
					}
					rtm.Output <- task.buildGlobalReduceResponse(message)
				}
			}
			// send a done signal
			close(task.Done)
			// delete the task from the tasks list
			rtm.Mutex.Lock()
			delete(rtm.Tasks, key)
			rtm.Mutex.Unlock()
		}()
		// start the global reduce operation
		globalReducer.GlobalReduce(ctx, request.GetPayload().GetKeys(), task.inputCh, task.outputCh)
		// close the output channel after the global reduce operation is completed
		close(task.outputCh)
		// wait for the output channel reader to complete
		wg.Wait()
	}()

	// send the datum to the task if the payload is not nil
	// also update the latest watermark
	if request.Payload != nil {
		task.inputCh <- buildDatum(request)
		task.latestWatermark.Store(request.Payload.Watermark.AsTime())
	}

	return nil
}

// AppendToTask writes the message to the reduce task.
// If the reduce task is not found, it will create a new reduce task and start the reduce operation.
func (rtm *globalReduceTaskManager) AppendToTask(request *v1.GlobalReduceRequest, globalReducer GlobalReducer) error {
	if len(request.Operation.Partitions) != 1 {
		return fmt.Errorf("invalid number of partitions")
	}

	rtm.Mutex.RLock()
	task, ok := rtm.Tasks[generateKey(request.Operation.Partitions[0], request.Payload.Keys)]
	rtm.Mutex.RUnlock()

	if !ok {
		return rtm.CreateTask(context.Background(), request, globalReducer)
	}

	if request.Payload != nil {
		task.inputCh <- buildDatum(request)
		task.latestWatermark.Store(request.Payload.Watermark.AsTime())
	}

	return nil
}

// CloseTask closes the reduce task input channel. The reduce task will be closed when all the messages are processed.
func (rtm *globalReduceTaskManager) CloseTask(request *v1.GlobalReduceRequest) {
	rtm.Mutex.Lock()
	tasksToBeClosed := make([]*globalReduceTask, 0, len(request.Operation.Partitions))
	for _, partition := range request.Operation.Partitions {
		key := generateKey(partition, request.Payload.Keys)
		task, ok := rtm.Tasks[key]
		if ok {
			tasksToBeClosed = append(tasksToBeClosed, task)
		}
	}
	rtm.Mutex.Unlock()

	for _, task := range tasksToBeClosed {
		close(task.inputCh)
	}
}

// OutputChannel returns the output channel for the reduce task manager to read the results.
func (rtm *globalReduceTaskManager) OutputChannel() <-chan *v1.GlobalReduceResponse {
	return rtm.Output
}

// WaitAll waits for all the reduce tasks to complete.
func (rtm *globalReduceTaskManager) WaitAll() {
	rtm.Mutex.RLock()
	tasks := make([]*globalReduceTask, 0, len(rtm.Tasks))
	for _, task := range rtm.Tasks {
		tasks = append(tasks, task)
	}
	rtm.Mutex.RUnlock()

	for _, task := range tasks {
		<-task.Done
	}
	// after all the tasks are completed, close the output channel
	close(rtm.Output)
}

// CloseAll closes all the reduce tasks.
func (rtm *globalReduceTaskManager) CloseAll() {
	rtm.Mutex.Lock()
	tasks := make([]*globalReduceTask, 0, len(rtm.Tasks))
	for _, task := range rtm.Tasks {
		tasks = append(tasks, task)
	}
	rtm.Mutex.Unlock()

	for _, task := range tasks {
		close(task.inputCh)
	}
}

func generateKey(partition *v1.Partition, keys []string) string {
	return fmt.Sprintf("%d:%d:%s",
		partition.GetStart().AsTime().UnixMilli(),
		partition.GetEnd().AsTime().UnixMilli(),
		strings.Join(keys, delimiter))
}

func buildDatum(request *v1.GlobalReduceRequest) Datum {
	return NewHandlerDatum(request.Payload.GetValue(), request.Payload.EventTime.AsTime(), request.Payload.Watermark.AsTime())
}
