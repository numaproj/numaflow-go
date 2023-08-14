package mapsvc

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/numaproj/numaflow-go/pkg/apis/proto/function/mapfn"
	"github.com/numaproj/numaflow-go/pkg/function"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestService_MapFn(t *testing.T) {
	type args struct {
		ctx context.Context
		d   *mapfn.MapRequest
	}
	tests := []struct {
		name    string
		handler function.MapHandler
		args    args
		want    *mapfn.MapResponseList
		wantErr bool
	}{
		{
			name: "map_fn_forward_msg",
			handler: function.MapFunc(func(ctx context.Context, keys []string, datum function.Datum) function.Messages {
				msg := datum.Value()
				return function.MessagesBuilder().Append(function.NewMessage(msg).WithKeys([]string{keys[0] + "_test"}))
			}),
			args: args{
				ctx: context.Background(),
				d: &mapfn.MapRequest{
					Keys:      []string{"client"},
					Value:     []byte(`test`),
					EventTime: &mapfn.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &mapfn.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			want: &mapfn.MapResponseList{
				Elements: []*mapfn.MapResponse{
					{
						Keys:  []string{"client_test"},
						Value: []byte(`test`),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "map_fn_forward_msg_forward_to_all",
			handler: function.MapFunc(func(ctx context.Context, keys []string, datum function.Datum) function.Messages {
				msg := datum.Value()
				return function.MessagesBuilder().Append(function.NewMessage(msg))
			}),
			args: args{
				ctx: context.Background(),
				d: &mapfn.MapRequest{
					Keys:      []string{"client"},
					Value:     []byte(`test`),
					EventTime: &mapfn.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &mapfn.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			want: &mapfn.MapResponseList{
				Elements: []*mapfn.MapResponse{
					{
						Value: []byte(`test`),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "map_fn_forward_msg_drop_msg",
			handler: function.MapFunc(func(ctx context.Context, keys []string, datum function.Datum) function.Messages {
				return function.MessagesBuilder().Append(function.MessageToDrop())
			}),
			args: args{
				ctx: context.Background(),
				d: &mapfn.MapRequest{
					Keys:      []string{"client"},
					Value:     []byte(`test`),
					EventTime: &mapfn.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &mapfn.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			want: &mapfn.MapResponseList{
				Elements: []*mapfn.MapResponse{
					{
						Tags:  []string{function.DROP},
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
				Mapper: tt.handler,
			}
			// here's a trick for testing:
			// because we are not using gRPC, we directly set a new incoming ctx
			// instead of the regular outgoing context in the real gRPC connection.
			ctx := context.Background()
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
