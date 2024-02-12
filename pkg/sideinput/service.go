package sideinput

import (
	"context"

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

// Service implements the proto gen server interface and contains the retrieve operation handler
type Service struct {
	sideinputpb.UnimplementedSideInputServer
	Retriever SideInputRetriever
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*sideinputpb.ReadyResponse, error) {
	return &sideinputpb.ReadyResponse{Ready: true}, nil
}

// RetrieveSideInput applies the function for each side input retrieval request.
func (fs *Service) RetrieveSideInput(ctx context.Context, _ *emptypb.Empty) (*sideinputpb.SideInputResponse, error) {
	messageSi := fs.Retriever.RetrieveSideInput(ctx)
	var element *sideinputpb.SideInputResponse
	element = &sideinputpb.SideInputResponse{
		Value:       messageSi.value,
		NoBroadcast: messageSi.noBroadcast,
	}
	return element, nil
}
