package function

import (
	"context"

	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client interface {
	CloseConn(ctx context.Context) error
	IsReady(ctx context.Context, in *emptypb.Empty) (bool, error)
	DoFn(ctx context.Context, datum *v1.Datum) ([]*v1.Datum, error)
	ReduceFn(ctx context.Context, datumStream []*v1.Datum) ([]*v1.Datum, error)
}

// MapHandler is the interface of map function implementation.
type MapHandler interface {
	// HandleDo is the function to process each coming message
	HandleDo(ctx context.Context, key string, msg []byte) (Messages, error)
}

// ReduceHandler is the interface of reduce function implementation.
type ReduceHandler interface {
	// TODO
}

// DoFunc is utility type used to convert a HandleDo function to a MapHandler.
type DoFunc func(ctx context.Context, key string, msg []byte) (Messages, error)

// HandleDo implements the function of map function.
func (df DoFunc) HandleDo(ctx context.Context, key string, msg []byte) (Messages, error) {
	return df(ctx, key, msg)
}
