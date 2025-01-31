package sideinput

import (
	"context"
	"errors"
	"log"
	"runtime/debug"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	sideinputpb "github.com/numaproj/numaflow-go/pkg/apis/proto/sideinput/v1"
)

const (
	uds                   = "unix"
	address               = "/var/run/numaflow/sideinput.sock"
	DirPath               = "/var/numaflow/side-inputs"
	defaultMaxMessageSize = 1024 * 1024 * 64 // 64MB
	serverInfoFilePath    = "/var/run/numaflow/sideinput-server-info"
)

var errSideInputHandlerPanic = errors.New("USER_CODE_ERROR(side input)")

// Service implements the proto gen server interface and contains the retrieve operation handler
type Service struct {
	sideinputpb.UnimplementedSideInputServer
	Retriever  SideInputRetriever
	shutdownCh chan<- struct{}
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*sideinputpb.ReadyResponse, error) {
	return &sideinputpb.ReadyResponse{Ready: true}, nil
}

// RetrieveSideInput applies the function for each side input retrieval request.
func (fs *Service) RetrieveSideInput(ctx context.Context, _ *emptypb.Empty) (resp *sideinputpb.SideInputResponse, err error) {
	// handle panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic inside sideinput handler: %v %v", r, string(debug.Stack()))
			fs.shutdownCh <- struct{}{}
			err = status.Errorf(codes.Internal, "%s: %v %v", errSideInputHandlerPanic, r, string(debug.Stack()))
		}
	}()
	messageSi := fs.Retriever.RetrieveSideInput(ctx)
	resp = &sideinputpb.SideInputResponse{
		Value:       messageSi.value,
		NoBroadcast: messageSi.noBroadcast,
	}
	return resp, nil
}
