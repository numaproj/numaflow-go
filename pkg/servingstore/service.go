package servingstore

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"sync"

	"golang.org/x/sync/errgroup"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	servingpb "github.com/numaproj/numaflow-go/pkg/apis/proto/serving/v1"
	"github.com/numaproj/numaflow-go/internal/shared"
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
	once         sync.Once
}

var errServingStorePanic = fmt.Errorf("UDF_EXECUTION_ERROR(%s)", shared.ContainerType)

// Put puts the payload into the Store.
func (s *Service) Put(ctx context.Context, request *servingpb.PutRequest) (*servingpb.PutResponse, error) {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic inside serving put handler: %v %v", r, string(debug.Stack()))
				st, _ := status.Newf(codes.Internal, "%s: %v", errServingStorePanic, r).WithDetails(&epb.DebugInfo{
					Detail: string(debug.Stack()),
				})
				g.Go(func() error { return st.Err() })
			}
		}()

		var payloads = make([]Payload, 0, len(request.Payloads))
		for _, payload := range request.Payloads {
			payloads = append(payloads, Payload{origin: payload.Origin, value: payload.Value})
		}

		s.ServingStore.Put(ctx, &PutRequest{id: request.Id, payloads: payloads})
		return nil
	})

	err := g.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		s.once.Do(func() {
			select {
			case s.shutdownCh <- struct{}{}:
				// signal enqueued
			default:
				log.Println("Shutdown signal already enqueued or watcher exited; skipping shutdown send")
			}
		})
	}
	return &servingpb.PutResponse{Success: err == nil}, err
}

// Get gets the data stored in the Store.
func (s *Service) Get(ctx context.Context, request *servingpb.GetRequest) (*servingpb.GetResponse, error) {
	g, ctx := errgroup.WithContext(ctx)

	var response *servingpb.GetResponse

	g.Go(func() error {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic inside serving get handler: %v %v", r, string(debug.Stack()))
				st, _ := status.Newf(codes.Internal, "%s: %v", errServingStorePanic, r).WithDetails(&epb.DebugInfo{
					Detail: string(debug.Stack()),
				})
				g.Go(func() error { return st.Err() })
			}
		}()

		storedResult := s.ServingStore.Get(ctx, &GetRequest{id: request.Id})

		var payloads = make([]*servingpb.Payload, 0)
		for _, payload := range storedResult.payloads {
			payloads = append(payloads, &servingpb.Payload{Origin: payload.origin, Value: payload.value})
		}

		response = &servingpb.GetResponse{Id: request.GetId(), Payloads: payloads}
		return nil
	})

	err := g.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		s.once.Do(func() {
			s.shutdownCh <- struct{}{}
		})
	}
	return response, err
}

// IsReady is used to indicate that the server is ready.
func (s *Service) IsReady(_ context.Context, _ *emptypb.Empty) (*servingpb.ReadyResponse, error) {
	return &servingpb.ReadyResponse{Ready: true}, nil
}
