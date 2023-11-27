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
	keyedWindow     *v1.KeyedWindow
	globalReducer   GlobalReducer
	inputCh         chan Datum
	outputCh        chan Message
	latestWatermark *atomic.Time
	doneCh          chan struct{}
}

// buildGlobalReduceResponse builds the global reduce response from the messages.
func (rt *globalReduceTask) buildGlobalReduceResponse(message Message) *v1.GlobalReduceResponse {

	// the event time is the latest watermark
	// since for global window, keyedWindow start and end time will be -1
	response := &v1.GlobalReduceResponse{
		Result: &v1.GlobalReduceResponse_Result{
			Keys:      message.Keys(),
			Value:     message.Value(),
			Tags:      message.Tags(),
			EventTime: timestamppb.New(rt.latestWatermark.Load()),
		},
		KeyedWindow: rt.keyedWindow,
	}

	return response
}

func (rt *globalReduceTask) buildEOFResponse() *v1.GlobalReduceResponse {
	response := &v1.GlobalReduceResponse{
		Result: &v1.GlobalReduceResponse_Result{
			EventTime: timestamppb.New(rt.latestWatermark.Load()),
		},
		KeyedWindow: rt.keyedWindow,
		EOF:         true,
	}

	return response
}

// uniqueKey returns the unique key for the reduce task to be used in the task manager to identify the task.
func (rt *globalReduceTask) uniqueKey() string {
	return fmt.Sprintf("%d:%d:%s",
		rt.keyedWindow.GetStart().AsTime().UnixMilli(),
		rt.keyedWindow.GetEnd().AsTime().UnixMilli(),
		strings.Join(rt.keyedWindow.GetKeys(), delimiter))
}

// globalReduceTaskManager manages the reduce tasks for a global reduce operation.
type globalReduceTaskManager struct {
	globalReducer GlobalReducer
	tasks         map[string]*globalReduceTask
	responseCh    chan *v1.GlobalReduceResponse
	rw            sync.RWMutex
}

func newReduceTaskManager(globalReducer GlobalReducer) *globalReduceTaskManager {
	return &globalReduceTaskManager{
		tasks:         make(map[string]*globalReduceTask),
		responseCh:    make(chan *v1.GlobalReduceResponse),
		globalReducer: globalReducer,
	}
}

// CreateTask creates a new reduce task and starts the global reduce operation.
func (rtm *globalReduceTaskManager) CreateTask(ctx context.Context, request *v1.GlobalReduceRequest) error {
	rtm.rw.Lock()

	if len(request.Operation.KeyedWindows) != 1 {
		return fmt.Errorf("create operation error: invalid number of windows")
	}

	task := &globalReduceTask{
		keyedWindow:     request.Operation.KeyedWindows[0],
		globalReducer:   rtm.globalReducer,
		inputCh:         make(chan Datum),
		outputCh:        make(chan Message),
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
			for msg := range task.outputCh {
				rtm.responseCh <- task.buildGlobalReduceResponse(msg)
			}
			// send the EOF response
			rtm.responseCh <- task.buildEOFResponse()
		}()

		// start the global reduce operation
		task.globalReducer.GlobalReduce(ctx, request.GetPayload().GetKeys(), task.inputCh, task.outputCh)
		// close the output channel after the global reduce operation is completed
		close(task.outputCh)
		// wait for the output channel reader to complete
		wg.Wait()
		// send a done signal
		close(task.doneCh)
		// delete the task from the tasks list
		rtm.rw.Lock()
		delete(rtm.tasks, key)
		rtm.rw.Unlock()
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
	if len(request.Operation.KeyedWindows) != 1 {
		return fmt.Errorf("append operation error: invalid number of windows")
	}

	rtm.rw.RLock()
	task, ok := rtm.tasks[generateKey(request.Operation.KeyedWindows[0])]
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
	tasksToBeClosed := make([]*globalReduceTask, 0, len(request.Operation.KeyedWindows))
	for _, window := range request.Operation.KeyedWindows {
		key := generateKey(window)
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
	return rtm.responseCh
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
	close(rtm.responseCh)
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

func generateKey(keyedWindow *v1.KeyedWindow) string {
	return fmt.Sprintf("%d:%d:%s",
		keyedWindow.GetStart().AsTime().UnixMilli(),
		keyedWindow.GetEnd().AsTime().UnixMilli(),
		strings.Join(keyedWindow.GetKeys(), delimiter))
}

func buildDatum(request *v1.GlobalReduceRequest) Datum {
	return NewHandlerDatum(request.Payload.GetValue(), request.Payload.EventTime.AsTime(), request.Payload.Watermark.AsTime())
}
