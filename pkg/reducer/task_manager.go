package reducer

import (
	"context"
	"fmt"
	"strings"

	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/reduce/v1"
)

// reduceTask represents a task for a performing reduceStream operation.
type reduceTask struct {
	keys     []string
	window   *v1.Window
	reducer  Reducer
	inputCh  chan Datum
	outputCh chan Message
	doneCh   chan struct{}
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

// reduceTaskManager manages the reduce tasks for a reduce operation.
type reduceTaskManager struct {
	reducerCreatorHandle ReducerCreator
	tasks                map[string]*reduceTask
	responseCh           chan *v1.ReduceResponse
}

func newReduceTaskManager(reducerCreatorHandle ReducerCreator) *reduceTaskManager {
	return &reduceTaskManager{
		reducerCreatorHandle: reducerCreatorHandle,
		tasks:                make(map[string]*reduceTask),
		responseCh:           make(chan *v1.ReduceResponse),
	}
}

// CreateTask creates a new reduce task and starts the reduce operation.
func (rtm *reduceTaskManager) CreateTask(ctx context.Context, request *v1.ReduceRequest) error {
	if len(request.Operation.Windows) != 1 {
		return fmt.Errorf("create operation error: invalid number of windows")
	}

	md := NewMetadata(NewIntervalWindow(request.Operation.Windows[0].GetStart().AsTime(),
		request.Operation.Windows[0].GetEnd().AsTime()))

	task := &reduceTask{
		keys:     request.GetPayload().GetKeys(),
		window:   request.Operation.Windows[0],
		inputCh:  make(chan Datum),
		outputCh: make(chan Message),
		doneCh:   make(chan struct{}),
	}

	key := task.uniqueKey()
	rtm.tasks[key] = task

	go func() {
		// invoke the reduce function
		// create a new reducer, since we got a new key
		reducerHandle := rtm.reducerCreatorHandle.Create()
		messages := reducerHandle.Reduce(ctx, request.GetPayload().GetKeys(), task.inputCh, md)

		for _, message := range messages {
			// write the output to the output channel, service will forward it to downstream
			rtm.responseCh <- task.buildReduceResponse(message)
		}
		// close the output channel after the reduce function is done
		close(task.outputCh)
		// send a done signal
		close(task.doneCh)
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

	task, ok := rtm.tasks[generateKey(request.Operation.Windows[0], request.Payload.Keys)]

	// if the task is not found, create a new task
	if !ok {
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
	var eofResponse *v1.ReduceResponse
	for _, task := range rtm.tasks {
		<-task.doneCh
		if eofResponse == nil {
			eofResponse = task.buildEOFResponse()
		}
	}
	rtm.responseCh <- eofResponse
	// after all the tasks are completed, close the output channel
	close(rtm.responseCh)
}

// CloseAll closes all the reduce tasks.
func (rtm *reduceTaskManager) CloseAll() {
	for _, task := range rtm.tasks {
		close(task.inputCh)
	}
}

func generateKey(window *v1.Window, keys []string) string {
	return fmt.Sprintf("%d:%d:%s",
		window.GetStart().AsTime().UnixMilli(),
		window.GetEnd().AsTime().UnixMilli(),
		strings.Join(keys, delimiter))
}

func buildDatum(request *v1.ReduceRequest) Datum {
	return NewHandlerDatum(
		request.GetPayload().GetValue(),
		request.GetPayload().GetEventTime().AsTime(),
		request.GetPayload().GetWatermark().AsTime(),
		request.GetPayload().GetHeaders(),
	)
}
