package reducer

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"strings"

	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/reduce/v1"
	"github.com/numaproj/numaflow-go/pkg/shared"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errReduceHandlerPanic = fmt.Errorf("UDF_EXECUTION_ERROR(%s)", shared.ContainerType)

// reduceTask represents a task for a performing reduceStream operation.
type reduceTask struct {
	keys    []string
	window  *v1.Window
	inputCh chan Datum
	doneCh  chan struct{}
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
	errorCh              chan error
}

func newReduceTaskManager(reducerCreatorHandle ReducerCreator) *reduceTaskManager {
	return &reduceTaskManager{
		reducerCreatorHandle: reducerCreatorHandle,
		tasks:                make(map[string]*reduceTask),
		responseCh:           make(chan *v1.ReduceResponse),
		errorCh:              make(chan error, 1),
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
		keys:    request.GetPayload().GetKeys(),
		window:  request.Operation.Windows[0],
		inputCh: make(chan Datum),
		doneCh:  make(chan struct{}),
	}

	key := task.uniqueKey()
	rtm.tasks[key] = task

	go func() {
		// handle panic
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic inside reduce handler: %v %v", r, string(debug.Stack()))
				st, _ := status.Newf(codes.Internal, "%s: %v", errReduceHandlerPanic, r).WithDetails(&epb.DebugInfo{
					Detail: string(debug.Stack()),
				})
				// Select is used here because errorCh is buffered (length 1).
				// If multiple panics occur before the error is read, only the first error will be sent;
				// This prevents the goroutine from blocking indefinitely if the channel is not being read.
				select {
				case rtm.errorCh <- st.Err():
				default:
					log.Printf("error channel full, shutting down")
				}

			}
		}()
		// invoke the reduce function
		// create a new reducer, since we got a new key
		reducerHandle := rtm.reducerCreatorHandle.Create()
		messages := reducerHandle.Reduce(ctx, request.GetPayload().GetKeys(), task.inputCh, md)
		for _, message := range messages {
			// write the output to the output channel, service will forward it to downstream
			rtm.responseCh <- task.buildReduceResponse(message)
		}
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

	// If the task is not found, create a new task
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

// Method to get the error channel
func (rtm *reduceTaskManager) ErrorChannel() <-chan error {
	return rtm.errorCh
}

// WaitAll waits for all the reduce tasks to complete.
// If the context is done, it will stop waiting and return.
func (rtm *reduceTaskManager) WaitAll(ctx context.Context) {
	var eofResponse *v1.ReduceResponse
	allDone := true
loop:
	for _, task := range rtm.tasks {
		select {
		case <-task.doneCh:
			if eofResponse == nil {
				eofResponse = task.buildEOFResponse()
			}
		case <-ctx.Done():
			allDone = false
			break loop
		}
	}
	if allDone {
		rtm.responseCh <- eofResponse
	}
	close(rtm.responseCh)
	close(rtm.errorCh)
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
