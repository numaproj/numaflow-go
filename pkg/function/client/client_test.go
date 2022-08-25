package client

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	"github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1/funcmock"
	"github.com/numaproj/numaflow-go/pkg/function/clienttest"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type rpcMsg struct {
	msg proto.Message
}

func (r *rpcMsg) Matches(msg interface{}) bool {
	m, ok := msg.(proto.Message)
	if !ok {
		return false
	}
	return proto.Equal(m, r.msg)
}

func (r *rpcMsg) String() string {
	return fmt.Sprintf("is %s", r.msg)
}

func TestIsReady(t *testing.T) {
	var ctx = context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := funcmock.NewMockUserDefinedFunctionClient(ctrl)
	mockClient.EXPECT().IsReady(gomock.Any(), gomock.Any()).Return(&functionpb.ReadyResponse{Ready: true}, nil)
	mockClient.EXPECT().IsReady(gomock.Any(), gomock.Any()).Return(&functionpb.ReadyResponse{Ready: false}, fmt.Errorf("mock connection refused"))

	testClient, err := clienttest.New(mockClient)
	assert.NoError(t, err)
	reflect.DeepEqual(testClient, &client{
		conn:    &grpc.ClientConn{},
		grpcClt: mockClient,
	})

	ready, err := testClient.IsReady(ctx, &emptypb.Empty{})
	assert.True(t, ready)
	assert.NoError(t, err)

	ready, err = testClient.IsReady(ctx, &emptypb.Empty{})
	assert.False(t, ready)
	assert.EqualError(t, err, "mock connection refused")
}

func TestDoFn(t *testing.T) {
	var ctx = context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := funcmock.NewMockUserDefinedFunctionClient(ctrl)
	testDatum := &functionpb.Datum{
		Key:            "test_success_key",
		Value:          []byte(`forward_message`),
		EventTime:      &functionpb.EventTime{EventTime: timestamppb.New(time.Unix(1661169600, 0))},
		IntervalWindow: &functionpb.IntervalWindow{StartTime: timestamppb.New(time.Unix(1661169600, 0)), EndTime: timestamppb.New(time.Unix(1661169660, 0))},
		PaneInfo:       &functionpb.PaneInfo{Watermark: timestamppb.New(time.Time{})},
	}
	mockClient.EXPECT().DoFn(gomock.Any(), &rpcMsg{msg: testDatum}).Return(&functionpb.DatumList{
		Elements: []*functionpb.Datum{
			testDatum,
		},
	}, nil)
	mockClient.EXPECT().DoFn(gomock.Any(), &rpcMsg{msg: testDatum}).Return(&functionpb.DatumList{
		Elements: []*functionpb.Datum{
			nil,
		},
	}, fmt.Errorf("mock DoFn error"))
	testClient, err := clienttest.New(mockClient)
	assert.NoError(t, err)
	reflect.DeepEqual(testClient, &client{
		conn:    &grpc.ClientConn{},
		grpcClt: mockClient,
	})

	got, err := testClient.DoFn(ctx, testDatum)
	reflect.DeepEqual(got, testDatum)
	assert.NoError(t, err)

	got, err = testClient.DoFn(ctx, testDatum)
	assert.Nil(t, got)
	assert.EqualError(t, err, "failed to execute c.grpcClt.DoFn(): mock DoFn error")
}
