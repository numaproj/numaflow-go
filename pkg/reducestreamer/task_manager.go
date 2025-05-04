package reducestreamer

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"strings"
	"sync"

	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/reduce/v1"
)

// reduceStreamTask represents a task for a performing reduceStream operation.
type reduceStreamTask struct {
	keys           []string
	window         *v1.Window
	reduceStreamer ReduceStreamer
	inputCh        chan Datum
	outputCh       chan Message
	doneCh         chan struct{}
	// Add an error field to track task errors
	err error
}

// buildReduceResponse builds the reduce response from the messages.
func (rt *reduceStreamTask) buildReduceResponse(message Message) *v1.ReduceResponse {

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

func (rt *reduceStreamTask) buildEOFResponse() *v1.ReduceResponse {
	response := &v1.ReduceResponse{
		Window: rt.window,
		EOF:    true,
	}

	return response
}

// windowContext holds a window and its associated tasks
type windowContext struct {
	window *v1.Window
	tasks  map[string]*reduceStreamTask
}

// reduceStreamTaskManager manages the reduceStream tasks.
type reduceStreamTaskManager struct {
	creatorHandle ReduceStreamerCreator
	// Replace tasks map with windowContexts map
	windowContexts map[string]*windowContext
	responseCh     chan *v1.ReduceResponse
	shutdownCh     chan<- struct{}
	// Add a wait group specifically for CloseTask goroutines
	closeWg sync.WaitGroup
}

func newReduceTaskManager(reduceStreamerCreator ReduceStreamerCreator, shutdownCh chan<- struct{}) *reduceStreamTaskManager {
	return &reduceStreamTaskManager{
		creatorHandle:  reduceStreamerCreator,
		windowContexts: make(map[string]*windowContext),
		responseCh:     make(chan *v1.ReduceResponse),
		shutdownCh:     shutdownCh,
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

// CreateTask creates a new reduceStream task and starts the reduceStream operation.
func (rtm *reduceStreamTaskManager) CreateTask(ctx context.Context, request *v1.ReduceRequest) error {
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
			tasks:  make(map[string]*reduceStreamTask),
		}
	}

	md := NewMetadata(NewIntervalWindow(window.GetStart().AsTime(),
		window.GetEnd().AsTime()))

	task := &reduceStreamTask{
		keys:     request.GetPayload().GetKeys(),
		window:   window,
		inputCh:  make(chan Datum),
		outputCh: make(chan Message),
		doneCh:   make(chan struct{}),
	}

	rtm.windowContexts[wKey].tasks[tKey] = task

	go func() {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for message := range task.outputCh {
				// write the output to the output channel, service will forward it to downstream
				rtm.responseCh <- task.buildReduceResponse(message)
			}
		}()

		// Always close doneCh when the goroutine exits, regardless of success or failure
		defer close(task.doneCh)

		// handle panic
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic inside reduce handler: %v %v", r, string(debug.Stack()))
				task.err = fmt.Errorf("panic in reducer: %v", r)
				rtm.shutdownCh <- struct{}{}
			}
		}()

		reduceStreamerHandle := rtm.creatorHandle.Create()
		// invoke the reduceStream function
		reduceStreamerHandle.ReduceStream(ctx, request.GetPayload().GetKeys(), task.inputCh, task.outputCh, md)
		// close the output channel after the reduceStream function is done
		close(task.outputCh)
		// wait for the output to be forwarded
		wg.Wait()
	}()

	// write the first message to the input channel
	task.inputCh <- buildDatum(request)
	return nil
}

// AppendToTask writes the message to the reduceStream task.
// If the task is not found, it creates a new task and starts the reduceStream operation.
func (rtm *reduceStreamTaskManager) AppendToTask(ctx context.Context, request *v1.ReduceRequest) error {
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

// CloseTask closes all reduceStream tasks matching the window's start and end time.
func (rtm *reduceStreamTaskManager) CloseTask(request *v1.ReduceRequest) error {
	if len(request.Operation.Windows) != 1 {
		return fmt.Errorf("close operation error: invalid number of windows")
	}

	window := request.Operation.Windows[0]
	wKey := windowKey(window)

	winCtx, ok := rtm.windowContexts[wKey]
	if !ok || len(winCtx.tasks) == 0 {
		return fmt.Errorf("close operation error: no tasks found for window")
	}

	tasksToWait := make([]*reduceStreamTask, 0, len(winCtx.tasks))

	for _, task := range winCtx.tasks {
		tasksToWait = append(tasksToWait, task)
		close(task.inputCh)
	}

	// Store the window for the goroutine
	win := winCtx.window

	// Remove the window context from the map
	delete(rtm.windowContexts, wKey)

	// Wait for all tasks to complete in a separate goroutine
	rtm.closeWg.Add(1)
	go func(tasks []*reduceStreamTask, win *v1.Window) {
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
				log.Printf("Error in reduceStream task: %v", t.err)
			}
		}

		if !hasErrors {
			// After all tasks are done, send the EOF response
			rtm.responseCh <- &v1.ReduceResponse{
				Window: win,
				EOF:    true,
			}
		}
	}(tasksToWait, win)

	return nil
}

// OutputChannel returns the output channel for the reduceStream task manager to read the results.
func (rtm *reduceStreamTaskManager) OutputChannel() <-chan *v1.ReduceResponse {
	return rtm.responseCh
}

// WaitAll waits for all the reduceStream tasks to complete.
func (rtm *reduceStreamTaskManager) WaitAll() {
	// First wait for all CloseTask goroutines to finish
	rtm.closeWg.Wait()

	windowContexts := make([]*windowContext, 0, len(rtm.windowContexts))
	for _, winCtx := range rtm.windowContexts {
		windowContexts = append(windowContexts, winCtx)
	}

	for _, winCtx := range windowContexts {
		var hasErrors bool

		for _, task := range winCtx.tasks {
			<-task.doneCh
			if task.err != nil {
				hasErrors = true
				log.Printf("Error in reduceStream task during WaitAll: %v", task.err)
			}
		}

		if !hasErrors {
			rtm.responseCh <- &v1.ReduceResponse{
				Window: winCtx.window,
				EOF:    true,
			}
		}
	}

	// after all the tasks are completed, close the output channel
	close(rtm.responseCh)
}

// CloseAll closes all the reduceStream tasks.
func (rtm *reduceStreamTaskManager) CloseAll() {
	for _, winCtx := range rtm.windowContexts {
		for _, task := range winCtx.tasks {
			close(task.inputCh)
		}
	}
}

func buildDatum(request *v1.ReduceRequest) Datum {
	return NewHandlerDatum(
		request.GetPayload().GetValue(),
		request.GetPayload().GetEventTime().AsTime(),
		request.GetPayload().GetWatermark().AsTime(),
		request.GetPayload().GetHeaders(),
	)
}
