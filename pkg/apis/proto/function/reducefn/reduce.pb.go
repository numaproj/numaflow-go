// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: pkg/apis/proto/function/reducefn/reduce.proto

package reducefn

import (
	common "github.com/numaproj/numaflow-go/pkg/apis/proto/common"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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
// ReduceRequest represents a request element.
type ReduceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keys      []string          `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
	Value     []byte            `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	EventTime *common.EventTime `protobuf:"bytes,3,opt,name=event_time,json=eventTime,proto3" json:"event_time,omitempty"`
	Watermark *common.Watermark `protobuf:"bytes,4,opt,name=watermark,proto3" json:"watermark,omitempty"`
}

func (x *ReduceRequest) Reset() {
	*x = ReduceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_function_reducefn_reduce_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReduceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReduceRequest) ProtoMessage() {}

func (x *ReduceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_function_reducefn_reduce_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReduceRequest.ProtoReflect.Descriptor instead.
func (*ReduceRequest) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_function_reducefn_reduce_proto_rawDescGZIP(), []int{0}
}

func (x *ReduceRequest) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *ReduceRequest) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *ReduceRequest) GetEventTime() *common.EventTime {
	if x != nil {
		return x.EventTime
	}
	return nil
}

func (x *ReduceRequest) GetWatermark() *common.Watermark {
	if x != nil {
		return x.Watermark
	}
	return nil
}

