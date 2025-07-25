// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.29.3
// source: pkg/apis/proto/sideinput/v1/sideinput.proto

package v1

import (
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
// SideInputResponse represents a response to a given side input retrieval request.
type SideInputResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// value represents the latest value of the side input payload
	Value []byte `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	// noBroadcast indicates whether the side input value should be broadcasted to all
	// True if value should not be broadcasted
	// False if value should be broadcasted
	NoBroadcast bool `protobuf:"varint,2,opt,name=no_broadcast,json=noBroadcast,proto3" json:"no_broadcast,omitempty"`
}

func (x *SideInputResponse) Reset() {
	*x = SideInputResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_sideinput_v1_sideinput_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SideInputResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SideInputResponse) ProtoMessage() {}

func (x *SideInputResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_sideinput_v1_sideinput_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SideInputResponse.ProtoReflect.Descriptor instead.
func (*SideInputResponse) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_sideinput_v1_sideinput_proto_rawDescGZIP(), []int{0}
}

func (x *SideInputResponse) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *SideInputResponse) GetNoBroadcast() bool {
	if x != nil {
		return x.NoBroadcast
	}
	return false
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
		mi := &file_pkg_apis_proto_sideinput_v1_sideinput_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadyResponse) ProtoMessage() {}

func (x *ReadyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_sideinput_v1_sideinput_proto_msgTypes[1]
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
	return file_pkg_apis_proto_sideinput_v1_sideinput_proto_rawDescGZIP(), []int{1}
}

func (x *ReadyResponse) GetReady() bool {
	if x != nil {
		return x.Ready
	}
	return false
}

var File_pkg_apis_proto_sideinput_v1_sideinput_proto protoreflect.FileDescriptor

var file_pkg_apis_proto_sideinput_v1_sideinput_proto_rawDesc = []byte{
	0x0a, 0x2b, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x73, 0x69, 0x64, 0x65, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x69,
	0x64, 0x65, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x73,
	0x69, 0x64, 0x65, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70,
	0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4c, 0x0a, 0x11, 0x53, 0x69, 0x64, 0x65,
	0x49, 0x6e, 0x70, 0x75, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x6e, 0x6f, 0x5f, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63,
	0x61, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x6e, 0x6f, 0x42, 0x72, 0x6f,
	0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x22, 0x25, 0x0a, 0x0d, 0x52, 0x65, 0x61, 0x64, 0x79, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x65, 0x61, 0x64, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x72, 0x65, 0x61, 0x64, 0x79, 0x32, 0x99, 0x01,
	0x0a, 0x09, 0x53, 0x69, 0x64, 0x65, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x4c, 0x0a, 0x11, 0x52,
	0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x53, 0x69, 0x64, 0x65, 0x49, 0x6e, 0x70, 0x75, 0x74,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1f, 0x2e, 0x73, 0x69, 0x64, 0x65, 0x69,
	0x6e, 0x70, 0x75, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x69, 0x64, 0x65, 0x49, 0x6e, 0x70, 0x75,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3e, 0x0a, 0x07, 0x49, 0x73, 0x52,
	0x65, 0x61, 0x64, 0x79, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1b, 0x2e, 0x73,
	0x69, 0x64, 0x65, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x61, 0x64,
	0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x3d, 0x5a, 0x3b, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x75, 0x6d, 0x61, 0x70, 0x72, 0x6f, 0x6a,
	0x2f, 0x6e, 0x75, 0x6d, 0x61, 0x66, 0x6c, 0x6f, 0x77, 0x2d, 0x67, 0x6f, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x69, 0x64, 0x65,
	0x69, 0x6e, 0x70, 0x75, 0x74, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_apis_proto_sideinput_v1_sideinput_proto_rawDescOnce sync.Once
	file_pkg_apis_proto_sideinput_v1_sideinput_proto_rawDescData = file_pkg_apis_proto_sideinput_v1_sideinput_proto_rawDesc
)

func file_pkg_apis_proto_sideinput_v1_sideinput_proto_rawDescGZIP() []byte {
	file_pkg_apis_proto_sideinput_v1_sideinput_proto_rawDescOnce.Do(func() {
		file_pkg_apis_proto_sideinput_v1_sideinput_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_apis_proto_sideinput_v1_sideinput_proto_rawDescData)
	})
	return file_pkg_apis_proto_sideinput_v1_sideinput_proto_rawDescData
}

var file_pkg_apis_proto_sideinput_v1_sideinput_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pkg_apis_proto_sideinput_v1_sideinput_proto_goTypes = []any{
	(*SideInputResponse)(nil), // 0: sideinput.v1.SideInputResponse
	(*ReadyResponse)(nil),     // 1: sideinput.v1.ReadyResponse
	(*emptypb.Empty)(nil),     // 2: google.protobuf.Empty
}
var file_pkg_apis_proto_sideinput_v1_sideinput_proto_depIdxs = []int32{
	2, // 0: sideinput.v1.SideInput.RetrieveSideInput:input_type -> google.protobuf.Empty
	2, // 1: sideinput.v1.SideInput.IsReady:input_type -> google.protobuf.Empty
	0, // 2: sideinput.v1.SideInput.RetrieveSideInput:output_type -> sideinput.v1.SideInputResponse
	1, // 3: sideinput.v1.SideInput.IsReady:output_type -> sideinput.v1.ReadyResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_apis_proto_sideinput_v1_sideinput_proto_init() }
func file_pkg_apis_proto_sideinput_v1_sideinput_proto_init() {
	if File_pkg_apis_proto_sideinput_v1_sideinput_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_apis_proto_sideinput_v1_sideinput_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*SideInputResponse); i {
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
		file_pkg_apis_proto_sideinput_v1_sideinput_proto_msgTypes[1].Exporter = func(v any, i int) any {
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
			RawDescriptor: file_pkg_apis_proto_sideinput_v1_sideinput_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_apis_proto_sideinput_v1_sideinput_proto_goTypes,
		DependencyIndexes: file_pkg_apis_proto_sideinput_v1_sideinput_proto_depIdxs,
		MessageInfos:      file_pkg_apis_proto_sideinput_v1_sideinput_proto_msgTypes,
	}.Build()
	File_pkg_apis_proto_sideinput_v1_sideinput_proto = out.File
	file_pkg_apis_proto_sideinput_v1_sideinput_proto_rawDesc = nil
	file_pkg_apis_proto_sideinput_v1_sideinput_proto_goTypes = nil
	file_pkg_apis_proto_sideinput_v1_sideinput_proto_depIdxs = nil
}
