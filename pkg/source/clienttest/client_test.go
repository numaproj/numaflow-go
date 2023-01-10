package clienttest

import (
	"context"
	"fmt"
	"github.com/numaproj/numaflow-go/pkg/apis/proto/source/v1/sourcemock"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	sourcepb "github.com/numaproj/numaflow-go/pkg/apis/proto/source/v1"
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

	mockClient := sourcemock.NewMockUserDefinedSourceTransformerClient(ctrl)
	mockClient.EXPECT().IsReady(gomock.Any(), gomock.Any()).Return(&sourcepb.ReadyResponse{Ready: true}, nil)
	mockClient.EXPECT().IsReady(gomock.Any(), gomock.Any()).Return(&sourcepb.ReadyResponse{Ready: false}, fmt.Errorf("mock connection refused"))

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

func TestTransformFn(t *testing.T) {
	var ctx = context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := sourcemock.NewMockUserDefinedSourceTransformerClient(ctrl)
	testDatum := &sourcepb.Datum{
		Key:       "test_success_key",
		Value:     []byte(`forward_message`),
		EventTime: &sourcepb.EventTime{EventTime: timestamppb.New(time.Unix(1661169600, 0))},
		Watermark: &sourcepb.Watermark{Watermark: timestamppb.New(time.Time{})},
	}
	mockClient.EXPECT().TransformFn(gomock.Any(), &rpcMsg{msg: testDatum}).Return(&sourcepb.DatumList{
		Elements: []*sourcepb.Datum{
			testDatum,
		},
	}, nil)
	mockClient.EXPECT().TransformFn(gomock.Any(), &rpcMsg{msg: testDatum}).Return(&sourcepb.DatumList{
		Elements: []*sourcepb.Datum{
			nil,
		},
	}, fmt.Errorf("mock TransformFn error"))
	testClient, err := New(mockClient)
	assert.NoError(t, err)
	reflect.DeepEqual(testClient, &client{
		grpcClt: mockClient,
	})
	got, err := testClient.TransformFn(ctx, testDatum)
	reflect.DeepEqual(got, testDatum)
	assert.NoError(t, err)

	got, err = testClient.TransformFn(ctx, testDatum)
	assert.Nil(t, got)
	assert.EqualError(t, err, "failed to execute c.grpcClt.TransformFn(): mock TransformFn error")
}
