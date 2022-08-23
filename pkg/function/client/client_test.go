package client

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1/funcmock"
	"google.golang.org/grpc"
)

func TestNewClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := funcmock.NewMockUserDefinedFunctionClient(ctrl)

	type args struct {
		inputOptions []Option
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		{
			name: "new_client_with_mock",
			args: args{
				inputOptions: []Option{WithMockGRPCClient(mockClient)},
			},
			want: &Client{
				conn:    &grpc.ClientConn{},
				grpcClt: mockClient,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.inputOptions...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() got = %v, want %v", got, tt.want)
			}
		})
	}
}
