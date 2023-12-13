package reducer

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	grpcmd "google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	reducepb "github.com/numaproj/numaflow-go/pkg/apis/proto/reduce/v1"
)

const (
	uds                   = "unix"
	defaultMaxMessageSize = 1024 * 1024 * 64
	address               = "/var/run/numaflow/reduce.sock"
	winStartTime          = "x-numaflow-win-start-time"
	winEndTime            = "x-numaflow-win-end-time"
	delimiter             = ":"
)

// Service implements the proto gen server interface and contains the reduce operation handler.
type Service struct {
	reducepb.UnimplementedReduceServer
	CreateReduceHandler CreateReducer
}

// IsReady returns true to indicate the gRPC connection is ready.
func (fs *Service) IsReady(context.Context, *emptypb.Empty) (*reducepb.ReadyResponse, error) {
	return &reducepb.ReadyResponse{Ready: true}, nil
}

// ReduceFn applies a reduce function to a request stream and returns a list of results.
func (fs *Service) ReduceFn(stream reducepb.Reduce_ReduceFnServer) error {
	var (
		md        Metadata
		err       error
		startTime int64
		endTime   int64
		ctx       = stream.Context()
		chanMap   = make(map[string]chan Datum)
		mu        sync.RWMutex
		g         errgroup.Group
	)

	grpcMD, ok := grpcmd.FromIncomingContext(ctx)
	if !ok {
		statusErr := status.Errorf(codes.InvalidArgument, "keys and window information are not passed in grpc metadata")
		return statusErr
	}

	// get window start and end time from grpc metadata
	var st, et string
	st, err = getValueFromMetadata(grpcMD, winStartTime)
	if err != nil {
		statusErr := status.Errorf(codes.InvalidArgument, err.Error())
		return statusErr
	}

	et, err = getValueFromMetadata(grpcMD, winEndTime)
	if err != nil {
		statusErr := status.Errorf(codes.InvalidArgument, err.Error())
		return statusErr
	}

	startTime, _ = strconv.ParseInt(st, 10, 64)
	endTime, _ = strconv.ParseInt(et, 10, 64)

	// create interval window interface using the start and end time
	iw := NewIntervalWindow(time.UnixMilli(startTime), time.UnixMilli(endTime))

	// create metadata using interval window interface
	md = NewMetadata(iw)

	// read messages from the stream and write the messages to corresponding channels
	// if the channel is not created, create the channel and invoke the reduceFn
	for {
		d, recvErr := stream.Recv()
		// if EOF, close all the channels
		if recvErr == io.EOF {
			closeChannels(chanMap)
			break
		}
		if recvErr != nil {
			closeChannels(chanMap)
			// the error here is returned by stream.Recv()
			// it's already a gRPC error
			return recvErr
		}
		unifiedKey := strings.Join(d.GetKeys(), delimiter)
		var hd = NewHandlerDatum(d.GetValue(), d.EventTime.AsTime(), d.Watermark.AsTime())

		ch, chok := chanMap[unifiedKey]
		if !chok {
			ch = make(chan Datum)
			chanMap[unifiedKey] = ch

			func(k []string, ch chan Datum) {
				g.Go(func() error {
					// we stream the messages to the user by writing messages
					// to channel and wait until we get the result and stream
					// the result back to the client (numaflow).

					// create a new reducer, since we got a new key
					reducer := fs.CreateReduceHandler.Create()
					messages := reducer.Reduce(ctx, k, ch, md)
					datumList := buildDatumList(messages)

					// stream.Send() is not thread safe.
					mu.Lock()
					defer mu.Unlock()
					sendErr := stream.Send(datumList)
					if sendErr != nil {
						// the error here is returned by stream.Send()
						// it's already a gRPC error
						return sendErr
					}
					return nil
				})
			}(d.GetKeys(), ch)
		}
		ch <- hd
	}

	// wait until all the mapfn return
	return g.Wait()
}

func buildDatumList(messages Messages) *reducepb.ReduceResponse {
	response := &reducepb.ReduceResponse{}
	for _, msg := range messages {
		response.Results = append(response.Results, &reducepb.ReduceResponse_Result{
			Keys:  msg.Keys(),
			Value: msg.Value(),
			Tags:  msg.Tags(),
		})
	}

	return response
}

func closeChannels(chanMap map[string]chan Datum) {
	for _, ch := range chanMap {
		close(ch)
	}
}

func getValueFromMetadata(md grpcmd.MD, k string) (string, error) {
	var value string

	keyValue := md.Get(k)

	if len(keyValue) > 1 {
		return value, fmt.Errorf("expected extactly one value for keys %s in metadata but got %d values, %s", k, len(keyValue), keyValue)
	} else if len(keyValue) == 1 {
		value = keyValue[0]
	} else {
		// the length equals zero is invalid for reduce
		// since we are using a global keys, and start and end time
		// cannot be empty
		return value, fmt.Errorf("expected non empty value for keys %s in metadata but got an empty value", k)
	}
	return value, nil
}
