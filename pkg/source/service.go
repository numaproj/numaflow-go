package source

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	sourcepb "github.com/KeranYang/numaflow-go/pkg/apis/proto/source/v1"
)

// Service implements the proto gen server interface
type Service struct {
	sourcepb.UnimplementedUserDefinedSourceServer
	PendingHandler PendingHandler
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*sourcepb.ReadyResponse, error) {
	return &sourcepb.ReadyResponse{Ready: true}, nil
}

func (fs *Service) Pending(ctx context.Context, _ *emptypb.Empty) (*sourcepb.PendingResponse, error) {
	return &sourcepb.PendingResponse{Result: &sourcepb.PendingResponse_Result{
		Count: fs.PendingHandler.HandleDo(ctx),
	}}, nil
}
