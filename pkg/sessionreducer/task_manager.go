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

// sessionReduceTask represents a session reduce task for a session reduce operation.
type sessionReduceTask struct {
	partition      *v1.Partition
	sessionReducer SessionReducer
	inputCh        chan Datum
	doneCh         chan struct{}
}

// buildSessionReduceResponse builds the session reduce response from the messages.
func (rt *sessionReduceTask) buildSessionReduceResponse(messages Messages) *v1.SessionReduceResponse {
	results := make([]*v1.SessionReduceResponse_Result, 0, len(messages))
	for _, message := range messages {
		results = append(results, &v1.SessionReduceResponse_Result{
			Keys:  message.Keys(),
			Value: message.Value(),
			Tags:  message.Tags(),
		})
	}

	response := &v1.SessionReduceResponse{
		Results:   results,
		Partition: rt.partition,
		EventTime: timestamppb.New(rt.partition.GetEnd().AsTime().Add(-1 * time.Millisecond)),
	}

	return response
}

// uniqueKey returns the unique key for the reduce task to be used in the task manager to identify the task.
func (rt *sessionReduceTask) uniqueKey() string {
	return fmt.Sprintf("%d:%d:%s",
		rt.partition.GetStart().AsTime().UnixMilli(),
		rt.partition.GetEnd().AsTime().UnixMilli(),
		strings.Join(rt.partition.GetKeys(), delimiter))
}

// sessionReduceTaskManager manages the reduce tasks for a session reduce operation.
type sessionReduceTaskManager struct {
	sessionReducerFactory CreateSessionReducer
	tasks                 map[string]*sessionReduceTask
	outputCh              chan *v1.SessionReduceResponse
	rw                    sync.RWMutex
}

func newReduceTaskManager(sessionReducerFactory CreateSessionReducer) *sessionReduceTaskManager {
	return &sessionReduceTaskManager{
		sessionReducerFactory: sessionReducerFactory,
		tasks:                 make(map[string]*sessionReduceTask),
		outputCh:              make(chan *v1.SessionReduceResponse),
	}
}

// CreateTask creates a new reduce task and starts the session reduce operation.
func (rtm *sessionReduceTaskManager) CreateTask(ctx context.Context, request *v1.SessionReduceRequest) error {
	println("create task was invoked with partitions - ", len(request.Operation.Partitions))
	for _, partition := range request.Operation.Partitions {
		println("partition -", partition.GetStart().AsTime().UnixMilli(), " ", partition.GetEnd().AsTime().UnixMilli())
	}

	rtm.rw.Lock()

	// for create operation, there should be exactly one partition
	if len(request.Operation.Partitions) != 1 {
		return fmt.Errorf("invalid number of partitions")
	}

	task := &sessionReduceTask{
		partition:      request.Operation.Partitions[0],
		sessionReducer: rtm.sessionReducerFactory(),
		inputCh:        make(chan Datum),
		doneCh:         make(chan struct{}),
	}

	key := task.uniqueKey()
	rtm.tasks[key] = task

	rtm.rw.Unlock()

	go func() {
		msgs := task.sessionReducer.SessionReduce(ctx, request.GetPayload().GetKeys(), task.inputCh)
		// write the output to the output channel, service will forward it to downstream
		rtm.outputCh <- task.buildSessionReduceResponse(msgs)
		// send a done signal
		close(task.doneCh)
		// delete the task from the tasks list
		rtm.rw.Lock()
		delete(rtm.tasks, key)
		rtm.rw.Unlock()
	}()

	// send the datum to the task if the payload is not nil
	if request.Payload != nil {
		task.inputCh <- buildDatum(request.Payload)
	}

	return nil
}

// AppendToTask writes the message to the reduce task.
// If the reduce task is not found, it will create a new reduce task and start the reduce operation.
func (rtm *sessionReduceTaskManager) AppendToTask(request *v1.SessionReduceRequest) error {
	println("append task was invoked with partitions - ", len(request.Operation.Partitions))
	for _, partition := range request.Operation.Partitions {
		println("partition -", partition.GetStart().AsTime().UnixMilli(), " ", partition.GetEnd().AsTime().UnixMilli())
	}

	// for append operation, there should be exactly one partition
	if len(request.Operation.Partitions) != 1 {
		return fmt.Errorf("invalid number of partitions")
	}

	rtm.rw.RLock()
	task, ok := rtm.tasks[generateKey(request.Operation.Partitions[0])]
	rtm.rw.RUnlock()

	// if the task is not found, create a new task and start the reduce operation
	if !ok {
		return rtm.CreateTask(context.Background(), request)
	}
	// send the datum to the task if the payload is not nil
	if request.Payload != nil {
		task.inputCh <- buildDatum(request.Payload)
	}
	return nil
}

