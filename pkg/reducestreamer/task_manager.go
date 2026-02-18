package reducestreamer

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"strings"
	"sync"

	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/reduce/v1"
	"github.com/numaproj/numaflow-go/pkg/internal/shared"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// reduceStreamTask represents a task for a performing reduceStream operation.
type reduceStreamTask struct {
	keys     []string
	window   *v1.Window
	inputCh  chan Datum
	outputCh chan Message
	doneCh   chan struct{}
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

var errReduceStreamHandlerPanic = fmt.Errorf("UDF_EXECUTION_ERROR(%s)", shared.ContainerType)

// reduceStreamTaskManager manages the reduceStream tasks.
type reduceStreamTaskManager struct {
	creatorHandle ReduceStreamerCreator
	tasks         map[string]*reduceStreamTask
	responseCh    chan *v1.ReduceResponse
	errorCh       chan error
	ctx           context.Context
}

func newReduceTaskManager(ctx context.Context, reduceStreamerCreator ReduceStreamerCreator) *reduceStreamTaskManager {
	return &reduceStreamTaskManager{
		creatorHandle: reduceStreamerCreator,
		tasks:         make(map[string]*reduceStreamTask),
		responseCh:    make(chan *v1.ReduceResponse),
		errorCh:       make(chan error, 1), // buffered channel to avoid blocking
		ctx:           ctx,
	}
}

// CreateTask creates a new reduceStream task and starts the  reduceStream operation.
func (rtm *reduceStreamTaskManager) CreateTask(request *v1.ReduceRequest) error {
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

		// handle panic and ensure doneCh is always closed
		defer func() {
			// Always close doneCh, even if panic occurs
			close(task.doneCh)

			if r := recover(); r != nil {
				log.Printf("panic inside reduce streamer handler: %v %v", r, string(debug.Stack()))
				st, _ := status.Newf(codes.Internal, "%s: %v", errReduceStreamHandlerPanic, r).WithDetails(&epb.DebugInfo{
					Detail: string(debug.Stack()),
				})
				// Non-blocking send - if channel is full or closed, we don't care since one panic is enough to trigger shutdown
				select {
				case rtm.errorCh <- st.Err():
				case <-rtm.ctx.Done():
					// Context is cancelled, don't try to send error
				default:
					// Channel is full or closed, its fine since we only need one panic to trigger shutdown
				}
			}
		}()

		reduceStreamerHandle := rtm.creatorHandle.Create()
		// invoke the reduceStream function
		reduceStreamerHandle.ReduceStream(rtm.ctx, request.GetPayload().GetKeys(), task.inputCh, task.outputCh, md)
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
func (rtm *reduceStreamTaskManager) AppendToTask(request *v1.ReduceRequest) error {
	if len(request.Operation.Windows) != 1 {
		return fmt.Errorf("append operation error: invalid number of windows")
	}

	gKey := generateKey(request.Operation.Windows[0], request.Payload.Keys)
	task, ok := rtm.tasks[gKey]

	// if the task is not found, create a new task
	if !ok {
		return rtm.CreateTask(request)
	}

	task.inputCh <- buildDatum(request)
	return nil
}

// OutputChannel returns the output channel for the reduceStream task manager to read the results.
func (rtm *reduceStreamTaskManager) OutputChannel() <-chan *v1.ReduceResponse {
	return rtm.responseCh
}

// ErrorChannel returns the error channel for the reduceStream task manager to read the errors.
func (rtm *reduceStreamTaskManager) ErrorChannel() <-chan error {
	return rtm.errorCh
}

// CloseErrorChannel closes the error channel to signal no more errors will be sent.
// This should be called after all tasks have completed.
func (rtm *reduceStreamTaskManager) CloseErrorChannel() {
	close(rtm.errorCh)
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
