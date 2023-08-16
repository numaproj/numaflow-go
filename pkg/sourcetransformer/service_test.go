package sourcetransformer

import (
	"context"
	"reflect"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/sourcetransformer/v1"
)

func TestService_MapTFn(t *testing.T) {
	type args struct {
		ctx context.Context
		d   *v1.SourceTransformerRequest
	}

	testTime := time.Date(2021, 8, 15, 14, 30, 45, 100, time.Local)
	tests := []struct {
		name    string
		handler SourceTransformer
		args    args
		want    *v1.SourceTransformerResponseList
		wantErr bool
	}{
		{
			name: "mapT_fn_forward_msg",
			handler: SourceTransformFunc(func(ctx context.Context, keys []string, datum Datum) Messages {
				msg := datum.Value()
				return MessagesBuilder().Append(NewMessage(msg, testTime).WithKeys([]string{keys[0] + "_test"}))
			}),
			args: args{
				ctx: context.Background(),
				d: &v1.SourceTransformerRequest{
					Keys:      []string{"client"},
					Value:     []byte(`test`),
					EventTime: &v1.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &v1.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			want: &v1.SourceTransformerResponseList{
				Elements: []*v1.SourceTransformerResponse{
					{
						EventTime: &v1.EventTime{EventTime: timestamppb.New(testTime)},
						Keys:      []string{"client_test"},
						Value:     []byte(`test`),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "mapT_fn_forward_msg_forward_to_all",
			handler: SourceTransformFunc(func(ctx context.Context, keys []string, datum Datum) Messages {
				msg := datum.Value()
				return MessagesBuilder().Append(NewMessage(msg, testTime))
			}),
			args: args{
				ctx: context.Background(),
				d: &v1.SourceTransformerRequest{
					Keys:      []string{"client"},
					Value:     []byte(`test`),
					EventTime: &v1.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &v1.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			want: &v1.SourceTransformerResponseList{
				Elements: []*v1.SourceTransformerResponse{
					{
						EventTime: &v1.EventTime{EventTime: timestamppb.New(testTime)},
						Value:     []byte(`test`),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "mapT_fn_forward_msg_drop_msg",
			handler: SourceTransformFunc(func(ctx context.Context, keys []string, datum Datum) Messages {
				return MessagesBuilder().Append(MessageToDrop())
			}),
			args: args{
				ctx: context.Background(),
				d: &v1.SourceTransformerRequest{
					Keys:      []string{"client"},
					Value:     []byte(`test`),
					EventTime: &v1.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &v1.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			want: &v1.SourceTransformerResponseList{
				Elements: []*v1.SourceTransformerResponse{
					{
						EventTime: &v1.EventTime{EventTime: timestamppb.New(time.Time{})},
						Tags:      []string{DROP},
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
				Transformer: tt.handler,
			}
			// here's a trick for testing:
			// because we are not using gRPC, we directly set a new incoming ctx
			// instead of the regular outgoing context in the real gRPC connection.
			ctx := context.Background()
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
