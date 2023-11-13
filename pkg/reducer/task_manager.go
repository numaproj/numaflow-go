package reducer

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/reduce/v1"
)

// reduceTask represents a reduce task for a  reduce operation.
type reduceTask struct {
	combinedKey   string
	partition     *v1.Partition
	reduceHandler Reducer
	inputCh       chan Datum
	Done          chan struct{}
}

// buildReduceResponse builds the reduce response from the messages.
func (rt *reduceTask) buildReduceResponse(messages Messages) *v1.ReduceResponse {
	results := make([]*v1.ReduceResponse_Result, 0, len(messages))
	for _, message := range messages {
		results = append(results, &v1.ReduceResponse_Result{
			Keys:  message.Keys(),
			Value: message.Value(),
			Tags:  message.Tags(),
		})
	}

	response := &v1.ReduceResponse{
		Results:   results,
		Partition: rt.partition,
		EventTime: timestamppb.New(rt.partition.GetEnd().AsTime().Add(-1 * time.Millisecond)),
	}

	return response
}

// uniqueKey returns the unique key for the reduce task to be used in the task manager to identify the task.
func (rt *reduceTask) uniqueKey() string {
	return fmt.Sprintf("%d:%d:%s",
		rt.partition.GetStart().AsTime().UnixMilli(),
		rt.partition.GetEnd().AsTime().UnixMilli(),
		rt.combinedKey)
}

// reduceTaskManager manages the reduce tasks for a  reduce operation.
type reduceTaskManager struct {
	Tasks  map[string]*reduceTask
	Output chan *v1.ReduceResponse
	Mutex  sync.RWMutex
}

func newReduceTaskManager() *reduceTaskManager {
	return &reduceTaskManager{
		Tasks:  make(map[string]*reduceTask),
		Output: make(chan *v1.ReduceResponse),
	}
}

// CreateTask creates a new reduce task and starts the  reduce operation.
func (rtm *reduceTaskManager) CreateTask(ctx context.Context, request *v1.ReduceRequest, reducer Reducer) error {
	rtm.Mutex.Lock()
	if len(request.Operation.Partitions) != 1 {
		return fmt.Errorf("invalid number of partitions")
	}

	md := NewMetadata(NewIntervalWindow(request.Operation.Partitions[0].GetStart().AsTime(),
		request.Operation.Partitions[0].GetEnd().AsTime()))

	task := &reduceTask{
		combinedKey: strings.Join(request.GetPayload().GetKeys(), delimiter),
		partition:   request.Operation.Partitions[0],
		inputCh:     make(chan Datum),
		Done:        make(chan struct{}),
	}

	key := task.uniqueKey()
	rtm.Tasks[key] = task

	rtm.Mutex.Unlock()

	go func() {
		msgs := reducer.Reduce(ctx, request.GetPayload().GetKeys(), task.inputCh, md)
		// write the output to the output channel, service will forward it to downstream
		rtm.Output <- task.buildReduceResponse(msgs)
		// send a done signal
		close(task.Done)
	}()

	// write the first message to the input channel
	task.inputCh <- buildDatum(request)
	return nil
}

// AppendToTask writes the message to the reduce task.
// If the task is not found, it creates a new task and starts the reduce operation.
func (rtm *reduceTaskManager) AppendToTask(request *v1.ReduceRequest, reducer Reducer) error {
	rtm.Mutex.RLock()
	gKey := generateKey(request.Operation.Partitions[0], request.Payload.Keys)
	task, ok := rtm.Tasks[gKey]
	rtm.Mutex.RUnlock()

	// if the task is not found, create a new task
	if !ok {
		return rtm.CreateTask(context.Background(), request, reducer)
	}

	task.inputCh <- buildDatum(request)
	return nil
}

// OutputChannel returns the output channel for the reduce task manager to read the results.
func (rtm *reduceTaskManager) OutputChannel() <-chan *v1.ReduceResponse {
	return rtm.Output
}

// WaitAll waits for all the reduce tasks to complete.
func (rtm *reduceTaskManager) WaitAll() {
	rtm.Mutex.RLock()
	tasks := make([]*reduceTask, 0, len(rtm.Tasks))
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
func (rtm *reduceTaskManager) CloseAll() {
	rtm.Mutex.Lock()
	tasks := make([]*reduceTask, 0, len(rtm.Tasks))
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

func buildDatum(request *v1.ReduceRequest) Datum {
	return NewHandlerDatum(request.Payload.GetValue(), request.Payload.EventTime.AsTime(), request.Payload.Watermark.AsTime())
}
