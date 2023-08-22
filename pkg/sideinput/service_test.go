package sideinput

import (
	"context"
	"reflect"
	"testing"

	"google.golang.org/protobuf/types/known/emptypb"

	sideinputpb "github.com/numaproj/numaflow-go/pkg/apis/proto/sideinput/v1"
)

func TestService_RetrieveSideInputFn(t *testing.T) {
	type args struct {
		ctx context.Context
		in1 *emptypb.Empty
	}

	tests := []struct {
		name      string
		retriever RetrieverSideInput
		args      args
		want      *sideinputpb.SideInputResponse
		wantErr   bool
	}{
		{
			name: "sideinput_retrieve_msg",
			retriever: RetrieverFunc(func(ctx context.Context) MessageSI {
				return NewMessageSI([]byte(`test`))
			}),
			args: args{
				ctx: context.Background(),
				in1: &emptypb.Empty{},
			},
			want: &sideinputpb.SideInputResponse{
				Value: []byte(`test`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &Service{
				Retriever: tt.retriever,
			}
			// here's a trick for testing:
			// because we are not using gRPC, we directly set a new incoming ctx
			// instead of the regular outgoing context in the real gRPC connection.
			ctx := context.Background()
			got, err := fs.RetrieveSideInput(ctx, tt.args.in1)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveSideInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapTFn() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_IsReady(t *testing.T) {
	type args struct {
		in0 context.Context
		in1 *emptypb.Empty
	}
	tests := []struct {
		name                     string
		retrieveSideInputHandler RetrieverSideInput
		args                     args
		want                     *sideinputpb.ReadyResponse
		wantErr                  bool
	}{
		{
			name: "is_ready",
			args: args{
				in0: nil,
				in1: &emptypb.Empty{},
			},
			want: &sideinputpb.ReadyResponse{
				Ready: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &Service{
				Retriever: tt.retrieveSideInputHandler,
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
