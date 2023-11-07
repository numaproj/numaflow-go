package sessionreducer

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/sessionreduce/v1"
)

// reduceTask represents a reduce task for a session reduce operation.
type reduceTask struct {
	combinedKey    string
	partition      *v1.Partition
	reduceStreamer SessionReducer
	watermark      time.Time
	inputCh        chan Datum
	Done           chan struct{}
}

// buildSessionReduceResponse builds the session reduce response from the messages.
func (rt *reduceTask) buildSessionReduceResponse(messages Messages) *v1.SessionReduceResponse {
	result := make([]*v1.SessionReduceResponse_Result, 0, len(messages))
	for _, message := range messages {
		result = append(result, &v1.SessionReduceResponse_Result{
			Keys:  message.Keys(),
			Value: message.Value(),
			Tags:  message.Tags(),
		})
	}

	response := &v1.SessionReduceResponse{
		Result:      result,
		Partition:   rt.partition,
		CombinedKey: rt.combinedKey,
		EventTime:   timestamppb.New(rt.partition.GetEnd().AsTime().Add(-1 * time.Millisecond)),
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

// reduceTaskManager manages the reduce tasks for a session reduce operation.
type reduceTaskManager struct {
	Tasks  map[string]*reduceTask
	Output chan *v1.SessionReduceResponse
	Mutex  sync.RWMutex
}

func newReduceTaskManager() *reduceTaskManager {
	return &reduceTaskManager{
		Tasks:  make(map[string]*reduceTask),
		Output: make(chan *v1.SessionReduceResponse),
	}
}

// CreateTask creates a new reduce task and starts the session reduce operation.
func (rtm *reduceTaskManager) CreateTask(ctx context.Context, request *v1.SessionReduceRequest, reduceStreamer SessionReducer) {
	rtm.Mutex.Lock()
	defer rtm.Mutex.Unlock()

	task := &reduceTask{
		combinedKey:    strings.Join(request.GetPayload().GetKeys(), delimiter),
		partition:      request.Operation.Partitions[0],
		reduceStreamer: reduceStreamer,
		inputCh:        make(chan Datum),
		Done:           make(chan struct{}),
	}

	go func() {
		defer close(task.Done)
		msgs := reduceStreamer.SessionReduce(ctx, request.GetPayload().GetKeys(), task.inputCh)
		rtm.Output <- task.buildSessionReduceResponse(msgs)
	}()

	rtm.Tasks[task.uniqueKey()] = task
}

// AppendToTask writes the message to the reduce task.
func (rtm *reduceTaskManager) AppendToTask(request *v1.SessionReduceRequest) error {
	rtm.Mutex.RLock()
	task, ok := rtm.Tasks[generateKey(request.Operation.Partitions[0], request.Payload.Keys)]
	rtm.Mutex.RUnlock()

	if !ok {
		return fmt.Errorf("task not found")
	}

	task.inputCh <- buildDatum(request)
	return nil
}

// CloseTask closes the reduce task input channel. The reduce task will be closed when all the messages are processed.
func (rtm *reduceTaskManager) CloseTask(request *v1.SessionReduceRequest) {
	rtm.Mutex.Lock()
	tasksToBeClosed := make([]*reduceTask, 0, len(request.Operation.Partitions))
	for _, partition := range request.Operation.Partitions {
		key := generateKey(partition, request.Payload.Keys)
		task, ok := rtm.Tasks[key]
		if ok {
			tasksToBeClosed = append(tasksToBeClosed, task)
		}
		delete(rtm.Tasks, key)
	}
	rtm.Mutex.Unlock()

	for _, task := range tasksToBeClosed {
		close(task.inputCh)
	}
}

// MergeTasks merges the session reduce tasks. It will invoke close on all the tasks except the main task. The main task will
// merge the aggregators from the other tasks. The main task will be first task in the list of tasks.
func (rtm *reduceTaskManager) MergeTasks(ctx context.Context, request *v1.SessionReduceRequest) error {
	rtm.Mutex.Lock()
	tasks := make([]*reduceTask, 0, len(request.Operation.Partitions))
	for _, partition := range request.Operation.Partitions {
		key := generateKey(partition, request.Payload.Keys)
		task, ok := rtm.Tasks[key]
		if !ok {
			rtm.Mutex.Unlock()
			return fmt.Errorf("task not found")
		}
		tasks = append(tasks, task)
	}
	rtm.Mutex.Unlock()

	if len(tasks) == 0 {
		return nil
	}

	mainTask := tasks[0]
	aggregators := make([][]byte, 0, len(tasks)-1)

	for _, task := range tasks[1:] {
		close(task.inputCh)
		aggregators = append(aggregators, task.reduceStreamer.Aggregator(ctx))
	}

	for _, aggregator := range aggregators {
		mainTask.reduceStreamer.MergeAggregator(ctx, aggregator)
	}

	if request.Payload != nil {
		mainTask.inputCh <- buildDatum(request)
	}

	return nil
}

// ExpandTask expands the reduce task. It will expand the window for the reduce task.
// deletes the old task and creates a new task with the new window.
// expects request.Operation.Partitions to have exactly two partitions. The first partition is the old partition and the second
// partition is the new partition.
func (rtm *reduceTaskManager) ExpandTask(request *v1.SessionReduceRequest) error {
	if len(request.Operation.Partitions) != 2 {
		return fmt.Errorf("expected exactly two partitions")
	}
	rtm.Mutex.Lock()
	key := generateKey(request.Operation.Partitions[0], request.Payload.Keys)
	task, ok := rtm.Tasks[key]
	task.partition = request.Operation.Partitions[1]
	delete(rtm.Tasks, key)
	rtm.Tasks[task.uniqueKey()] = task
	rtm.Mutex.Unlock()

	if !ok {
		return fmt.Errorf("task not found")
	}

	return nil
}

// OutputChannel returns the output channel for the reduce task manager to read the results.
func (rtm *reduceTaskManager) OutputChannel() <-chan *v1.SessionReduceResponse {
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

func buildDatum(request *v1.SessionReduceRequest) Datum {
	return NewHandlerDatum(request.Payload.GetValue(), request.Payload.EventTime.AsTime(), request.Payload.Watermark.AsTime())
}
