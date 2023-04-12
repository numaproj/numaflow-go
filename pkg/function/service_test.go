package function

import (
	"context"
	"io"
	"reflect"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	"google.golang.org/grpc"
	grpcmd "google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type fields struct {
	mapper  MapHandler
	mapperT MapTHandler
	reducer ReduceHandler
}

type UserDefinedFunctionReduceFnServerTest struct {
	ctx      context.Context
	inputCh  chan *functionpb.DatumRequest
	outputCh chan *functionpb.DatumResponseList
	grpc.ServerStream
}

func NewUserDefinedFunctionReduceFnServerTest(ctx context.Context,
	inputCh chan *functionpb.DatumRequest,
	outputCh chan *functionpb.DatumResponseList) *UserDefinedFunctionReduceFnServerTest {
	return &UserDefinedFunctionReduceFnServerTest{
		ctx:      ctx,
		inputCh:  inputCh,
		outputCh: outputCh,
	}
}

func (u *UserDefinedFunctionReduceFnServerTest) Send(list *functionpb.DatumResponseList) error {
	u.outputCh <- list
	return nil
}

func (u *UserDefinedFunctionReduceFnServerTest) Recv() (*functionpb.DatumRequest, error) {
	val, ok := <-u.inputCh
	if !ok {
		return val, io.EOF
	}
	return val, nil
}

func (u *UserDefinedFunctionReduceFnServerTest) Context() context.Context {
	return u.ctx
}

