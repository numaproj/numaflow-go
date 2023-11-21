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
	partition       *v1.Partition
	globalReducer   GlobalReducer
	inputCh         chan Datum
	outputCh        chan Messages
	latestWatermark *atomic.Time
	doneCh          chan struct{}
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
		Results:   results,
		Partition: rt.partition,
		EventTime: timestamppb.New(rt.latestWatermark.Load()),
	}

	return response
}

// uniqueKey returns the unique key for the reduce task to be used in the task manager to identify the task.
func (rt *globalReduceTask) uniqueKey() string {
	return fmt.Sprintf("%d:%d:%s",
		rt.partition.GetStart().AsTime().UnixMilli(),
		rt.partition.GetEnd().AsTime().UnixMilli(),
		strings.Join(rt.partition.GetKeys(), delimiter))
}

// globalReduceTaskManager manages the reduce tasks for a global reduce operation.
type globalReduceTaskManager struct {
	globalReducer GlobalReducer
	tasks         map[string]*globalReduceTask
	outputCh      chan *v1.GlobalReduceResponse
	rw            sync.RWMutex
}

func newReduceTaskManager(globalReducer GlobalReducer) *globalReduceTaskManager {
	return &globalReduceTaskManager{
		tasks:         make(map[string]*globalReduceTask),
		outputCh:      make(chan *v1.GlobalReduceResponse),
		globalReducer: globalReducer,
	}
}

// CreateTask creates a new reduce task and starts the global reduce operation.
func (rtm *globalReduceTaskManager) CreateTask(ctx context.Context, request *v1.GlobalReduceRequest) error {
	rtm.rw.Lock()

	if len(request.Operation.Partitions) != 1 {
		return fmt.Errorf("invalid number of partitions")
	}

	task := &globalReduceTask{
		partition:       request.Operation.Partitions[0],
		globalReducer:   rtm.globalReducer,
		inputCh:         make(chan Datum),
		outputCh:        make(chan Messages),
		doneCh:          make(chan struct{}),
		latestWatermark: atomic.NewTime(time.Time{}),
	}

	key := task.uniqueKey()
	rtm.tasks[key] = task

	rtm.rw.Unlock()

	go func() {
		defer close(task.doneCh)
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
					rtm.outputCh <- task.buildGlobalReduceResponse(message)
				}
			}
			// send a done signal
			close(task.doneCh)
			// delete the task from the tasks list
			rtm.rw.Lock()
			delete(rtm.tasks, key)
			rtm.rw.Unlock()
		}()
		// start the global reduce operation
		task.globalReducer.GlobalReduce(ctx, request.GetPayload().GetKeys(), task.inputCh, task.outputCh)
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
func (rtm *globalReduceTaskManager) AppendToTask(ctx context.Context, request *v1.GlobalReduceRequest) error {
	if len(request.Operation.Partitions) != 1 {
		return fmt.Errorf("invalid number of partitions")
	}

	rtm.rw.RLock()
	task, ok := rtm.tasks[generateKey(request.Operation.Partitions[0])]
	rtm.rw.RUnlock()

	if !ok {
		return rtm.CreateTask(ctx, request)
	}

	// send the datum to the task if the payload is not nil
	if request.Payload != nil {
		task.inputCh <- buildDatum(request)
		task.latestWatermark.Store(request.Payload.Watermark.AsTime())
	}

	return nil
}

// CloseTask closes the reduce task input channel. The reduce task will be closed when all the messages are processed.
func (rtm *globalReduceTaskManager) CloseTask(request *v1.GlobalReduceRequest) {
	rtm.rw.RLock()
	tasksToBeClosed := make([]*globalReduceTask, 0, len(request.Operation.Partitions))
	for _, partition := range request.Operation.Partitions {
		key := generateKey(partition)
		task, ok := rtm.tasks[key]
		if ok {
			tasksToBeClosed = append(tasksToBeClosed, task)
		}
	}
	rtm.rw.RUnlock()

	// close the input channel for all the tasks
	for _, task := range tasksToBeClosed {
		close(task.inputCh)
	}
}

// OutputChannel returns the output channel for the reduce task manager to read the results.
func (rtm *globalReduceTaskManager) OutputChannel() <-chan *v1.GlobalReduceResponse {
	return rtm.outputCh
}

// WaitAll waits for all the reduce tasks to complete.
func (rtm *globalReduceTaskManager) WaitAll() {
	rtm.rw.RLock()
	tasks := make([]*globalReduceTask, 0, len(rtm.tasks))
	for _, task := range rtm.tasks {
		tasks = append(tasks, task)
	}
	rtm.rw.RUnlock()

	for _, task := range tasks {
		<-task.doneCh
	}
	// after all the tasks are completed, close the output channel
	close(rtm.outputCh)
}

// CloseAll closes all the reduce tasks.
func (rtm *globalReduceTaskManager) CloseAll() {
	rtm.rw.Lock()
	tasks := make([]*globalReduceTask, 0, len(rtm.tasks))
	for _, task := range rtm.tasks {
		tasks = append(tasks, task)
	}
	rtm.rw.Unlock()

	for _, task := range tasks {
		close(task.inputCh)
	}
}

func generateKey(partition *v1.Partition) string {
	return fmt.Sprintf("%d:%d:%s",
		partition.GetStart().AsTime().UnixMilli(),
		partition.GetEnd().AsTime().UnixMilli(),
		strings.Join(partition.GetKeys(), delimiter))
}

func buildDatum(request *v1.GlobalReduceRequest) Datum {
	return NewHandlerDatum(request.Payload.GetValue(), request.Payload.EventTime.AsTime(), request.Payload.Watermark.AsTime())
}
