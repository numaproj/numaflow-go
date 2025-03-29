package sideinput

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"sync"

	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	sideinputpb "github.com/numaproj/numaflow-go/pkg/apis/proto/sideinput/v1"
	"github.com/numaproj/numaflow-go/pkg/shared"
)

const (
	uds                   = "unix"
	address               = "/var/run/numaflow/sideinput.sock"
	DirPath               = "/var/numaflow/side-inputs"
	defaultMaxMessageSize = 1024 * 1024 * 64 // 64MB
	serverInfoFilePath    = "/var/run/numaflow/sideinput-server-info"
)

var containerType = func() string {
	if val, exists := os.LookupEnv(shared.EnvUDContainerType); exists {
		return val
	}
	return "unknown-container"
}()

var errSideInputHandlerPanic = fmt.Errorf("UDF_EXECUTION_ERROR(%s)", containerType)

// Service implements the proto gen server interface and contains the retrieve operation handler
type Service struct {
	sideinputpb.UnimplementedSideInputServer
	Retriever  SideInputRetriever
	shutdownCh chan<- struct{}
	once       sync.Once
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
			fs.once.Do(func() {
				fs.shutdownCh <- struct{}{}
			})
			st, _ := status.Newf(codes.Internal, "%s: %v", errSideInputHandlerPanic, r).WithDetails(&epb.DebugInfo{
				Detail: string(debug.Stack()),
			})
			err = st.Err()
		}
	}()
	messageSi := fs.Retriever.RetrieveSideInput(ctx)
	resp = &sideinputpb.SideInputResponse{
		Value:       messageSi.value,
		NoBroadcast: messageSi.noBroadcast,
	}
	return resp, nil
}
