// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: pkg/apis/proto/sink/v1/sink.proto

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
// SinkRequest represents a request element.
type SinkRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keys      []string   `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
	Value     []byte     `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	EventTime *EventTime `protobuf:"bytes,3,opt,name=event_time,json=eventTime,proto3" json:"event_time,omitempty"`
	Watermark *Watermark `protobuf:"bytes,4,opt,name=watermark,proto3" json:"watermark,omitempty"`
	Id        string     `protobuf:"bytes,5,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *SinkRequest) Reset() {
	*x = SinkRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SinkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SinkRequest) ProtoMessage() {}

func (x *SinkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SinkRequest.ProtoReflect.Descriptor instead.
func (*SinkRequest) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_sink_v1_sink_proto_rawDescGZIP(), []int{0}
}

func (x *SinkRequest) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *SinkRequest) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *SinkRequest) GetEventTime() *EventTime {
	if x != nil {
		return x.EventTime
	}
	return nil
}

func (x *SinkRequest) GetWatermark() *Watermark {
	if x != nil {
		return x.Watermark
	}
	return nil
}

func (x *SinkRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
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
		mi := &file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadyResponse) ProtoMessage() {}

func (x *ReadyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[1]
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
	return file_pkg_apis_proto_sink_v1_sink_proto_rawDescGZIP(), []int{1}
}

func (x *ReadyResponse) GetReady() bool {
	if x != nil {
		return x.Ready
	}
	return false
}

// *
// SinkResponse is the individual response of each message written to the sink.
type SinkResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id is the ID of the message, can be used to uniquely identify the message.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// success denotes the status of persisting to disk. if set to false, it means writing to sink for the message failed.
	Success bool `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
	// err_msg is the error message, set it if success is set to false.
	ErrMsg string `protobuf:"bytes,3,opt,name=err_msg,json=errMsg,proto3" json:"err_msg,omitempty"`
}

func (x *SinkResponse) Reset() {
	*x = SinkResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SinkResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SinkResponse) ProtoMessage() {}

func (x *SinkResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SinkResponse.ProtoReflect.Descriptor instead.
func (*SinkResponse) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_sink_v1_sink_proto_rawDescGZIP(), []int{2}
}

func (x *SinkResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *SinkResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *SinkResponse) GetErrMsg() string {
	if x != nil {
		return x.ErrMsg
	}
	return ""
}

// *
// SinkResponseList is the list of responses. The number of elements in this list will be equal to the number of requests
// passed to the SinkFn.
type SinkResponseList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Responses []*SinkResponse `protobuf:"bytes,1,rep,name=responses,proto3" json:"responses,omitempty"`
}

func (x *SinkResponseList) Reset() {
	*x = SinkResponseList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SinkResponseList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SinkResponseList) ProtoMessage() {}

func (x *SinkResponseList) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SinkResponseList.ProtoReflect.Descriptor instead.
func (*SinkResponseList) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_sink_v1_sink_proto_rawDescGZIP(), []int{3}
}

func (x *SinkResponseList) GetResponses() []*SinkResponse {
	if x != nil {
		return x.Responses
	}
	return nil
}

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
		mi := &file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventTime) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventTime) ProtoMessage() {}

func (x *EventTime) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[4]
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
	return file_pkg_apis_proto_sink_v1_sink_proto_rawDescGZIP(), []int{4}
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
		mi := &file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Watermark) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Watermark) ProtoMessage() {}

func (x *Watermark) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[5]
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
	return file_pkg_apis_proto_sink_v1_sink_proto_rawDescGZIP(), []int{5}
}

func (x *Watermark) GetWatermark() *timestamppb.Timestamp {
	if x != nil {
		return x.Watermark
	}
	return nil
}

var File_pkg_apis_proto_sink_v1_sink_proto protoreflect.FileDescriptor