// *
// ReduceResponse represents a response element.
type ReduceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keys  []string `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
	Value []byte   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Tags  []string `protobuf:"bytes,3,rep,name=tags,proto3" json:"tags,omitempty"`
}

func (x *ReduceResponse) Reset() {
	*x = ReduceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_function_reducefn_reduce_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReduceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReduceResponse) ProtoMessage() {}

func (x *ReduceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_function_reducefn_reduce_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReduceResponse.ProtoReflect.Descriptor instead.
func (*ReduceResponse) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_function_reducefn_reduce_proto_rawDescGZIP(), []int{1}
}

func (x *ReduceResponse) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *ReduceResponse) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *ReduceResponse) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

// *
// ReduceResponseList represents a list of response elements.
type ReduceResponseList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Elements []*ReduceResponse `protobuf:"bytes,1,rep,name=elements,proto3" json:"elements,omitempty"`
}

func (x *ReduceResponseList) Reset() {
	*x = ReduceResponseList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_function_reducefn_reduce_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReduceResponseList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReduceResponseList) ProtoMessage() {}

func (x *ReduceResponseList) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_function_reducefn_reduce_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReduceResponseList.ProtoReflect.Descriptor instead.
func (*ReduceResponseList) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_function_reducefn_reduce_proto_rawDescGZIP(), []int{2}
}

func (x *ReduceResponseList) GetElements() []*ReduceResponse {
	if x != nil {
		return x.Elements
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
		mi := &file_pkg_apis_proto_function_reducefn_reduce_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadyResponse) ProtoMessage() {}

func (x *ReadyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_function_reducefn_reduce_proto_msgTypes[3]
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
	return file_pkg_apis_proto_function_reducefn_reduce_proto_rawDescGZIP(), []int{3}
}

func (x *ReadyResponse) GetReady() bool {
	if x != nil {
		return x.Ready
	}
	return false
}

var File_pkg_apis_proto_function_reducefn_reduce_proto protoreflect.FileDescriptor

var file_pkg_apis_proto_function_reducefn_reduce_proto_rawDesc = []byte{
	0x0a, 0x2d, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65,
	0x66, 0x6e, 0x2f, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x11, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65,
	0x66, 0x6e, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x22, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x9c, 0x01, 0x0a, 0x0d, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12,
	0x30, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d,
	0x65, 0x12, 0x2f, 0x0a, 0x09, 0x77, 0x61, 0x74, 0x65, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x57, 0x61,
	0x74, 0x65, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x52, 0x09, 0x77, 0x61, 0x74, 0x65, 0x72, 0x6d, 0x61,
	0x72, 0x6b, 0x22, 0x4e, 0x0a, 0x0e, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x74, 0x61, 0x67, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x74, 0x61,
	0x67, 0x73, 0x22, 0x53, 0x0a, 0x12, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x3d, 0x0a, 0x08, 0x65, 0x6c, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x66, 0x75, 0x6e,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x66, 0x6e, 0x2e, 0x52,
	0x65, 0x64, 0x75, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x08, 0x65,
	0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x25, 0x0a, 0x0d, 0x52, 0x65, 0x61, 0x64, 0x79,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x65, 0x61, 0x64,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x72, 0x65, 0x61, 0x64, 0x79, 0x32, 0xa6,
	0x01, 0x0a, 0x06, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x12, 0x57, 0x0a, 0x08, 0x52, 0x65, 0x64,
	0x75, 0x63, 0x65, 0x46, 0x6e, 0x12, 0x20, 0x2e, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x66, 0x6e, 0x2e, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x66, 0x6e, 0x2e, 0x52, 0x65, 0x64, 0x75,
	0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x28, 0x01,
	0x30, 0x01, 0x12, 0x43, 0x0a, 0x07, 0x49, 0x73, 0x52, 0x65, 0x61, 0x64, 0x79, 0x12, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x20, 0x2e, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x66, 0x6e, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x79, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x42, 0x5a, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x75, 0x6d, 0x61, 0x70, 0x72, 0x6f, 0x6a, 0x2f, 0x6e,
	0x75, 0x6d, 0x61, 0x66, 0x6c, 0x6f, 0x77, 0x2d, 0x67, 0x6f, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61,
	0x70, 0x69, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x2f, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x66, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_pkg_apis_proto_function_reducefn_reduce_proto_rawDescOnce sync.Once
	file_pkg_apis_proto_function_reducefn_reduce_proto_rawDescData = file_pkg_apis_proto_function_reducefn_reduce_proto_rawDesc
)

func file_pkg_apis_proto_function_reducefn_reduce_proto_rawDescGZIP() []byte {
	file_pkg_apis_proto_function_reducefn_reduce_proto_rawDescOnce.Do(func() {
		file_pkg_apis_proto_function_reducefn_reduce_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_apis_proto_function_reducefn_reduce_proto_rawDescData)
	})
	return file_pkg_apis_proto_function_reducefn_reduce_proto_rawDescData
}

var file_pkg_apis_proto_function_reducefn_reduce_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_pkg_apis_proto_function_reducefn_reduce_proto_goTypes = []interface{}{
	(*ReduceRequest)(nil),      // 0: function.reducefn.ReduceRequest
	(*ReduceResponse)(nil),     // 1: function.reducefn.ReduceResponse
	(*ReduceResponseList)(nil), // 2: function.reducefn.ReduceResponseList
	(*ReadyResponse)(nil),      // 3: function.reducefn.ReadyResponse
	(*common.EventTime)(nil),   // 4: common.EventTime
	(*common.Watermark)(nil),   // 5: common.Watermark
	(*emptypb.Empty)(nil),      // 6: google.protobuf.Empty
}
var file_pkg_apis_proto_function_reducefn_reduce_proto_depIdxs = []int32{
	4, // 0: function.reducefn.ReduceRequest.event_time:type_name -> common.EventTime
	5, // 1: function.reducefn.ReduceRequest.watermark:type_name -> common.Watermark
	1, // 2: function.reducefn.ReduceResponseList.elements:type_name -> function.reducefn.ReduceResponse
	0, // 3: function.reducefn.Reduce.ReduceFn:input_type -> function.reducefn.ReduceRequest
	6, // 4: function.reducefn.Reduce.IsReady:input_type -> google.protobuf.Empty
	2, // 5: function.reducefn.Reduce.ReduceFn:output_type -> function.reducefn.ReduceResponseList
	3, // 6: function.reducefn.Reduce.IsReady:output_type -> function.reducefn.ReadyResponse
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_pkg_apis_proto_function_reducefn_reduce_proto_init() }
func file_pkg_apis_proto_function_reducefn_reduce_proto_init() {
	if File_pkg_apis_proto_function_reducefn_reduce_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_apis_proto_function_reducefn_reduce_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReduceRequest); i {
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
		file_pkg_apis_proto_function_reducefn_reduce_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReduceResponse); i {
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
		file_pkg_apis_proto_function_reducefn_reduce_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReduceResponseList); i {
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
		file_pkg_apis_proto_function_reducefn_reduce_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_apis_proto_function_reducefn_reduce_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_apis_proto_function_reducefn_reduce_proto_goTypes,
		DependencyIndexes: file_pkg_apis_proto_function_reducefn_reduce_proto_depIdxs,
		MessageInfos:      file_pkg_apis_proto_function_reducefn_reduce_proto_msgTypes,
	}.Build()
	File_pkg_apis_proto_function_reducefn_reduce_proto = out.File
	file_pkg_apis_proto_function_reducefn_reduce_proto_rawDesc = nil
	file_pkg_apis_proto_function_reducefn_reduce_proto_goTypes = nil
	file_pkg_apis_proto_function_reducefn_reduce_proto_depIdxs = nil
}
