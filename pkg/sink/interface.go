package sink

import (
	"context"
	"time"

	sinkpb "github.com/numaproj/numaflow-go/pkg/apis/proto/sink/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Datum interface {
	Value() []byte
	EventTime() time.Time
	Watermark() time.Time
	ID() string
}

// Client contains methods to call a gRPC client.
type Client interface {
	CloseConn(ctx context.Context) error
	IsReady(ctx context.Context, in *emptypb.Empty) (bool, error)
	SinkFn(ctx context.Context, datumStreamCh <-chan *sinkpb.Datum) ([]*sinkpb.Response, error)
}

// SinkHandler is the interface of sink function implementation.
type SinkHandler interface {
	// HandleDo is the function to process a list of incoming messages
	HandleDo(ctx context.Context, datumStreamCh <-chan Datum) Responses
}

// SinkFunc is utility type used to convert a HandleDo function to a SinkHandler.
type SinkFunc func(ctx context.Context, datumStreamCh <-chan Datum) Responses

// HandleDo implements the function of sink function.
func (sf SinkFunc) HandleDo(ctx context.Context, datumStreamCh <-chan Datum) Responses {
	return sf(ctx, datumStreamCh)
}
