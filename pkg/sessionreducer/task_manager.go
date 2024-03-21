package sessionreducer

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"go.uber.org/atomic"

	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/sessionreduce/v1"
)

// sessionReduceTask represents a task for a performing session reduce operation.
type sessionReduceTask struct {
	keyedWindow    *v1.KeyedWindow
	sessionReducer SessionReducer
	inputCh        chan Datum
	outputCh       chan Message
	doneCh         chan struct{}
	merged         *atomic.Bool
	lock           sync.RWMutex
}

// buildSessionReduceResponse builds the session reduce response from the messages.
func (rt *sessionReduceTask) buildSessionReduceResponse(message Message) *v1.SessionReduceResponse {
	rt.lock.RLock()
	defer rt.lock.RUnlock()
	response := &v1.SessionReduceResponse{
		Result: &v1.SessionReduceResponse_Result{
			Keys:  message.Keys(),
			Value: message.Value(),
			Tags:  message.Tags(),
		},
		KeyedWindow: rt.keyedWindow,
	}
	return response
}

// buildEOFResponse builds the EOF response for the session reduce task.
func (rt *sessionReduceTask) buildEOFResponse() *v1.SessionReduceResponse {
	rt.lock.RLock()
	defer rt.lock.RUnlock()
	response := &v1.SessionReduceResponse{
		KeyedWindow: rt.keyedWindow,
		EOF:         true,
	}
	return response
}

// assignNewKeyedWindow updates the window for the current session reduce task.
func (rt *sessionReduceTask) assignNewKeyedWindow(kw *v1.KeyedWindow) {
	rt.lock.Lock()
	defer rt.lock.Unlock()
	rt.keyedWindow = kw
}

// getKeyedWindow returns the window for the current session reduce task.
func (rt *sessionReduceTask) getKeyedWindow() *v1.KeyedWindow {
	rt.lock.RLock()
	defer rt.lock.RUnlock()
	return rt.keyedWindow
}

// uniqueKey returns the unique key for the session reduce task to be used in the task manager to identify the task.
func (rt *sessionReduceTask) uniqueKey() string {
	rt.lock.RLock()
	defer rt.lock.RUnlock()
	return fmt.Sprintf("%d:%d:%s",
		rt.keyedWindow.GetStart().AsTime().UnixMilli(),
		rt.keyedWindow.GetEnd().AsTime().UnixMilli(),
		strings.Join(rt.keyedWindow.GetKeys(), delimiter))
}

// sessionReduceTaskManager manages the tasks for a session reduce operation.
type sessionReduceTaskManager struct {
	creatorHandle SessionReducerCreator
	tasks         map[string]*sessionReduceTask
	responseCh    chan *v1.SessionReduceResponse
	rw            sync.RWMutex
}

func newReduceTaskManager(sessionReducerFactory SessionReducerCreator) *sessionReduceTaskManager {
	return &sessionReduceTaskManager{
		creatorHandle: sessionReducerFactory,
		tasks:         make(map[string]*sessionReduceTask),
		responseCh:    make(chan *v1.SessionReduceResponse),
	}
}

// CreateTask creates a new task and starts the session reduce operation.
func (rtm *sessionReduceTaskManager) CreateTask(ctx context.Context, request *v1.SessionReduceRequest) error {
	rtm.rw.Lock()

	// for create operation, there should be exactly one keyedWindow
	if len(request.Operation.KeyedWindows) != 1 {
		return fmt.Errorf("create operation error: invalid number of windows in the request - %d", len(request.Operation.KeyedWindows))
	}

	task := &sessionReduceTask{
		keyedWindow:    request.Operation.KeyedWindows[0],
		sessionReducer: rtm.creatorHandle.Create(),
		inputCh:        make(chan Datum),
		outputCh:       make(chan Message),
		doneCh:         make(chan struct{}),
		merged:         atomic.NewBool(false),
	}

	// add the task to the tasks list
	rtm.tasks[task.uniqueKey()] = task

	rtm.rw.Unlock()

	go func() {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for message := range task.outputCh {
				if !task.merged.Load() {
					// write the output to the output channel, service will forward it to downstream
					// if the task is merged to another task, we don't need to send the response
					rtm.responseCh <- task.buildSessionReduceResponse(message)
				}
			}
			if !task.merged.Load() {
				// send EOF
				rtm.responseCh <- task.buildEOFResponse()
			}
		}()

		task.sessionReducer.SessionReduce(ctx, task.getKeyedWindow().GetKeys(), task.inputCh, task.outputCh)
		// close the output channel and wait for the response to be forwarded
		close(task.outputCh)
		wg.Wait()
		// send a done signal
		close(task.doneCh)
		// delete the task from the tasks list
		rtm.rw.Lock()
		delete(rtm.tasks, task.uniqueKey())
		rtm.rw.Unlock()
	}()

	// send the datum to the task if the payload is not nil
	if request.Payload != nil {
		task.inputCh <- buildDatum(request.Payload)
	}

	return nil
}

