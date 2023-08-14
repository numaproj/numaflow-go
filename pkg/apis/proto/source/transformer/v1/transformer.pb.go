// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: pkg/apis/proto/source/transformer/v1/transformer.proto

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

type EventTime struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// event_time is the time associated with each datum.
	EventTime *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=event_time,json=eventTime,proto3" json:"event_time,omitempty"`
}

func (x *EventTime) Reset() {
	*x = EventTime{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_source_transformer_v1_transformer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventTime) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventTime) ProtoMessage() {}

func (x *EventTime) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_source_transformer_v1_transformer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventTime.ProtoReflect.Descriptor instead.
func (*EventTime) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_source_transformer_v1_transformer_proto_rawDescGZIP(), []int{0}
}

func (x *EventTime) GetEventTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EventTime
	}
	return nil
}

type Watermark struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// watermark is the monotonically increasing time which denotes completeness for the given time for the given vertex.
	Watermark *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=watermark,proto3" json:"watermark,omitempty"` // future we can add LATE, ON_TIME etc.
}

func (x *Watermark) Reset() {
	*x = Watermark{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_source_transformer_v1_transformer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Watermark) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Watermark) ProtoMessage() {}

func (x *Watermark) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_source_transformer_v1_transformer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Watermark.ProtoReflect.Descriptor instead.
func (*Watermark) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_source_transformer_v1_transformer_proto_rawDescGZIP(), []int{1}
}

func (x *Watermark) GetWatermark() *timestamppb.Timestamp {
	if x != nil {
		return x.Watermark
	}
	return nil
}

// *
// MapRequest represents a datum request element.
type SourceTransformerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keys      []string   `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
	Value     []byte     `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	EventTime *EventTime `protobuf:"bytes,3,opt,name=event_time,json=eventTime,proto3" json:"event_time,omitempty"`
	Watermark *Watermark `protobuf:"bytes,4,opt,name=watermark,proto3" json:"watermark,omitempty"`
}

func (x *SourceTransformerRequest) Reset() {
	*x = SourceTransformerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_source_transformer_v1_transformer_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SourceTransformerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SourceTransformerRequest) ProtoMessage() {}

func (x *SourceTransformerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_source_transformer_v1_transformer_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SourceTransformerRequest.ProtoReflect.Descriptor instead.
func (*SourceTransformerRequest) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_source_transformer_v1_transformer_proto_rawDescGZIP(), []int{2}
}

func (x *SourceTransformerRequest) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *SourceTransformerRequest) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *SourceTransformerRequest) GetEventTime() *EventTime {
	if x != nil {
		return x.EventTime
	}
	return nil
}

func (x *SourceTransformerRequest) GetWatermark() *Watermark {
	if x != nil {
		return x.Watermark
	}
	return nil
}

// *
// DatumResponse represents a datum response element.
type SourceTransformerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keys      []string   `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
	Value     []byte     `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	EventTime *EventTime `protobuf:"bytes,3,opt,name=event_time,json=eventTime,proto3" json:"event_time,omitempty"`
	Tags      []string   `protobuf:"bytes,4,rep,name=tags,proto3" json:"tags,omitempty"`
}

func (x *SourceTransformerResponse) Reset() {
	*x = SourceTransformerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_source_transformer_v1_transformer_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SourceTransformerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SourceTransformerResponse) ProtoMessage() {}

func (x *SourceTransformerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_source_transformer_v1_transformer_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SourceTransformerResponse.ProtoReflect.Descriptor instead.
func (*SourceTransformerResponse) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_source_transformer_v1_transformer_proto_rawDescGZIP(), []int{3}
}

func (x *SourceTransformerResponse) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *SourceTransformerResponse) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *SourceTransformerResponse) GetEventTime() *EventTime {
	if x != nil {
		return x.EventTime
	}
	return nil
}

func (x *SourceTransformerResponse) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

// *
// DatumResponseList represents a list of datum response elements.
type SourceTransformerResponseList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Elements []*SourceTransformerResponse `protobuf:"bytes,1,rep,name=elements,proto3" json:"elements,omitempty"`
}

