package reducer

import (
	"context"
	"fmt"
	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/reduce/v1"
	"log"
	"runtime/debug"
	"strings"
	"sync"
)

// reduceTask represents a task for a performing reduceStream operation.
type reduceTask struct {
	keys    []string
	window  *v1.Window
	reducer Reducer
	inputCh chan Datum
	doneCh  chan struct{}
	// Add an error field to track task errors
	err error
}

// buildReduceResponse builds the reduce response from the messages.
func (rt *reduceTask) buildReduceResponse(message Message) *v1.ReduceResponse {

	response := &v1.ReduceResponse{
		Result: &v1.ReduceResponse_Result{
			Keys:  message.Keys(),
			Value: message.Value(),
			Tags:  message.Tags(),
		},
		Window: rt.window,
	}

	return response
}

func (rt *reduceTask) buildEOFResponse() *v1.ReduceResponse {
	response := &v1.ReduceResponse{
		Window: rt.window,
		EOF:    true,
	}

	return response
}

// uniqueKey returns the unique key for the reduce task to be used in the task manager to identify the task.
func (rt *reduceTask) uniqueKey() string {
	return fmt.Sprintf("%d:%d:%s",
		rt.window.GetStart().AsTime().UnixMilli(),
		rt.window.GetEnd().AsTime().UnixMilli(),
		strings.Join(rt.keys, delimiter))
}

// windowContext holds a window and its associated tasks
type windowContext struct {
	window *v1.Window
	tasks  map[string]*reduceTask
}

// reduceTaskManager manages the reduce tasks for a reduce operation.
type reduceTaskManager struct {
	reducerCreatorHandle ReducerCreator
	// windowContexts maps window keys to their context (window object + tasks)
	windowContexts map[string]*windowContext
	responseCh     chan *v1.ReduceResponse
	shutdownCh     chan<- struct{}
	// Add a wait group specifically for CloseTask goroutines
	closeWg sync.WaitGroup
}

func newReduceTaskManager(reducerCreatorHandle ReducerCreator, shutdownCh chan<- struct{}) *reduceTaskManager {
	return &reduceTaskManager{
		reducerCreatorHandle: reducerCreatorHandle,
		windowContexts:       make(map[string]*windowContext),
		responseCh:           make(chan *v1.ReduceResponse),
		shutdownCh:           shutdownCh,
	}
}

// windowKey generates a key for the window based on start and end times
func windowKey(window *v1.Window) string {
	return fmt.Sprintf("%d:%d",
		window.GetStart().AsTime().UnixMilli(),
		window.GetEnd().AsTime().UnixMilli())
}

// taskKey generates a key for the task based on keys
func taskKey(keys []string) string {
	return strings.Join(keys, delimiter)
}

// CreateTask creates a new reduce task and starts the reduce operation.
func (rtm *reduceTaskManager) CreateTask(ctx context.Context, request *v1.ReduceRequest) error {
	if len(request.Operation.Windows) != 1 {
		return fmt.Errorf("create operation error: invalid number of windows")
	}

	window := request.Operation.Windows[0]
	wKey := windowKey(window)
	tKey := taskKey(request.GetPayload().GetKeys())

	// Initialize the window context if it doesn't exist
	if _, exists := rtm.windowContexts[wKey]; !exists {
		rtm.windowContexts[wKey] = &windowContext{
			window: window,
			tasks:  make(map[string]*reduceTask),
		}
	}

	md := NewMetadata(NewIntervalWindow(window.GetStart().AsTime(),
		window.GetEnd().AsTime()))

	task := &reduceTask{
		keys:    request.GetPayload().GetKeys(),
		window:  window,
		inputCh: make(chan Datum),
		doneCh:  make(chan struct{}),
	}

	rtm.windowContexts[wKey].tasks[tKey] = task

	go func() {
		// Always close doneCh when the goroutine exits, regardless of success or failure
		defer close(task.doneCh)

		// handle panic
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())
				log.Printf("panic inside reduce handler: %v %v", r, stack)
				task.err = fmt.Errorf("panic in reducer: %v", r)
				rtm.shutdownCh <- struct{}{}
			}
		}()

		// invoke the reduce function
		reducerHandle := rtm.reducerCreatorHandle.Create()

		messages := reducerHandle.Reduce(ctx, request.GetPayload().GetKeys(), task.inputCh, md)

		for _, message := range messages {
			// write the output to the output channel, service will forward it to downstream
			rtm.responseCh <- task.buildReduceResponse(message)
		}
	}()

	// write the first message to the input channel
	task.inputCh <- buildDatum(request)
	return nil
}