// AppendToTask writes the message to the reduce task.
// If the task is not found, it will create a new task and starts the session reduce operation.
func (rtm *sessionReduceTaskManager) AppendToTask(ctx context.Context, request *v1.SessionReduceRequest) error {

	// for append operation, there should be exactly one keyedWindow
	if len(request.Operation.KeyedWindows) != 1 {
		return fmt.Errorf("append operation error: invalid number of windows in the request - %d", len(request.Operation.KeyedWindows))
	}

	rtm.rw.RLock()
	task, ok := rtm.tasks[generateKey(request.Operation.KeyedWindows[0])]
	rtm.rw.RUnlock()

	// if the task is not found, create a new task and start the session reduce operation
	if !ok {
		return rtm.CreateTask(ctx, request)
	}

	// send the datum to the task if the payload is not nil
	if request.Payload != nil {
		task.inputCh <- buildDatum(request.Payload)
	}
	return nil
}

// CloseTask closes the input channel of the session reduce tasks.
func (rtm *sessionReduceTaskManager) CloseTask(request *v1.SessionReduceRequest) {
	rtm.rw.RLock()
	tasksToBeClosed := make([]*sessionReduceTask, 0, len(request.Operation.KeyedWindows))
	for _, window := range request.Operation.KeyedWindows {
		key := generateKey(window)
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

// MergeTasks merges the session reduce tasks. It will create a new task with the merged window and
// merges the accumulators from the other tasks to the merged task.
func (rtm *sessionReduceTaskManager) MergeTasks(ctx context.Context, request *v1.SessionReduceRequest) error {
	rtm.rw.Lock()
	mergedWindow := request.Operation.KeyedWindows[0]

	tasks := make([]*sessionReduceTask, 0, len(request.Operation.KeyedWindows))

	// merge the aggregators from the other tasks
	for _, window := range request.Operation.KeyedWindows {
		key := generateKey(window)
		task, ok := rtm.tasks[key]
		if !ok {
			rtm.rw.Unlock()
			return fmt.Errorf("merge operation error: task not found for %s", key)
		}
		task.merged.Store(true)
		tasks = append(tasks, task)

		// mergedWindow will be the largest window which contains all the windows
		if window.GetStart().AsTime().Before(mergedWindow.GetStart().AsTime()) {
			mergedWindow.Start = window.Start
		}

		if window.GetEnd().AsTime().After(mergedWindow.GetEnd().AsTime()) {
			mergedWindow.End = window.End
		}
	}

	rtm.rw.Unlock()

	accumulators := make([][]byte, 0, len(tasks))
	// close all the tasks and collect the accumulators
	for _, task := range tasks {
		close(task.inputCh)
		// wait for the task to complete
		<-task.doneCh
		accumulators = append(accumulators, task.sessionReducer.Accumulator(ctx))
	}

	// create a new task with the merged keyedWindow
	err := rtm.CreateTask(ctx, &v1.SessionReduceRequest{
		Payload: nil,
		Operation: &v1.SessionReduceRequest_WindowOperation{
			Event:        v1.SessionReduceRequest_WindowOperation_OPEN,
			KeyedWindows: []*v1.KeyedWindow{mergedWindow},
		},
	})
	if err != nil {
		return err
	}

	rtm.rw.RLock()
	mergedTask, ok := rtm.tasks[generateKey(mergedWindow)]
	rtm.rw.RUnlock()
	if !ok {
		return fmt.Errorf("merge operation error: merged task not found for key %s", mergedWindow.String())
	}
	// merge the accumulators using the merged task
	for _, aggregator := range accumulators {
		mergedTask.sessionReducer.MergeAccumulator(ctx, aggregator)
	}

	return nil
}

// ExpandTask expands session reduce task. It will update the keyedWindow of the task
// expects request.Operation.KeyedWindows to have exactly two windows. The first is the old window and the second
// is the new expanded window.
func (rtm *sessionReduceTaskManager) ExpandTask(request *v1.SessionReduceRequest) error {
	// for expand operation, there should be exactly two windows
	if len(request.Operation.KeyedWindows) != 2 {
		return fmt.Errorf("expand operation error: expected exactly two windows")
	}

	rtm.rw.Lock()
	key := generateKey(request.Operation.KeyedWindows[0])
	task, ok := rtm.tasks[key]
	if !ok {
		rtm.rw.Unlock()
		return fmt.Errorf("expand operation error: task not found for key - %s", key)
	}

	// assign the new keyedWindow to the task
	task.assignNewKeyedWindow(request.Operation.KeyedWindows[1])

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

// OutputChannel returns the output channel of the task manager to read the results.
func (rtm *sessionReduceTaskManager) OutputChannel() <-chan *v1.SessionReduceResponse {
	return rtm.responseCh
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
	close(rtm.responseCh)
}

func generateKey(keyedWindow *v1.KeyedWindow) string {
	return fmt.Sprintf("%d:%d:%s",
		keyedWindow.GetStart().AsTime().UnixMilli(),
		keyedWindow.GetEnd().AsTime().UnixMilli(),
		strings.Join(keyedWindow.GetKeys(), delimiter))
}

func buildDatum(payload *v1.SessionReduceRequest_Payload) Datum {
	return NewHandlerDatum(
		payload.GetValue(),
		payload.GetEventTime().AsTime(),
		payload.GetWatermark().AsTime(),
		payload.GetHeaders(),
	)
}
