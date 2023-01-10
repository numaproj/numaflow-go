package source

import (
	"context"
	"fmt"
	"time"

	sourcepb "github.com/numaproj/numaflow-go/pkg/apis/proto/source/v1"
	grpcmd "google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

// handlerDatum implements the Datum interface and is used in the transform handlers.
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

// Service implements the proto gen server interface and contains the transform operation handler.
type Service struct {
	sourcepb.UnimplementedUserDefinedSourceTransformerServer
	Transformer TransformHandler
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*sourcepb.ReadyResponse, error) {
	return &sourcepb.ReadyResponse{Ready: true}, nil
}

// TransformFn applies a transformation to each datum element
func (fs *Service) TransformFn(ctx context.Context, d *sourcepb.Datum) (*sourcepb.DatumList, error) {
	var key string

	// get key from gPRC metadata
	if grpcMD, ok := grpcmd.FromIncomingContext(ctx); ok {
		keyValue := grpcMD.Get(DatumKey)
		if len(keyValue) > 1 {
			return nil, fmt.Errorf("expect extact one key but got %d keys", len(keyValue))
		} else if len(keyValue) == 1 {
			key = keyValue[0]
		} else {
			// do nothing: the length equals zero is valid, meaning the key is an empty string ""
		}
	}
	var hd = handlerDatum{
		value:     d.GetValue(),
		eventTime: d.GetEventTime().EventTime.AsTime(),
		watermark: d.GetWatermark().Watermark.AsTime(),
	}
	messages := fs.Transformer.HandleTransform(ctx, key, &hd)
	var elements []*sourcepb.Datum
	for _, m := range messages.Items() {
		elements = append(elements, &sourcepb.Datum{
			// TODO - extract event time from message
			EventTime: nil,
			Key:       m.Key,
			Value:     m.Value,
		})
	}
	datumList := &sourcepb.DatumList{
		Elements: elements,
	}
	return datumList, nil
}
