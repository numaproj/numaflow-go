package function

import "context"

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
