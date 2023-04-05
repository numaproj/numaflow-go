package sink

import (
	"context"
	"io"
	"sync"
	"time"

	sinkpb "github.com/numaproj/numaflow-go/pkg/apis/proto/sink/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// handlerDatum implements the Datum interface and is used in the sink handlers.
type handlerDatum struct {
	keys      []string
	value     []byte
	eventTime time.Time
	watermark time.Time
	metadata  handlerDatumMetadata
}

func (h *handlerDatum) Keys() []string {
	return h.keys
}

func (h *handlerDatum) Value() []byte {
	return h.value
}

func (h *handlerDatum) EventTime() time.Time {
	return h.eventTime
}

func (h *handlerDatum) Watermark() time.Time {
	return h.watermark
}

func (h *handlerDatum) Metadata() DatumMetadata {
	return h.metadata
}

// handlerDatumMetadata implements the DatumMetadata interface and is used in the sink handlers.
type handlerDatumMetadata struct {
	id           string
	numDelivered uint64
}

func (h handlerDatumMetadata) ID() string {
	return h.id
}

func (h handlerDatumMetadata) NumDelivered() uint64 {
	return h.numDelivered
}

// Service implements the proto gen server interface and contains the sink operation handler.
type Service struct {
	sinkpb.UnimplementedUserDefinedSinkServer

	Sinker SinkHandler
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*sinkpb.ReadyResponse, error) {
	return &sinkpb.ReadyResponse{Ready: true}, nil
}

// SinkFn applies a function to a list of datum element.
func (fs *Service) SinkFn(stream sinkpb.UserDefinedSink_SinkFnServer) error {
	var (
		responseList  []*sinkpb.Response
		wg            sync.WaitGroup
		datumStreamCh = make(chan Datum)
		ctx           = stream.Context()
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		messages := fs.Sinker.HandleDo(ctx, datumStreamCh)
		for _, msg := range messages {
			responseList = append(responseList, &sinkpb.Response{
				Id:      msg.ID,
				Success: msg.Success,
				ErrMsg:  msg.Err,
			})
		}
	}()

	for {
		d, err := stream.Recv()
		if err == io.EOF {
			close(datumStreamCh)
			break
		}
		if err != nil {
			close(datumStreamCh)
			// TODO: research on gRPC errors and revisit the error handler
			return err
		}
		var hd = &handlerDatum{
			value:     d.GetValue(),
			eventTime: d.GetEventTime().EventTime.AsTime(),
			watermark: d.GetWatermark().Watermark.AsTime(),
			metadata: handlerDatumMetadata{
				id:           d.GetMetadata().GetId(),
				numDelivered: d.GetMetadata().GetNumDelivered(),
			},
		}
		datumStreamCh <- hd
	}

	wg.Wait()
	return stream.SendAndClose(&sinkpb.ResponseList{
		Responses: responseList,
	})
}