// AppendToTask writes the message to the reduce task.
// If the task is not found, it creates a new task and starts the reduce operation.
func (rtm *reduceTaskManager) AppendToTask(ctx context.Context, request *v1.ReduceRequest) error {
	if len(request.Operation.Windows) != 1 {
		return fmt.Errorf("append operation error: invalid number of windows")
	}

	window := request.Operation.Windows[0]
	wKey := windowKey(window)
	tKey := taskKey(request.GetPayload().GetKeys())

	// Check if we have a window context for this window
	winCtx, windowExists := rtm.windowContexts[wKey]
	if !windowExists {
		return rtm.CreateTask(ctx, request)
	}

	// Check if we have a task for these keys
	task, taskExists := winCtx.tasks[tKey]
	if !taskExists {
		return rtm.CreateTask(ctx, request)
	}

	task.inputCh <- buildDatum(request)
	return nil
}

// OutputChannel returns the output channel for the reduce task manager to read the results.
func (rtm *reduceTaskManager) OutputChannel() <-chan *v1.ReduceResponse {
	return rtm.responseCh
}

// WaitAll waits for all the reduce tasks to complete.
func (rtm *reduceTaskManager) WaitAll() {
	for _, winCtx := range rtm.windowContexts {
		var hasErrors bool

		for _, task := range winCtx.tasks {
			<-task.doneCh
			if task.err != nil {
				hasErrors = true
				log.Printf("Error in reduce task during WaitAll: %v", task.err)
			}
		}

		if !hasErrors {
			rtm.responseCh <- &v1.ReduceResponse{
				Window: winCtx.window,
				EOF:    true,
			}
		}
	}

	rtm.closeWg.Wait()
	// after all the tasks are completed, close the output channel
	close(rtm.responseCh)
}

// CloseAll closes all the reduce tasks.
func (rtm *reduceTaskManager) CloseAll() {
	for _, winCtx := range rtm.windowContexts {
		for _, task := range winCtx.tasks {
			close(task.inputCh)
		}
	}
}

// CloseTask closes all reduce tasks matching the window's start and end time.
func (rtm *reduceTaskManager) CloseTask(request *v1.ReduceRequest) error {
	if len(request.Operation.Windows) != 1 {
		return fmt.Errorf("close operation error: invalid number of windows")
	}

	window := request.Operation.Windows[0]
	wKey := windowKey(window)
	winCtx, ok := rtm.windowContexts[wKey]
	if !ok || len(winCtx.tasks) == 0 {
		return fmt.Errorf("close operation error: no tasks found for window")
	}

	tasksToWait := make([]*reduceTask, 0, len(winCtx.tasks))

	for _, task := range winCtx.tasks {
		tasksToWait = append(tasksToWait, task)
		close(task.inputCh)
	}

	// Wait for all tasks to complete in a separate goroutine
	rtm.closeWg.Add(1)
	go func(tasks []*reduceTask, win *v1.Window) {
		defer rtm.closeWg.Done()
		// Wait for all tasks to complete
		for _, t := range tasks {
			<-t.doneCh
		}

		// Check if any tasks had errors
		var hasErrors bool
		for _, t := range tasks {
			if t.err != nil {
				hasErrors = true
				log.Printf("Error in reduce task: %v", t.err)
			}
		}

		if !hasErrors {
			// After all tasks are done, send the EOF response
			rtm.responseCh <- &v1.ReduceResponse{
				Window: win,
				EOF:    true,
			}
		}
	}(tasksToWait, winCtx.window)

	// Remove the window context from the map
	delete(rtm.windowContexts, wKey)
	return nil
}

func buildDatum(request *v1.ReduceRequest) Datum {
	return NewHandlerDatum(
		request.GetPayload().GetValue(),
		request.GetPayload().GetEventTime().AsTime(),
		request.GetPayload().GetWatermark().AsTime(),
		request.GetPayload().GetHeaders(),
	)
}