var file_pkg_apis_proto_sink_v1_sink_proto_rawDesc = []byte{
	0x0a, 0x21, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x73, 0x69, 0x6e, 0x6b, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x07, 0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d,
	0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xac, 0x01, 0x0a, 0x0b, 0x53,
	0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65,
	0x79, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x12, 0x31, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x69,
	0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x73, 0x69, 0x6e, 0x6b, 0x2e,
	0x76, 0x31, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x52, 0x09, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x30, 0x0a, 0x09, 0x77, 0x61, 0x74, 0x65, 0x72,
	0x6d, 0x61, 0x72, 0x6b, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x73, 0x69, 0x6e,
	0x6b, 0x2e, 0x76, 0x31, 0x2e, 0x57, 0x61, 0x74, 0x65, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x52, 0x09,
	0x77, 0x61, 0x74, 0x65, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x25, 0x0a, 0x0d, 0x52, 0x65, 0x61,
	0x64, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x65,
	0x61, 0x64, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x72, 0x65, 0x61, 0x64, 0x79,
	0x22, 0x51, 0x0a, 0x0c, 0x53, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x17, 0x0a, 0x07, 0x65, 0x72,
	0x72, 0x5f, 0x6d, 0x73, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x65, 0x72, 0x72,
	0x4d, 0x73, 0x67, 0x22, 0x47, 0x0a, 0x10, 0x53, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x33, 0x0a, 0x09, 0x72, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x73, 0x69, 0x6e,
	0x6b, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x52, 0x09, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x73, 0x22, 0x46, 0x0a, 0x09,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x54, 0x69, 0x6d, 0x65, 0x22, 0x45, 0x0a, 0x09, 0x57, 0x61, 0x74, 0x65, 0x72, 0x6d, 0x61, 0x72,
	0x6b, 0x12, 0x38, 0x0a, 0x09, 0x77, 0x61, 0x74, 0x65, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x09, 0x77, 0x61, 0x74, 0x65, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x32, 0x7e, 0x0a, 0x04, 0x53,
	0x69, 0x6e, 0x6b, 0x12, 0x3b, 0x0a, 0x06, 0x53, 0x69, 0x6e, 0x6b, 0x46, 0x6e, 0x12, 0x14, 0x2e,
	0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x69,
	0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x28, 0x01,
	0x12, 0x39, 0x0a, 0x07, 0x49, 0x73, 0x52, 0x65, 0x61, 0x64, 0x79, 0x12, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65,
	0x61, 0x64, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x38, 0x5a, 0x36, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x75, 0x6d, 0x61, 0x70, 0x72,
	0x6f, 0x6a, 0x2f, 0x6e, 0x75, 0x6d, 0x61, 0x66, 0x6c, 0x6f, 0x77, 0x2d, 0x67, 0x6f, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x69,
	0x6e, 0x6b, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_apis_proto_sink_v1_sink_proto_rawDescOnce sync.Once
	file_pkg_apis_proto_sink_v1_sink_proto_rawDescData = file_pkg_apis_proto_sink_v1_sink_proto_rawDesc
)

func file_pkg_apis_proto_sink_v1_sink_proto_rawDescGZIP() []byte {
	file_pkg_apis_proto_sink_v1_sink_proto_rawDescOnce.Do(func() {
		file_pkg_apis_proto_sink_v1_sink_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_apis_proto_sink_v1_sink_proto_rawDescData)
	})
	return file_pkg_apis_proto_sink_v1_sink_proto_rawDescData
}

var file_pkg_apis_proto_sink_v1_sink_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_pkg_apis_proto_sink_v1_sink_proto_goTypes = []interface{}{
	(*SinkRequest)(nil),           // 0: sink.v1.SinkRequest
	(*ReadyResponse)(nil),         // 1: sink.v1.ReadyResponse
	(*SinkResponse)(nil),          // 2: sink.v1.SinkResponse
	(*SinkResponseList)(nil),      // 3: sink.v1.SinkResponseList
	(*EventTime)(nil),             // 4: sink.v1.EventTime
	(*Watermark)(nil),             // 5: sink.v1.Watermark
	(*timestamppb.Timestamp)(nil), // 6: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),         // 7: google.protobuf.Empty
}
var file_pkg_apis_proto_sink_v1_sink_proto_depIdxs = []int32{
	4, // 0: sink.v1.SinkRequest.event_time:type_name -> sink.v1.EventTime
	5, // 1: sink.v1.SinkRequest.watermark:type_name -> sink.v1.Watermark
	2, // 2: sink.v1.SinkResponseList.responses:type_name -> sink.v1.SinkResponse
	6, // 3: sink.v1.EventTime.event_time:type_name -> google.protobuf.Timestamp
	6, // 4: sink.v1.Watermark.watermark:type_name -> google.protobuf.Timestamp
	0, // 5: sink.v1.Sink.SinkFn:input_type -> sink.v1.SinkRequest
	7, // 6: sink.v1.Sink.IsReady:input_type -> google.protobuf.Empty
	3, // 7: sink.v1.Sink.SinkFn:output_type -> sink.v1.SinkResponseList
	1, // 8: sink.v1.Sink.IsReady:output_type -> sink.v1.ReadyResponse
	7, // [7:9] is the sub-list for method output_type
	5, // [5:7] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_pkg_apis_proto_sink_v1_sink_proto_init() }
func file_pkg_apis_proto_sink_v1_sink_proto_init() {
	if File_pkg_apis_proto_sink_v1_sink_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SinkRequest); i {
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
		file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SinkResponse); i {
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
		file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SinkResponseList); i {
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
		file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
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
		file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_apis_proto_sink_v1_sink_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_apis_proto_sink_v1_sink_proto_goTypes,
		DependencyIndexes: file_pkg_apis_proto_sink_v1_sink_proto_depIdxs,
		MessageInfos:      file_pkg_apis_proto_sink_v1_sink_proto_msgTypes,
	}.Build()
	File_pkg_apis_proto_sink_v1_sink_proto = out.File
	file_pkg_apis_proto_sink_v1_sink_proto_rawDesc = nil
	file_pkg_apis_proto_sink_v1_sink_proto_goTypes = nil
	file_pkg_apis_proto_sink_v1_sink_proto_depIdxs = nil
}
