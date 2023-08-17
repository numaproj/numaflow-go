// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: pkg/apis/proto/mapstream/v1/mapstream.proto

package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// *
// MapStreamRequest represents a request element.
type MapStreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keys      []string               `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
	Value     []byte                 `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	EventTime *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=event_time,json=eventTime,proto3" json:"event_time,omitempty"`
	Watermark *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=watermark,proto3" json:"watermark,omitempty"`
}

func (x *MapStreamRequest) Reset() {
	*x = MapStreamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_mapstream_v1_mapstream_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MapStreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MapStreamRequest) ProtoMessage() {}

func (x *MapStreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_mapstream_v1_mapstream_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MapStreamRequest.ProtoReflect.Descriptor instead.
func (*MapStreamRequest) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_mapstream_v1_mapstream_proto_rawDescGZIP(), []int{0}
}

func (x *MapStreamRequest) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *MapStreamRequest) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *MapStreamRequest) GetEventTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EventTime
	}
	return nil
}

func (x *MapStreamRequest) GetWatermark() *timestamppb.Timestamp {
	if x != nil {
		return x.Watermark
	}
	return nil
}

// *
// MapStreamResponse represents a response element.
type MapStreamResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Results []*MapStreamResponse_Result `protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty"`
}

func (x *MapStreamResponse) Reset() {
	*x = MapStreamResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_mapstream_v1_mapstream_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MapStreamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MapStreamResponse) ProtoMessage() {}

func (x *MapStreamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_mapstream_v1_mapstream_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MapStreamResponse.ProtoReflect.Descriptor instead.
func (*MapStreamResponse) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_mapstream_v1_mapstream_proto_rawDescGZIP(), []int{1}
}

func (x *MapStreamResponse) GetResults() []*MapStreamResponse_Result {
	if x != nil {
		return x.Results
	}
	return nil
}

// *
// ReadyResponse is the health check result.
type ReadyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ready bool `protobuf:"varint,1,opt,name=ready,proto3" json:"ready,omitempty"`
}

func (x *ReadyResponse) Reset() {
	*x = ReadyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_mapstream_v1_mapstream_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadyResponse) ProtoMessage() {}

func (x *ReadyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_mapstream_v1_mapstream_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadyResponse.ProtoReflect.Descriptor instead.
func (*ReadyResponse) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_mapstream_v1_mapstream_proto_rawDescGZIP(), []int{2}
}

func (x *ReadyResponse) GetReady() bool {
	if x != nil {
		return x.Ready
	}
	return false
}

