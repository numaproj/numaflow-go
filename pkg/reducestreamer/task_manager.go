package reducestreamer

import (
	"context"
	"fmt"
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

// uniqueKey returns the unique key for the reduceStream task to be used in the task manager to identify the task.
func (rt *reduceStreamTask) uniqueKey() string {
	return fmt.Sprintf("%d:%d:%s",
		rt.window.GetStart().AsTime().UnixMilli(),
		rt.window.GetEnd().AsTime().UnixMilli(),
		strings.Join(rt.keys, delimiter))
}

// reduceStreamTaskManager manages the reduceStream tasks.
type reduceStreamTaskManager struct {
	creatorHandle ReduceStreamerCreator
	tasks         map[string]*reduceStreamTask
	responseCh    chan *v1.ReduceResponse
}

func newReduceTaskManager(reduceStreamerCreator ReduceStreamerCreator) *reduceStreamTaskManager {
	return &reduceStreamTaskManager{
		creatorHandle: reduceStreamerCreator,
		tasks:         make(map[string]*reduceStreamTask),
		responseCh:    make(chan *v1.ReduceResponse),
	}
}

// CreateTask creates a new reduceStream task and starts the  reduceStream operation.
func (rtm *reduceStreamTaskManager) CreateTask(ctx context.Context, request *v1.ReduceRequest) error {
	if len(request.Operation.Windows) != 1 {
		return fmt.Errorf("create operation error: invalid number of windows")
	}

	md := NewMetadata(NewIntervalWindow(request.Operation.Windows[0].GetStart().AsTime(),
		request.Operation.Windows[0].GetEnd().AsTime()))

	task := &reduceStreamTask{
		keys:     request.GetPayload().GetKeys(),
		window:   request.Operation.Windows[0],
		inputCh:  make(chan Datum),
		outputCh: make(chan Message),
		doneCh:   make(chan struct{}),
	}

	key := task.uniqueKey()
	rtm.tasks[key] = task

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

		reduceStreamerHandle := rtm.creatorHandle.Create()
		// invoke the reduceStream function
		reduceStreamerHandle.ReduceStream(ctx, request.GetPayload().GetKeys(), task.inputCh, task.outputCh, md)
		// close the output channel after the reduceStream function is done
		close(task.outputCh)
		// wait for the output to be forwarded
		wg.Wait()
		// send a done signal
		close(task.doneCh)
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

	gKey := generateKey(request.Operation.Windows[0], request.Payload.Keys)
	task, ok := rtm.tasks[gKey]

	// if the task is not found, create a new task
	if !ok {
		return rtm.CreateTask(ctx, request)
	}

	task.inputCh <- buildDatum(request)
	return nil
}

// OutputChannel returns the output channel for the reduceStream task manager to read the results.
func (rtm *reduceStreamTaskManager) OutputChannel() <-chan *v1.ReduceResponse {
	return rtm.responseCh
}

// WaitAll waits for all the reduceStream tasks to complete.
func (rtm *reduceStreamTaskManager) WaitAll() {
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

// CloseAll closes all the reduceStream tasks.
func (rtm *reduceStreamTaskManager) CloseAll() {
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
