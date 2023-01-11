package function

import (
	"context"
	"github.com/numaproj/numaflow-go/pkg/function/types"
	"reflect"
	"testing"
	"time"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	grpcmd "google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type fields struct {
	mapper  MapHandler
	mapperT MapTHandler
	reducer ReduceHandler
}

func TestService_MapFn(t *testing.T) {
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
				mapper: MapFunc(func(ctx context.Context, key string, datum Datum) types.Messages {
					msg := datum.Value()
					return types.MessagesBuilder().Append(types.MessageTo(key+"_test", msg))
				}),
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
		{
			name: "map_fn_forward_msg_forward_to_all",
			fields: fields{
				mapper: MapFunc(func(ctx context.Context, key string, datum Datum) types.Messages {
					msg := datum.Value()
					return types.MessagesBuilder().Append(types.MessageToAll(msg))
				}),
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
						Key:   types.ALL,
						Value: []byte(`test`),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "map_fn_forward_msg_drop_msg",
			fields: fields{
				mapper: MapFunc(func(ctx context.Context, key string, datum Datum) types.Messages {
					return types.MessagesBuilder().Append(types.MessageToDrop())
				}),
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
						Key:   types.DROP,
						Value: []byte{},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &Service{
				Mapper:  tt.fields.mapper,
				MapperT: tt.fields.mapperT,
				Reducer: tt.fields.reducer,
			}
			// here's a trick for testing:
			// because we are not using gRPC, we directly set a new incoming ctx
			// instead of the regular outgoing context in the real gRPC connection.
			ctx := grpcmd.NewIncomingContext(context.Background(), grpcmd.New(map[string]string{DatumKey: "client"}))
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

func TestService_MapTFn(t *testing.T) {
	type args struct {
		ctx context.Context
		d   *functionpb.Datum
	}

	testTime := time.Date(2021, 8, 15, 14, 30, 45, 100, time.Local)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *functionpb.DatumList
		wantErr bool
	}{
		{
			name: "mapT_fn_forward_msg",
			fields: fields{
				mapperT: MapTFunc(func(ctx context.Context, key string, datum Datum) types.MessageTs {
					msg := datum.Value()
					return types.MessageTsBuilder().Append(types.MessageTTo(testTime, key+"_test", msg))
				}),
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
						EventTime: &functionpb.EventTime{EventTime: timestamppb.New(testTime)},
						Key:       "client_test",
						Value:     []byte(`test`),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "mapT_fn_forward_msg_forward_to_all",
			fields: fields{
				mapperT: MapTFunc(func(ctx context.Context, key string, datum Datum) types.MessageTs {
					msg := datum.Value()
					return types.MessageTsBuilder().Append(types.MessageTToAll(testTime, msg))
				}),
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
						EventTime: &functionpb.EventTime{EventTime: timestamppb.New(testTime)},
						Key:       types.ALL,
						Value:     []byte(`test`),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "mapT_fn_forward_msg_drop_msg",
			fields: fields{
				mapperT: MapTFunc(func(ctx context.Context, key string, datum Datum) types.MessageTs {
					return types.MessageTsBuilder().Append(types.MessageTToDrop())
				}),
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
						EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
						Key:       types.DROP,
						Value:     []byte{},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &Service{
				Mapper:  tt.fields.mapper,
				MapperT: tt.fields.mapperT,
				Reducer: tt.fields.reducer,
			}
			// here's a trick for testing:
			// because we are not using gRPC, we directly set a new incoming ctx
			// instead of the regular outgoing context in the real gRPC connection.
			ctx := grpcmd.NewIncomingContext(context.Background(), grpcmd.New(map[string]string{DatumKey: "client"}))
			got, err := fs.MapTFn(ctx, tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapTFn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapTFn() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO - add TestService_ReduceFn

func TestService_IsReady(t *testing.T) {
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
				mapper:  nil,
				reducer: nil,
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
				Mapper:  tt.fields.mapper,
				Reducer: tt.fields.reducer,
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
