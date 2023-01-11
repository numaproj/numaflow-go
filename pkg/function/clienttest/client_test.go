package clienttest

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	"github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1/funcmock"
	"github.com/stretchr/testify/assert"
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

	testClient, err := New(mockClient)
	assert.NoError(t, err)
	reflect.DeepEqual(testClient, &client{
		grpcClt: mockClient,
	})

	ready, err := testClient.IsReady(ctx, &emptypb.Empty{})
	assert.True(t, ready)
	assert.NoError(t, err)

	ready, err = testClient.IsReady(ctx, &emptypb.Empty{})
	assert.False(t, ready)
	assert.EqualError(t, err, "mock connection refused")
}

func TestMapFn(t *testing.T) {
	var ctx = context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := funcmock.NewMockUserDefinedFunctionClient(ctrl)
	testDatum := &functionpb.Datum{
		Key:       "test_success_key",
		Value:     []byte(`forward_message`),
		EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Unix(1661169600, 0))},
		Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
	}
	mockClient.EXPECT().MapFn(gomock.Any(), &rpcMsg{msg: testDatum}).Return(&functionpb.DatumList{
		Elements: []*functionpb.Datum{
			testDatum,
		},
	}, nil)
	mockClient.EXPECT().MapFn(gomock.Any(), &rpcMsg{msg: testDatum}).Return(&functionpb.DatumList{
		Elements: []*functionpb.Datum{
			nil,
		},
	}, fmt.Errorf("mock MapFn error"))
	testClient, err := New(mockClient)
	assert.NoError(t, err)
	reflect.DeepEqual(testClient, &client{
		grpcClt: mockClient,
	})
	got, err := testClient.MapFn(ctx, testDatum)
	reflect.DeepEqual(got, testDatum)
	assert.NoError(t, err)

	got, err = testClient.MapFn(ctx, testDatum)
	assert.Nil(t, got)
	assert.EqualError(t, err, "failed to execute c.grpcClt.MapFn(): mock MapFn error")
}

func TestMapTFn(t *testing.T) {
	var ctx = context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := funcmock.NewMockUserDefinedFunctionClient(ctrl)
	testDatum := &functionpb.Datum{
		Key:       "test_success_key",
		Value:     []byte(`forward_message`),
		EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Unix(1661169600, 0))},
		Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
	}
	mockClient.EXPECT().MapTFn(gomock.Any(), &rpcMsg{msg: testDatum}).Return(&functionpb.DatumList{
		Elements: []*functionpb.Datum{
			testDatum,
		},
	}, nil)
	mockClient.EXPECT().MapTFn(gomock.Any(), &rpcMsg{msg: testDatum}).Return(&functionpb.DatumList{
		Elements: []*functionpb.Datum{
			nil,
		},
	}, fmt.Errorf("mock MapTFn error"))
	testClient, err := New(mockClient)
	assert.NoError(t, err)
	reflect.DeepEqual(testClient, &client{
		grpcClt: mockClient,
	})
	got, err := testClient.MapTFn(ctx, testDatum)
	reflect.DeepEqual(got, testDatum)
	assert.NoError(t, err)

	got, err = testClient.MapTFn(ctx, testDatum)
	assert.Nil(t, got)
	assert.EqualError(t, err, "failed to execute c.grpcClt.MapTFn(): mock MapTFn error")
}
