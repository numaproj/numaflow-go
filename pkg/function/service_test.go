package function

import (
	"context"
	"reflect"
	"testing"
	"time"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestService_DoFn(t *testing.T) {
	type fields struct {
		Mapper  MapHandler
		Reducer ReduceHandler
	}
	type args struct {
		ctx context.Context
		d   *functionpb.Datum
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *functionpb.DatumList
		wantErr bool
	}{
		{
			name: "do_fn_forward_msg",
			fields: fields{
				Mapper: DoFunc(func(ctx context.Context, key string, msg []byte) (Messages, error) {
					return MessagesBuilder().Append(MessageTo(key+"_test", msg)), nil
				}),
				Reducer: nil,
			},
			args: args{
				ctx: nil,
				d: &functionpb.Datum{
					Key:   "client",
					Value: []byte(`test`),
					EventTime: &functionpb.EventTime{
						EventTime: timestamppb.New(time.Unix(1661169600, 0)),
					},
					IntervalWindow: &functionpb.IntervalWindow{
						StartTime: timestamppb.New(time.Unix(1661169600, 0)),
						EndTime:   timestamppb.New(time.Unix(1661169660, 0)),
					},
					PaneInfo: &functionpb.PaneInfo{
						// TODO: need to update once we've finalized the datum data type
						Watermark: timestamppb.New(time.Time{}),
					},
				},
			},
			want: &functionpb.DatumList{
				Elements: []*functionpb.Datum{
					{
						Key:   "client_test",
						Value: []byte(`test`),
						EventTime: &functionpb.EventTime{
							EventTime: timestamppb.New(time.Unix(1661169600, 0)),
						},
						IntervalWindow: &functionpb.IntervalWindow{
							StartTime: timestamppb.New(time.Unix(1661169600, 0)),
							EndTime:   timestamppb.New(time.Unix(1661169660, 0)),
						},
						PaneInfo: &functionpb.PaneInfo{
							// TODO: need to update once we've finalized the datum data type
							Watermark: timestamppb.New(time.Time{}),
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &Service{
				Mapper:  tt.fields.Mapper,
				Reducer: tt.fields.Reducer,
			}
			got, err := fs.DoFn(tt.args.ctx, tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoFn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoFn() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_IsReady(t *testing.T) {
	type fields struct {
		Mapper  MapHandler
		Reducer ReduceHandler
	}
	type args struct {
		in0 context.Context
		in1 *emptypb.Empty
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *functionpb.ReadyResponse
		wantErr bool
	}{
		{
			name: "is_ready",
			fields: fields{
				Mapper:  nil,
				Reducer: nil,
			},
			args: args{
				in0: nil,
				in1: &emptypb.Empty{},
			},
			want: &functionpb.ReadyResponse{
				Ready: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &Service{
				Mapper:  tt.fields.Mapper,
				Reducer: tt.fields.Reducer,
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

func TestService_ReduceFn(t *testing.T) {
	type fields struct {
		UnimplementedUserDefinedFunctionServer functionpb.UnimplementedUserDefinedFunctionServer
		Mapper                                 MapHandler
		Reducer                                ReduceHandler
	}
	type args struct {
		fnServer functionpb.UserDefinedFunction_ReduceFnServer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &Service{
				UnimplementedUserDefinedFunctionServer: tt.fields.UnimplementedUserDefinedFunctionServer,
				Mapper:                                 tt.fields.Mapper,
				Reducer:                                tt.fields.Reducer,
			}
			if err := fs.ReduceFn(tt.args.fnServer); (err != nil) != tt.wantErr {
				t.Errorf("ReduceFn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