type MapStreamResponse_Result struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keys  []string `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
	Value []byte   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Tags  []string `protobuf:"bytes,3,rep,name=tags,proto3" json:"tags,omitempty"`
}

func (x *MapStreamResponse_Result) Reset() {
	*x = MapStreamResponse_Result{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_mapstream_v1_mapstream_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MapStreamResponse_Result) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MapStreamResponse_Result) ProtoMessage() {}

func (x *MapStreamResponse_Result) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_mapstream_v1_mapstream_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MapStreamResponse_Result.ProtoReflect.Descriptor instead.
func (*MapStreamResponse_Result) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_mapstream_v1_mapstream_proto_rawDescGZIP(), []int{1, 0}
}

func (x *MapStreamResponse_Result) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *MapStreamResponse_Result) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *MapStreamResponse_Result) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

var File_pkg_apis_proto_mapstream_v1_mapstream_proto protoreflect.FileDescriptor

var file_pkg_apis_proto_mapstream_v1_mapstream_proto_rawDesc = []byte{
	0x0a, 0x2b, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x6d, 0x61, 0x70, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x61,
	0x70, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x6d,
	0x61, 0x70, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70,
	0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb1, 0x01, 0x0a, 0x10, 0x4d, 0x61,
	0x70, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65,
	0x79, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x38, 0x0a, 0x09, 0x77, 0x61, 0x74, 0x65, 0x72, 0x6d, 0x61, 0x72, 0x6b,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x09, 0x77, 0x61, 0x74, 0x65, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x22, 0x9d, 0x01,
	0x0a, 0x11, 0x4d, 0x61, 0x70, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x40, 0x0a, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x6d, 0x61, 0x70, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x61, 0x70, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x07, 0x72, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x73, 0x1a, 0x46, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6b,
	0x65, 0x79, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x61, 0x67,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x22, 0x25, 0x0a,
	0x0d, 0x52, 0x65, 0x61, 0x64, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x72, 0x65, 0x61, 0x64, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x72,
	0x65, 0x61, 0x64, 0x79, 0x32, 0xa4, 0x01, 0x0a, 0x09, 0x4d, 0x61, 0x70, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x12, 0x57, 0x0a, 0x0b, 0x4d, 0x61, 0x70, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x46,
	0x6e, 0x12, 0x1e, 0x2e, 0x6d, 0x61, 0x70, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x76, 0x31,
	0x2e, 0x4d, 0x61, 0x70, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x26, 0x2e, 0x6d, 0x61, 0x70, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x76, 0x31,
	0x2e, 0x4d, 0x61, 0x70, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x30, 0x01, 0x12, 0x3e, 0x0a, 0x07, 0x49,
	0x73, 0x52, 0x65, 0x61, 0x64, 0x79, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1b,
	0x2e, 0x6d, 0x61, 0x70, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65,
	0x61, 0x64, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x3d, 0x5a, 0x3b, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x75, 0x6d, 0x61, 0x70, 0x72,
	0x6f, 0x6a, 0x2f, 0x6e, 0x75, 0x6d, 0x61, 0x66, 0x6c, 0x6f, 0x77, 0x2d, 0x67, 0x6f, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x61,
	0x70, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_pkg_apis_proto_mapstream_v1_mapstream_proto_rawDescOnce sync.Once
	file_pkg_apis_proto_mapstream_v1_mapstream_proto_rawDescData = file_pkg_apis_proto_mapstream_v1_mapstream_proto_rawDesc
)

func file_pkg_apis_proto_mapstream_v1_mapstream_proto_rawDescGZIP() []byte {
	file_pkg_apis_proto_mapstream_v1_mapstream_proto_rawDescOnce.Do(func() {
		file_pkg_apis_proto_mapstream_v1_mapstream_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_apis_proto_mapstream_v1_mapstream_proto_rawDescData)
	})
	return file_pkg_apis_proto_mapstream_v1_mapstream_proto_rawDescData
}

var file_pkg_apis_proto_mapstream_v1_mapstream_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_pkg_apis_proto_mapstream_v1_mapstream_proto_goTypes = []interface{}{
	(*MapStreamRequest)(nil),         // 0: mapstream.v1.MapStreamRequest
	(*MapStreamResponse)(nil),        // 1: mapstream.v1.MapStreamResponse
	(*ReadyResponse)(nil),            // 2: mapstream.v1.ReadyResponse
	(*MapStreamResponse_Result)(nil), // 3: mapstream.v1.MapStreamResponse.Result
	(*timestamppb.Timestamp)(nil),    // 4: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),            // 5: google.protobuf.Empty
}
var file_pkg_apis_proto_mapstream_v1_mapstream_proto_depIdxs = []int32{
	4, // 0: mapstream.v1.MapStreamRequest.event_time:type_name -> google.protobuf.Timestamp
	4, // 1: mapstream.v1.MapStreamRequest.watermark:type_name -> google.protobuf.Timestamp
	3, // 2: mapstream.v1.MapStreamResponse.results:type_name -> mapstream.v1.MapStreamResponse.Result
	0, // 3: mapstream.v1.MapStream.MapStreamFn:input_type -> mapstream.v1.MapStreamRequest
	5, // 4: mapstream.v1.MapStream.IsReady:input_type -> google.protobuf.Empty
	3, // 5: mapstream.v1.MapStream.MapStreamFn:output_type -> mapstream.v1.MapStreamResponse.Result
	2, // 6: mapstream.v1.MapStream.IsReady:output_type -> mapstream.v1.ReadyResponse
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_pkg_apis_proto_mapstream_v1_mapstream_proto_init() }
func file_pkg_apis_proto_mapstream_v1_mapstream_proto_init() {
	if File_pkg_apis_proto_mapstream_v1_mapstream_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_apis_proto_mapstream_v1_mapstream_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MapStreamRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_apis_proto_mapstream_v1_mapstream_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MapStreamResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_apis_proto_mapstream_v1_mapstream_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadyResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_apis_proto_mapstream_v1_mapstream_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MapStreamResponse_Result); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_apis_proto_mapstream_v1_mapstream_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_apis_proto_mapstream_v1_mapstream_proto_goTypes,
		DependencyIndexes: file_pkg_apis_proto_mapstream_v1_mapstream_proto_depIdxs,
		MessageInfos:      file_pkg_apis_proto_mapstream_v1_mapstream_proto_msgTypes,
	}.Build()
	File_pkg_apis_proto_mapstream_v1_mapstream_proto = out.File
	file_pkg_apis_proto_mapstream_v1_mapstream_proto_rawDesc = nil
	file_pkg_apis_proto_mapstream_v1_mapstream_proto_goTypes = nil
	file_pkg_apis_proto_mapstream_v1_mapstream_proto_depIdxs = nil
}