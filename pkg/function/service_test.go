package function

import (
	"context"
	"reflect"
	"testing"
	"time"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	"github.com/numaproj/numaflow-go/pkg/configs"
	"github.com/numaproj/numaflow-go/pkg/datum"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestService_MapFn(t *testing.T) {
	type fields struct {
		Mapper  MapHandler
		Reducer ReduceHandler
	}
	type args struct {
		ctx context.Context
		d   *functionpb.Datum
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *functionpb.DatumList
		wantErr bool
	}{
		{
			name: "map_fn_forward_msg",
			fields: fields{
				Mapper: MapFunc(func(ctx context.Context, key string, datum datum.Datum) Messages {
					msg := datum.Value()
					return MessagesBuilder().Append(MessageTo(key+"_test", msg))
				}),
				Reducer: nil,
			},
			args: args{
				ctx: context.Background(),
				d: &functionpb.Datum{
					Key:       "client",
					Value:     []byte(`test`),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			want: &functionpb.DatumList{
				Elements: []*functionpb.Datum{
					{
						Key:   "client_test",
						Value: []byte(`test`),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &Service{
				Mapper:  tt.fields.Mapper,
				Reducer: tt.fields.Reducer,
			}
			// here's a trick for testing:
			// because we are not using gRPC, we directly set a new incoming ctx
			// instead of the regular outgoing context in the real gRPC connection.
			ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{configs.DatumKey: "client"}))
			got, err := fs.MapFn(ctx, tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapFn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapFn() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_IsReady(t *testing.T) {
	type fields struct {
		Mapper  MapHandler
		Reducer ReduceHandler
	}
	type args struct {
		in0 context.Context
		in1 *emptypb.Empty
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *functionpb.ReadyResponse
		wantErr bool
	}{
		{
			name: "is_ready",
			fields: fields{
				Mapper:  nil,
				Reducer: nil,
			},
			args: args{
				in0: nil,
				in1: &emptypb.Empty{},
			},
			want: &functionpb.ReadyResponse{
				Ready: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &Service{
				Mapper:  tt.fields.Mapper,
				Reducer: tt.fields.Reducer,
			}
			got, err := fs.IsReady(tt.args.in0, tt.args.in1)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsReady() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsReady() got = %v, want %v", got, tt.want)
			}
		})
	}
}
