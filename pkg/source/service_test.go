package source

import (
	"context"
	"reflect"
	"testing"
	"time"

	sourcepb "github.com/numaproj/numaflow-go/pkg/apis/proto/source/v1"

	grpcmd "google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestService_TransformFn(t *testing.T) {
	type fields struct {
		Transformer TransformHandler
	}
	type args struct {
		ctx context.Context
		d   *sourcepb.Datum
	}

	testTime := time.Date(2021, 8, 15, 14, 30, 45, 100, time.Local)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *sourcepb.DatumList
		wantErr bool
	}{
		{
			name: "map_fn_forward_msg",
			fields: fields{
				Transformer: TransformFunc(func(ctx context.Context, key string, datum Datum) Messages {
					msg := datum.Value()
					return MessagesBuilder().Append(MessageTo(testTime, key+"_test", msg))
				}),
			},
			args: args{
				ctx: context.Background(),
				d: &sourcepb.Datum{
					Key:       "client",
					Value:     []byte(`test`),
					EventTime: &sourcepb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &sourcepb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			want: &sourcepb.DatumList{
				Elements: []*sourcepb.Datum{
					{
						EventTime: &sourcepb.EventTime{EventTime: timestamppb.New(testTime)},
						Key:       "client_test",
						Value:     []byte(`test`),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "map_fn_forward_msg_forward_to_all",
			fields: fields{
				Transformer: TransformFunc(func(ctx context.Context, key string, datum Datum) Messages {
					msg := datum.Value()
					return MessagesBuilder().Append(MessageToAll(testTime, msg))
				}),
			},
			args: args{
				ctx: context.Background(),
				d: &sourcepb.Datum{
					Key:       "client",
					Value:     []byte(`test`),
					EventTime: &sourcepb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &sourcepb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			want: &sourcepb.DatumList{
				Elements: []*sourcepb.Datum{
					{
						EventTime: &sourcepb.EventTime{EventTime: timestamppb.New(testTime)},
						Key:       ALL,
						Value:     []byte(`test`),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "map_fn_drop_msg",
			fields: fields{
				Transformer: TransformFunc(func(ctx context.Context, key string, datum Datum) Messages {
					return MessagesBuilder().Append(MessageToDrop())
				}),
			},
			args: args{
				ctx: context.Background(),
				d: &sourcepb.Datum{
					Key:       "client",
					Value:     []byte(`test`),
					EventTime: &sourcepb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &sourcepb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			want: &sourcepb.DatumList{
				Elements: []*sourcepb.Datum{
					{
						EventTime: &sourcepb.EventTime{EventTime: timestamppb.New(time.Time{})},
						Key:       DROP,
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
				Transformer: tt.fields.Transformer,
			}
			// here's a trick for testing:
			// because we are not using gRPC, we directly set a new incoming ctx
			// instead of the regular outgoing context in the real gRPC connection.
			ctx := grpcmd.NewIncomingContext(context.Background(), grpcmd.New(map[string]string{DatumKey: "client"}))
			got, err := fs.TransformFn(ctx, tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapFn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransformFn() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_IsReady(t *testing.T) {
	type fields struct {
		Transformer TransformHandler
	}
	type args struct {
		in0 context.Context
		in1 *emptypb.Empty
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *sourcepb.ReadyResponse
		wantErr bool
	}{
		{
			name: "is_ready",
			fields: fields{
				Transformer: nil,
			},
			args: args{
				in0: nil,
				in1: &emptypb.Empty{},
			},
			want: &sourcepb.ReadyResponse{
				Ready: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &Service{
				Transformer: tt.fields.Transformer,
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
