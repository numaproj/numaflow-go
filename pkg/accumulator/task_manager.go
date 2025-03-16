package accumulator

import (
	"context"
	"log"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/accumulator/v1"
)

const delimiter = ":"

// accumulateTask represents a task for performing accumulate operation.
type accumulateTask struct {
	keys            []string
	accumulator     Accumulator
	inputCh         chan Datum
	outputCh        chan Datum
	latestWatermark time.Time
}

// uniqueKey returns the unique key for the accumulate task to be used in the task manager to identify the task.
func (at *accumulateTask) uniqueKey() string {
	return strings.Join(at.keys, delimiter)
}

// accumulatorTaskManager manages the accumulate tasks for an accumulate operation.
type accumulatorTaskManager struct {
	accumulatorCreatorHandle AccumulatorCreator
	tasks                    map[string]*accumulateTask
	responseCh               chan *v1.AccumulatorResponse
	shutdownCh               chan<- struct{}
	mu                       sync.RWMutex
	eg                       *errgroup.Group
	ctx                      context.Context
}

// newAccumulatorTaskManager creates a new accumulator task manager.
func newAccumulatorTaskManager(ctx context.Context, eg *errgroup.Group, accumulatorCreatorHandle AccumulatorCreator) *accumulatorTaskManager {
	return &accumulatorTaskManager{
		accumulatorCreatorHandle: accumulatorCreatorHandle,
		tasks:                    make(map[string]*accumulateTask),
		responseCh:               make(chan *v1.AccumulatorResponse),
		eg:                       eg,
		ctx:                      ctx,
	}
}

// CreateTask creates a new accumulate task and starts the accumulate operation.
func (atm *accumulatorTaskManager) CreateTask(request *v1.AccumulatorRequest) {
	atm.mu.Lock()
	defer atm.mu.Unlock()

	task := &accumulateTask{
		keys:            request.GetPayload().GetKeys(),
		inputCh:         make(chan Datum, 500),
		outputCh:        make(chan Datum, 500),
		latestWatermark: time.UnixMilli(-1),
	}

	key := task.uniqueKey()
	atm.tasks[key] = task

	// starts the accumulate operation in a goroutine, any panic inside the accumulator handler is recovered and sent as an error response.
	atm.eg.Go(func() (err error) {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			// read responses from the output channel and send to the response channel
			for {
				select {
				case <-atm.ctx.Done():
					return
				case output, ok := <-task.outputCh:
					if !ok {
						// send EOF response to the response channel
						atm.responseCh <- &v1.AccumulatorResponse{
							EOF: true,
						}
						return
					}

					// update the latest watermark
					if output.Watermark().After(task.latestWatermark) {
						task.latestWatermark = output.Watermark()
					}

					select {
					case <-atm.ctx.Done():
						return
					case atm.responseCh <- &v1.AccumulatorResponse{
						Payload: &v1.Payload{
							Keys:      output.Keys(),
							Value:     output.Value(),
							EventTime: timestamppb.New(output.EventTime()),
							Watermark: timestamppb.New(output.Watermark()),
							Headers:   output.Headers(),
						},
						Window: &v1.Window{
							Start: timestamppb.New(time.UnixMilli(0)),
							// window end time is considered the latest watermark, based on the window end time, the compaction happens
							// on the client side.
							End:  timestamppb.New(task.latestWatermark),
							Slot: "slot-0",
						},
						Id:  output.ID(),
						EOF: false,
					}:
					}
				}
			}
		}()

		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic inside accumulator handler: %v %v", r, string(debug.Stack()))
				st, _ := status.Newf(codes.Internal, "%s: %v", errAccumulatorPanic, r).WithDetails(&epb.DebugInfo{
					Detail: string(debug.Stack()),
				})
				err = st.Err()
			}
		}()

		accumulatorHandle := atm.accumulatorCreatorHandle.Create()
		accumulatorHandle.Accumulate(atm.ctx, task.inputCh, task.outputCh)
		close(task.outputCh)

		wg.Wait()

		// remove the task from the task manager
		atm.mu.Lock()
		delete(atm.tasks, key)
		atm.mu.Unlock()
		return nil
	})

	select {
	case <-atm.ctx.Done():
		return
	case task.inputCh <- buildDatum(request):
	}
}

// AppendToTask writes the message to the accumulate task.
// If the task is not found, it creates a new task and starts the accumulate operation.
func (atm *accumulatorTaskManager) AppendToTask(request *v1.AccumulatorRequest) {
	atm.mu.RLock()
	task, ok := atm.tasks[generateKey(request.GetPayload().GetKeys())]
	atm.mu.RUnlock()

	if !ok {
		atm.CreateTask(request)
		return
	}
	select {
	case <-atm.ctx.Done():
		return
	case task.inputCh <- buildDatum(request):
	}
}

// CloseTask closes the accumulate task by closing the input channel.
func (atm *accumulatorTaskManager) CloseTask(request *v1.AccumulatorRequest) {
	atm.mu.Lock()
	defer atm.mu.Unlock()

	key := strings.Join(request.GetPayload().GetKeys(), delimiter)
	task, ok := atm.tasks[key]
	if !ok {
		return
	}

	close(task.inputCh)
}

// CloseAll closes all the accumulate tasks.
func (atm *accumulatorTaskManager) CloseAll() {
	atm.mu.Lock()
	defer atm.mu.Unlock()

	for _, task := range atm.tasks {
		close(task.inputCh)
	}
}

func (atm *accumulatorTaskManager) OutputChannel() <-chan *v1.AccumulatorResponse {
	return atm.responseCh
}

func generateKey(keys []string) string {
	return strings.Join(keys, delimiter)
}

func buildDatum(request *v1.AccumulatorRequest) Datum {
	return &handlerDatum{
		value:     request.GetPayload().GetValue(),
		eventTime: request.GetPayload().GetEventTime().AsTime(),
		watermark: request.GetPayload().GetWatermark().AsTime(),
		keys:      request.GetPayload().GetKeys(),
		headers:   request.GetPayload().GetHeaders(),
		id:        request.GetId(),
	}
}