func (x *SourceTransformerResponseList) Reset() {
	*x = SourceTransformerResponseList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_source_transformer_v1_transformer_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SourceTransformerResponseList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SourceTransformerResponseList) ProtoMessage() {}

func (x *SourceTransformerResponseList) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_source_transformer_v1_transformer_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SourceTransformerResponseList.ProtoReflect.Descriptor instead.
func (*SourceTransformerResponseList) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_source_transformer_v1_transformer_proto_rawDescGZIP(), []int{4}
}

func (x *SourceTransformerResponseList) GetElements() []*SourceTransformerResponse {
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
		mi := &file_pkg_apis_proto_source_transformer_v1_transformer_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadyResponse) ProtoMessage() {}

func (x *ReadyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_source_transformer_v1_transformer_proto_msgTypes[5]
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
	return file_pkg_apis_proto_source_transformer_v1_transformer_proto_rawDescGZIP(), []int{5}
}

func (x *ReadyResponse) GetReady() bool {
	if x != nil {
		return x.Ready
	}
	return false
}

var File_pkg_apis_proto_source_transformer_v1_transformer_proto protoreflect.FileDescriptor

var file_pkg_apis_proto_source_transformer_v1_transformer_proto_rawDesc = []byte{
	0x0a, 0x36, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72,
	0x6d, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d,
	0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x46, 0x0a, 0x09, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x39,
	0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x45, 0x0a, 0x09, 0x57, 0x61, 0x74,
	0x65, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x12, 0x38, 0x0a, 0x09, 0x77, 0x61, 0x74, 0x65, 0x72, 0x6d,
	0x61, 0x72, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x77, 0x61, 0x74, 0x65, 0x72, 0x6d, 0x61, 0x72, 0x6b,
	0x22, 0xad, 0x01, 0x0a, 0x18, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65, 0x79,
	0x73, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x33, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d,
	0x65, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x32, 0x0a, 0x09,
	0x77, 0x61, 0x74, 0x65, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x14, 0x2e, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x57, 0x61, 0x74, 0x65,
	0x72, 0x6d, 0x61, 0x72, 0x6b, 0x52, 0x09, 0x77, 0x61, 0x74, 0x65, 0x72, 0x6d, 0x61, 0x72, 0x6b,
	0x22, 0x8e, 0x01, 0x0a, 0x19, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65,
	0x79, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x33, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x69,
	0x6d, 0x65, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x74, 0x61, 0x67, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x74, 0x61, 0x67,
	0x73, 0x22, 0x61, 0x0a, 0x1d, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4c, 0x69,
	0x73, 0x74, 0x12, 0x40, 0x0a, 0x08, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x08, 0x65, 0x6c, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x73, 0x22, 0x25, 0x0a, 0x0d, 0x52, 0x65, 0x61, 0x64, 0x79, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x65, 0x61, 0x64, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x72, 0x65, 0x61, 0x64, 0x79, 0x32, 0xb4, 0x01, 0x0a, 0x11,
	0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x65,
	0x72, 0x12, 0x62, 0x0a, 0x11, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72, 0x12, 0x23, 0x2e, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f,
	0x72, 0x6d, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x72,
	0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x3b, 0x0a, 0x07, 0x49, 0x73, 0x52, 0x65, 0x61, 0x64, 0x79,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x18, 0x2e, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0x46, 0x5a, 0x44, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x6e, 0x75, 0x6d, 0x61, 0x70, 0x72, 0x6f, 0x6a, 0x2f, 0x6e, 0x75, 0x6d, 0x61, 0x66, 0x6c,
	0x6f, 0x77, 0x2d, 0x67, 0x6f, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2f, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_pkg_apis_proto_source_transformer_v1_transformer_proto_rawDescOnce sync.Once
	file_pkg_apis_proto_source_transformer_v1_transformer_proto_rawDescData = file_pkg_apis_proto_source_transformer_v1_transformer_proto_rawDesc
)

