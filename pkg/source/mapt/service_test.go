package mapt

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/numaproj/numaflow-go/pkg/apis/proto/source/transformerfn"
	"github.com/numaproj/numaflow-go/pkg/source"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestService_MapTFn(t *testing.T) {
	type args struct {
		ctx context.Context
		d   *transformerfn.SourceTransformerRequest
	}

	testTime := time.Date(2021, 8, 15, 14, 30, 45, 100, time.Local)
	tests := []struct {
		name    string
		handler source.MapTHandler
		args    args
		want    *transformerfn.SourceTransformerResponseList
		wantErr bool
	}{
		{
			name: "mapT_fn_forward_msg",
			handler: source.MapTFunc(func(ctx context.Context, keys []string, datum source.Datum) source.MessageTs {
				msg := datum.Value()
				return source.MessageTsBuilder().Append(source.NewMessageT(msg, testTime).WithKeys([]string{keys[0] + "_test"}))
			}),
			args: args{
				ctx: context.Background(),
				d: &transformerfn.SourceTransformerRequest{
					Keys:      []string{"client"},
					Value:     []byte(`test`),
					EventTime: &transformerfn.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &transformerfn.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			want: &transformerfn.SourceTransformerResponseList{
				Elements: []*transformerfn.SourceTransformerResponse{
					{
						EventTime: &transformerfn.EventTime{EventTime: timestamppb.New(testTime)},
						Keys:      []string{"client_test"},
						Value:     []byte(`test`),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "mapT_fn_forward_msg_forward_to_all",
			handler: source.MapTFunc(func(ctx context.Context, keys []string, datum source.Datum) source.MessageTs {
				msg := datum.Value()
				return source.MessageTsBuilder().Append(source.NewMessageT(msg, testTime))
			}),
			args: args{
				ctx: context.Background(),
				d: &transformerfn.SourceTransformerRequest{
					Keys:      []string{"client"},
					Value:     []byte(`test`),
					EventTime: &transformerfn.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &transformerfn.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			want: &transformerfn.SourceTransformerResponseList{
				Elements: []*transformerfn.SourceTransformerResponse{
					{
						EventTime: &transformerfn.EventTime{EventTime: timestamppb.New(testTime)},
						Value:     []byte(`test`),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "mapT_fn_forward_msg_drop_msg",
			handler: source.MapTFunc(func(ctx context.Context, keys []string, datum source.Datum) source.MessageTs {
				return source.MessageTsBuilder().Append(source.MessageTToDrop())
			}),
			args: args{
				ctx: context.Background(),
				d: &transformerfn.SourceTransformerRequest{
					Keys:      []string{"client"},
					Value:     []byte(`test`),
					EventTime: &transformerfn.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &transformerfn.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			want: &transformerfn.SourceTransformerResponseList{
				Elements: []*transformerfn.SourceTransformerResponse{
					{
						EventTime: &transformerfn.EventTime{EventTime: timestamppb.New(time.Time{})},
						Tags:      []string{source.DROP},
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
				MapperT: tt.handler,
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
