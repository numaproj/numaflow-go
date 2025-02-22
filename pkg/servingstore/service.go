package servingstore

import (
	"context"
	"errors"
	"log"
	"runtime/debug"

	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	servingpb "github.com/numaproj/numaflow-go/pkg/apis/proto/serving/v1"
)

const (
	uds                   = "unix"
	address               = "/var/run/numaflow/serving.sock"
	defaultMaxMessageSize = 1024 * 1024 * 64 // 64MB
	serverInfoFilePath    = "/var/run/numaflow/serving-server-info"
)

// Service implements the proto gen server interface
type Service struct {
	servingpb.UnimplementedServingStoreServer
	ServingStore ServingStorer
	shutdownCh   chan<- struct{}
}

var errServingStorePanic = errors.New("UDF_EXECUTION_ERROR(serving)")

func handlePanic() (err error) {
	if r := recover(); r != nil {
		log.Printf("panic inside map handler: %v %v", r, string(debug.Stack()))
		st, _ := status.Newf(codes.Internal, "%s: %v", errServingStorePanic, r).WithDetails(&epb.DebugInfo{
			Detail: string(debug.Stack()),
		})
		err = st.Err()
	}

	return err
}

// Put puts tine payload into the Store.
func (s *Service) Put(ctx context.Context, request *servingpb.PutRequest) (*servingpb.PutResponse, error) {
	var err error
	// handle panic
	defer func() { err = handlePanic() }()

	var payloads = make([][]byte, 0, len(request.Payloads))
	for _, payload := range request.Payloads {
		payloads = append(payloads, payload.Value)
	}

	s.ServingStore.Put(ctx, &PutRequest{origin: request.Origin, payloads: payloads})

	return &servingpb.PutResponse{Success: true}, err
}

// Get gets the data stored in the Store.
func (s *Service) Get(ctx context.Context, request *servingpb.GetRequest) (*servingpb.GetResponse, error) {
	var err error
	// handle panic
	defer func() { err = handlePanic() }()

	storedResults := s.ServingStore.Get(ctx, &GetRequest{id: request.Id})

	items := storedResults.Items()
	var payloads = make([]*servingpb.OriginalPayload, 0, len(items))

	for _, storedResult := range items {
		var p = make([]*servingpb.Payload, 0)
		for _, payload := range storedResult.payloads {
			p = append(p, &servingpb.Payload{Id: request.GetId(), Value: payload.value})
		}
		payloads = append(payloads, &servingpb.OriginalPayload{
			Origin:   storedResult.origin,
			Payloads: p,
		})
	}

	return &servingpb.GetResponse{Payloads: payloads}, err
}

// IsReady is used to indicate that the server is ready.
func (s *Service) IsReady(_ context.Context, _ *emptypb.Empty) (*servingpb.ReadyResponse, error) {
	return &servingpb.ReadyResponse{Ready: true}, nil
}