// CloseTask closes the reduce task input channel. The reduce task will be closed when all the messages are processed.
func (rtm *sessionReduceTaskManager) CloseTask(request *v1.SessionReduceRequest) {
	println("close task was invoked with partitions - ", len(request.Operation.Partitions))
	for _, partition := range request.Operation.Partitions {
		println("partition -", partition.GetStart().AsTime().UnixMilli(), " ", partition.GetEnd().AsTime().UnixMilli())
	}

	rtm.rw.RLock()
	tasksToBeClosed := make([]*sessionReduceTask, 0, len(request.Operation.Partitions))
	for _, partition := range request.Operation.Partitions {
		key := generateKey(partition)
		task, ok := rtm.tasks[key]
		if ok {
			tasksToBeClosed = append(tasksToBeClosed, task)
		}
	}
	rtm.rw.RUnlock()

	for _, task := range tasksToBeClosed {
		close(task.inputCh)
	}
}

// MergeTasks merges the session reduce tasks. It will invoke close on all the tasks except the main task. The main task will
// merge the aggregators from the other tasks. The main task will be first task in the list of tasks.
func (rtm *sessionReduceTaskManager) MergeTasks(ctx context.Context, request *v1.SessionReduceRequest) error {
	println("merge task was invoked with partitions - ", len(request.Operation.Partitions))
	for _, partition := range request.Operation.Partitions {
		println("partition -", partition.GetStart().AsTime().UnixMilli(), " ", partition.GetEnd().AsTime().UnixMilli())
	}

	rtm.rw.Lock()
	mergedPartition := request.Operation.Partitions[0]

	tasks := make([]*sessionReduceTask, 0, len(request.Operation.Partitions))
	task, ok := rtm.tasks[generateKey(mergedPartition)]
	if !ok {
		return fmt.Errorf("task not found")
	}

	// send the datum to first task if the payload is not nil
	if request.Payload != nil {
		task.inputCh <- buildDatum(request.Payload)
	}

	// merge the aggregators from the other tasks
	for _, partition := range request.Operation.Partitions[1:] {
		task, ok = rtm.tasks[generateKey(partition)]
		if !ok {
			rtm.rw.Unlock()
			return fmt.Errorf("task not found")
		}
		tasks = append(tasks, task)

		// mergedPartition will be the smallest window which contains all the windows
		if partition.GetStart().AsTime().Before(mergedPartition.GetStart().AsTime()) {
			mergedPartition.Start = partition.Start
		}

		if partition.GetEnd().AsTime().After(mergedPartition.GetEnd().AsTime()) {
			mergedPartition.End = partition.End
		}
	}

	// create a new task with the merged partition
	mergedTask := &sessionReduceTask{
		partition:      mergedPartition,
		sessionReducer: rtm.sessionReducerFactory(),
		inputCh:        make(chan Datum),
		doneCh:         make(chan struct{}),
	}

	// add merged task to the tasks map
	rtm.tasks[mergedTask.uniqueKey()] = mergedTask
	rtm.rw.Unlock()

	accumulators := make([][]byte, 0, len(tasks))
	// close all the tasks and collect the accumulators
	for _, task = range tasks {
		close(task.inputCh)
		accumulators = append(accumulators, task.sessionReducer.Accumulator(ctx))
	}

	// merge the accumulators using the merged task
	for _, aggregator := range accumulators {
		mergedTask.sessionReducer.MergeAccumulator(ctx, aggregator)
	}

	return nil
}

// ExpandTask expands the reduce task. It will expand the window for the reduce task.
// deletes the old task and creates a new task with the new window.
// expects request.Operation.Partitions to have exactly two partitions. The first partition is the old partition and the second
// partition is the new partition.
func (rtm *sessionReduceTaskManager) ExpandTask(request *v1.SessionReduceRequest) error {
	println("expand task was invoked with partitions - ", len(request.Operation.Partitions))
	for _, partition := range request.Operation.Partitions {
		println("partition -", partition.GetStart().AsTime().UnixMilli(), " ", partition.GetEnd().AsTime().UnixMilli())
	}

	// for expand operation, there should be exactly two partitions
	if len(request.Operation.Partitions) != 2 {
		return fmt.Errorf("expected exactly two partitions")
	}

	rtm.rw.Lock()
	key := generateKey(request.Operation.Partitions[0])
	task, ok := rtm.tasks[key]
	if !ok {
		return fmt.Errorf("task not found")
	}

	// assign the new partition to the task
	task.partition = request.Operation.Partitions[1]

	// delete the old entry from the tasks map and add the new entry
	delete(rtm.tasks, key)
	rtm.tasks[task.uniqueKey()] = task
	rtm.rw.Unlock()

	// send the datum to the task if the payload is not nil
	if request.Payload != nil {
		task.inputCh <- buildDatum(request.GetPayload())
	}

	return nil
}

// OutputChannel returns the output channel for the reduce task manager to read the results.
func (rtm *sessionReduceTaskManager) OutputChannel() <-chan *v1.SessionReduceResponse {
	return rtm.outputCh
}

// WaitAll waits for all the pending reduce tasks to complete.
func (rtm *sessionReduceTaskManager) WaitAll() {
	rtm.rw.RLock()
	tasks := make([]*sessionReduceTask, 0, len(rtm.tasks))
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
func (rtm *sessionReduceTaskManager) CloseAll() {
	rtm.rw.Lock()
	tasks := make([]*sessionReduceTask, 0, len(rtm.tasks))
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

func buildDatum(payload *v1.SessionReduceRequest_Payload) Datum {
	return NewHandlerDatum(payload.GetValue(), payload.EventTime.AsTime(), payload.Watermark.AsTime())
}
