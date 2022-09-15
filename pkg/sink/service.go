package sink

import (
	"context"
	"time"

	sinkpb "github.com/numaproj/numaflow-go/pkg/apis/proto/sink/v1"
	"github.com/numaproj/numaflow-go/pkg/datum"
	"google.golang.org/protobuf/types/known/emptypb"
)

// handlerDatum implements the Datum interface and is used in the map and reduce handlers.
type handlerDatum struct {
	value     []byte
	eventTime time.Time
	watermark time.Time
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

// Service implements the proto gen server interface and contains the map operation handler and the reduce operation handler.
type Service struct {
	sinkpb.UnimplementedUserDefinedSinkServer

	Sinker SinkHandler
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*sinkpb.ReadyResponse, error) {
	return &sinkpb.ReadyResponse{Ready: true}, nil
}

// SinkFn applies a function to a list of datum element.
func (fs *Service) SinkFn(ctx context.Context, datumList *sinkpb.DatumList) (*sinkpb.ResponseList, error) {
	var hdList []datum.Datum
	for _, d := range datumList.GetElements() {
		hdList = append(hdList, &handlerDatum{
			value:     d.GetValue(),
			eventTime: d.GetEventTime().EventTime.AsTime(),
			watermark: d.GetWatermark().Watermark.AsTime(),
		})
	}
	messages := fs.Sinker.HandleDo(ctx, hdList)
	var responses []*sinkpb.Response
	for _, m := range messages.Items() {
		responses = append(responses, &sinkpb.Response{
			Id:      m.ID,
			Success: m.Success,
			ErrMsg:  m.Err,
		})
	}
	responseList := &sinkpb.ResponseList{
		Responses: responses,
	}
	return responseList, nil
}