func TestService_MapFn(t *testing.T) {
	type args struct {
		ctx context.Context
		d   *functionpb.DatumRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *functionpb.DatumResponseList
		wantErr bool
	}{
		{
			name: "map_fn_forward_msg",
			fields: fields{
				mapper: MapFunc(func(ctx context.Context, keys []string, datum Datum) Messages {
					msg := datum.Value()
					return MessagesBuilder().Append(NewMessage(msg).WithKeys([]string{keys[0] + "_test"}))
				}),
			},
			args: args{
				ctx: context.Background(),
				d: &functionpb.DatumRequest{
					Keys:      []string{"client"},
					Value:     []byte(`test`),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			want: &functionpb.DatumResponseList{
				Elements: []*functionpb.DatumResponse{
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
			fields: fields{
				mapper: MapFunc(func(ctx context.Context, keys []string, datum Datum) Messages {
					msg := datum.Value()
					return MessagesBuilder().Append(NewMessage(msg))
				}),
			},
			args: args{
				ctx: context.Background(),
				d: &functionpb.DatumRequest{
					Keys:      []string{"client"},
					Value:     []byte(`test`),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			want: &functionpb.DatumResponseList{
				Elements: []*functionpb.DatumResponse{
					{
						Value: []byte(`test`),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "map_fn_forward_msg_drop_msg",
			fields: fields{
				mapper: MapFunc(func(ctx context.Context, keys []string, datum Datum) Messages {
					return MessagesBuilder().Append(MessageToDrop())
				}),
			},
			args: args{
				ctx: context.Background(),
				d: &functionpb.DatumRequest{
					Keys:      []string{"client"},
					Value:     []byte(`test`),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			want: &functionpb.DatumResponseList{
				Elements: []*functionpb.DatumResponse{
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
				Mapper:  tt.fields.mapper,
				MapperT: tt.fields.mapperT,
				Reducer: tt.fields.reducer,
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

func TestService_MapTFn(t *testing.T) {
	type args struct {
		ctx context.Context
		d   *functionpb.DatumRequest
	}

	testTime := time.Date(2021, 8, 15, 14, 30, 45, 100, time.Local)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *functionpb.DatumResponseList
		wantErr bool
	}{
		{
			name: "mapT_fn_forward_msg",
			fields: fields{
				mapperT: MapTFunc(func(ctx context.Context, keys []string, datum Datum) MessageTs {
					msg := datum.Value()
					return MessageTsBuilder().Append(NewMessageT(testTime, msg).WithKeys([]string{keys[0] + "_test"}))
				}),
			},
			args: args{
				ctx: context.Background(),
				d: &functionpb.DatumRequest{
					Keys:      []string{"client"},
					Value:     []byte(`test`),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			want: &functionpb.DatumResponseList{
				Elements: []*functionpb.DatumResponse{
					{
						EventTime: &functionpb.EventTime{EventTime: timestamppb.New(testTime)},
						Keys:      []string{"client_test"},
						Value:     []byte(`test`),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "mapT_fn_forward_msg_forward_to_all",
			fields: fields{
				mapperT: MapTFunc(func(ctx context.Context, keys []string, datum Datum) MessageTs {
					msg := datum.Value()
					return MessageTsBuilder().Append(NewMessageT(testTime, msg))
				}),
			},
			args: args{
				ctx: context.Background(),
				d: &functionpb.DatumRequest{
					Keys:      []string{"client"},
					Value:     []byte(`test`),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			want: &functionpb.DatumResponseList{
				Elements: []*functionpb.DatumResponse{
					{
						EventTime: &functionpb.EventTime{EventTime: timestamppb.New(testTime)},
						Value:     []byte(`test`),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "mapT_fn_forward_msg_drop_msg",
			fields: fields{
				mapperT: MapTFunc(func(ctx context.Context, keys []string, datum Datum) MessageTs {
					return MessageTsBuilder().Append(MessageTToDrop())
				}),
			},
			args: args{
				ctx: context.Background(),
				d: &functionpb.DatumRequest{
					Keys:      []string{"client"},
					Value:     []byte(`test`),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			want: &functionpb.DatumResponseList{
				Elements: []*functionpb.DatumResponse{
					{
						EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
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
				Mapper:  tt.fields.mapper,
				MapperT: tt.fields.mapperT,
				Reducer: tt.fields.reducer,
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

func TestService_ReduceFn(t *testing.T) {

	tests := []struct {
		name        string
		fields      fields
		input       []*functionpb.DatumRequest
		expected    *functionpb.DatumResponseList
		expectedErr bool
	}{
		{
			name: "reduce_fn_forward_msg_same_keys",
			fields: fields{
				reducer: ReduceFunc(func(ctx context.Context, keys []string, rch <-chan Datum, md Metadata) Messages {
					sum := 0
					for val := range rch {
						msgVal, _ := strconv.Atoi(string(val.Value()))
						sum += msgVal
					}
					return MessagesBuilder().Append(NewMessage([]byte(strconv.Itoa(sum))).WithKeys([]string{keys[0] + "_test"}))
				}),
			},
			input: []*functionpb.DatumRequest{
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(10)),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(20)),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(30)),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			expected: &functionpb.DatumResponseList{
				Elements: []*functionpb.DatumResponse{
					{
						Keys:  []string{"client_test"},
						Value: []byte(strconv.Itoa(60)),
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "reduce_fn_forward_msg_multiple_keys",
			fields: fields{
				reducer: ReduceFunc(func(ctx context.Context, keys []string, rch <-chan Datum, md Metadata) Messages {
					sum := 0
					for val := range rch {
						msgVal, _ := strconv.Atoi(string(val.Value()))
						sum += msgVal
					}
					return MessagesBuilder().Append(NewMessage([]byte(strconv.Itoa(sum))).WithKeys([]string{keys[0] + "_test"}))
				}),
			},
			input: []*functionpb.DatumRequest{
				{
					Keys:      []string{"client1"},
					Value:     []byte(strconv.Itoa(10)),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Keys:      []string{"client2"},
					Value:     []byte(strconv.Itoa(20)),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Keys:      []string{"client3"},
					Value:     []byte(strconv.Itoa(30)),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Keys:      []string{"client1"},
					Value:     []byte(strconv.Itoa(10)),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Keys:      []string{"client2"},
					Value:     []byte(strconv.Itoa(20)),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Keys:      []string{"client3"},
					Value:     []byte(strconv.Itoa(30)),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			expected: &functionpb.DatumResponseList{
				Elements: []*functionpb.DatumResponse{
					{
						Keys:  []string{"client1_test"},
						Value: []byte(strconv.Itoa(20)),
					},
					{
						Keys:  []string{"client2_test"},
						Value: []byte(strconv.Itoa(40)),
					},
					{
						Keys:  []string{"client3_test"},
						Value: []byte(strconv.Itoa(60)),
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "reduce_fn_forward_msg_forward_to_all",
			fields: fields{
				reducer: ReduceFunc(func(ctx context.Context, keys []string, rch <-chan Datum, md Metadata) Messages {
					sum := 0
					for val := range rch {
						msgVal, _ := strconv.Atoi(string(val.Value()))
						sum += msgVal
					}
					return MessagesBuilder().Append(NewMessage([]byte(strconv.Itoa(sum))))
				}),
			},
			input: []*functionpb.DatumRequest{
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(10)),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(20)),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(30)),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			expected: &functionpb.DatumResponseList{
				Elements: []*functionpb.DatumResponse{
					{
						Value: []byte(strconv.Itoa(60)),
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "reduce_fn_forward_msg_drop_msg",
			fields: fields{
				reducer: ReduceFunc(func(ctx context.Context, keys []string, rch <-chan Datum, md Metadata) Messages {
					sum := 0
					for val := range rch {
						msgVal, _ := strconv.Atoi(string(val.Value()))
						sum += msgVal
					}
					return MessagesBuilder().Append(MessageToDrop())
				}),
			},
			input: []*functionpb.DatumRequest{
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(10)),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(20)),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
				{
					Keys:      []string{"client"},
					Value:     []byte(strconv.Itoa(30)),
					EventTime: &functionpb.EventTime{EventTime: timestamppb.New(time.Time{})},
					Watermark: &functionpb.Watermark{Watermark: timestamppb.New(time.Time{})},
				},
			},
			expected: &functionpb.DatumResponseList{
				Elements: []*functionpb.DatumResponse{
					{
						Tags:  []string{DROP},
						Value: []byte{},
					},
				},
			},
			expectedErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &Service{
				Mapper:  tt.fields.mapper,
				MapperT: tt.fields.mapperT,
				Reducer: tt.fields.reducer,
			}
			// here's a trick for testing:
			// because we are not using gRPC, we directly set a new incoming ctx
			// instead of the regular outgoing context in the real gRPC connection.
			ctx := grpcmd.NewIncomingContext(context.Background(), grpcmd.New(map[string]string{WinStartTime: "60000", WinEndTime: "120000"}))

			inputCh := make(chan *functionpb.DatumRequest)
			outputCh := make(chan *functionpb.DatumResponseList)
			result := &functionpb.DatumResponseList{}

			udfReduceFnStream := NewUserDefinedFunctionReduceFnServerTest(ctx, inputCh, outputCh)

			var wg sync.WaitGroup
			var err error

			wg.Add(1)
			go func() {
				defer wg.Done()
				err = fs.ReduceFn(udfReduceFnStream)
				close(outputCh)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				for msg := range outputCh {
					result.Elements = append(result.Elements, msg.Elements...)
				}
			}()

			for _, val := range tt.input {
				udfReduceFnStream.inputCh <- val
			}
			close(udfReduceFnStream.inputCh)
			wg.Wait()

			if (err != nil) != tt.expectedErr {
				t.Errorf("ReduceFn() error = %v, wantErr %v", err, tt.expectedErr)
				return
			}

			//sort and compare, since order of the output doesn't matter
			sort.Slice(result.Elements, func(i, j int) bool {
				return string(result.Elements[i].Value) < string(result.Elements[j].Value)
			})

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ReduceFn() got = %v, want %v", result, tt.expected)
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
		name    string
		fields  fields
		args    args
		want    *functionpb.ReadyResponse
		wantErr bool
	}{
		{
			name: "is_ready",
			fields: fields{
				mapper:  nil,
				mapperT: nil,
				reducer: nil,
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
				Mapper:  tt.fields.mapper,
				MapperT: tt.fields.mapperT,
				Reducer: tt.fields.reducer,
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