func file_pkg_apis_proto_source_transformer_v1_transformer_proto_rawDescGZIP() []byte {
	file_pkg_apis_proto_source_transformer_v1_transformer_proto_rawDescOnce.Do(func() {
		file_pkg_apis_proto_source_transformer_v1_transformer_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_apis_proto_source_transformer_v1_transformer_proto_rawDescData)
	})
	return file_pkg_apis_proto_source_transformer_v1_transformer_proto_rawDescData
}

var file_pkg_apis_proto_source_transformer_v1_transformer_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_pkg_apis_proto_source_transformer_v1_transformer_proto_goTypes = []interface{}{
	(*EventTime)(nil),                     // 0: source.v1.EventTime
	(*Watermark)(nil),                     // 1: source.v1.Watermark
	(*SourceTransformerRequest)(nil),      // 2: source.v1.SourceTransformerRequest
	(*SourceTransformerResponse)(nil),     // 3: source.v1.SourceTransformerResponse
	(*SourceTransformerResponseList)(nil), // 4: source.v1.SourceTransformerResponseList
	(*ReadyResponse)(nil),                 // 5: source.v1.ReadyResponse
	(*timestamppb.Timestamp)(nil),         // 6: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),                 // 7: google.protobuf.Empty
}
var file_pkg_apis_proto_source_transformer_v1_transformer_proto_depIdxs = []int32{
	6, // 0: source.v1.EventTime.event_time:type_name -> google.protobuf.Timestamp
	6, // 1: source.v1.Watermark.watermark:type_name -> google.protobuf.Timestamp
	0, // 2: source.v1.SourceTransformerRequest.event_time:type_name -> source.v1.EventTime
	1, // 3: source.v1.SourceTransformerRequest.watermark:type_name -> source.v1.Watermark
	0, // 4: source.v1.SourceTransformerResponse.event_time:type_name -> source.v1.EventTime
	3, // 5: source.v1.SourceTransformerResponseList.elements:type_name -> source.v1.SourceTransformerResponse
	2, // 6: source.v1.SourceTransformer.SourceTransformer:input_type -> source.v1.SourceTransformerRequest
	7, // 7: source.v1.SourceTransformer.IsReady:input_type -> google.protobuf.Empty
	4, // 8: source.v1.SourceTransformer.SourceTransformer:output_type -> source.v1.SourceTransformerResponseList
	5, // 9: source.v1.SourceTransformer.IsReady:output_type -> source.v1.ReadyResponse
	8, // [8:10] is the sub-list for method output_type
	6, // [6:8] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_pkg_apis_proto_source_transformer_v1_transformer_proto_init() }
func file_pkg_apis_proto_source_transformer_v1_transformer_proto_init() {
	if File_pkg_apis_proto_source_transformer_v1_transformer_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_apis_proto_source_transformer_v1_transformer_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventTime); i {
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
		file_pkg_apis_proto_source_transformer_v1_transformer_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Watermark); i {
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
		file_pkg_apis_proto_source_transformer_v1_transformer_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SourceTransformerRequest); i {
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
		file_pkg_apis_proto_source_transformer_v1_transformer_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SourceTransformerResponse); i {
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
		file_pkg_apis_proto_source_transformer_v1_transformer_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SourceTransformerResponseList); i {
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
		file_pkg_apis_proto_source_transformer_v1_transformer_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
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
			RawDescriptor: file_pkg_apis_proto_source_transformer_v1_transformer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_apis_proto_source_transformer_v1_transformer_proto_goTypes,
		DependencyIndexes: file_pkg_apis_proto_source_transformer_v1_transformer_proto_depIdxs,
		MessageInfos:      file_pkg_apis_proto_source_transformer_v1_transformer_proto_msgTypes,
	}.Build()
	File_pkg_apis_proto_source_transformer_v1_transformer_proto = out.File
	file_pkg_apis_proto_source_transformer_v1_transformer_proto_rawDesc = nil
	file_pkg_apis_proto_source_transformer_v1_transformer_proto_goTypes = nil
	file_pkg_apis_proto_source_transformer_v1_transformer_proto_depIdxs = nil
}
