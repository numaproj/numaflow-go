package _map

import (
	"context"
	"reflect"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	mappb "github.com/numaproj/numaflow-go/pkg/apis/proto/map/v1"
)

func TestService_MapFn(t *testing.T) {
	type args struct {
		ctx context.Context
		d   *mappb.MapRequest
	}
	tests := []struct {
		name    string
		handler Mapper
		args    args
		want    *mappb.MapResponseList
		wantErr bool
	}{
		{
			name: "map_fn_forward_msg",
			handler: MapperFunc(func(ctx context.Context, keys []string, datum Datum) Messages {
				msg := datum.Value()
				return MessagesBuilder().Append(NewMessage(msg).WithKeys([]string{keys[0] + "_test"}))
			}),
			args: args{
				ctx: context.Background(),
				d: &mappb.MapRequest{
					Keys:      []string{"client"},
					Value:     []byte(`test`),
					EventTime: &mappb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &mappb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			want: &mappb.MapResponseList{
				Elements: []*mappb.MapResponse{
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
			handler: MapperFunc(func(ctx context.Context, keys []string, datum Datum) Messages {
				msg := datum.Value()
				return MessagesBuilder().Append(NewMessage(msg))
			}),
			args: args{
				ctx: context.Background(),
				d: &mappb.MapRequest{
					Keys:      []string{"client"},
					Value:     []byte(`test`),
					EventTime: &mappb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &mappb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			want: &mappb.MapResponseList{
				Elements: []*mappb.MapResponse{
					{
						Value: []byte(`test`),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "map_fn_forward_msg_drop_msg",
			handler: MapperFunc(func(ctx context.Context, keys []string, datum Datum) Messages {
				return MessagesBuilder().Append(MessageToDrop())
			}),
			args: args{
				ctx: context.Background(),
				d: &mappb.MapRequest{
					Keys:      []string{"client"},
					Value:     []byte(`test`),
					EventTime: &mappb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &mappb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			want: &mappb.MapResponseList{
				Elements: []*mappb.MapResponse{
					{
						Tags:  []string{DROP},
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
