// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v4.25.1
// source: pkg/apis/proto/sessionreduce/v1/sessionreduce.proto

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

type SessionReduceRequest_WindowOperation_Event int32

const (
	SessionReduceRequest_WindowOperation_OPEN   SessionReduceRequest_WindowOperation_Event = 0
	SessionReduceRequest_WindowOperation_CLOSE  SessionReduceRequest_WindowOperation_Event = 1
	SessionReduceRequest_WindowOperation_EXPAND SessionReduceRequest_WindowOperation_Event = 2
	SessionReduceRequest_WindowOperation_MERGE  SessionReduceRequest_WindowOperation_Event = 3
	SessionReduceRequest_WindowOperation_APPEND SessionReduceRequest_WindowOperation_Event = 4
)

// Enum value maps for SessionReduceRequest_WindowOperation_Event.
var (
	SessionReduceRequest_WindowOperation_Event_name = map[int32]string{
		0: "OPEN",
		1: "CLOSE",
		2: "EXPAND",
		3: "MERGE",
		4: "APPEND",
	}
	SessionReduceRequest_WindowOperation_Event_value = map[string]int32{
		"OPEN":   0,
		"CLOSE":  1,
		"EXPAND": 2,
		"MERGE":  3,
		"APPEND": 4,
	}
)

func (x SessionReduceRequest_WindowOperation_Event) Enum() *SessionReduceRequest_WindowOperation_Event {
	p := new(SessionReduceRequest_WindowOperation_Event)
	*p = x
	return p
}

func (x SessionReduceRequest_WindowOperation_Event) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SessionReduceRequest_WindowOperation_Event) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_enumTypes[0].Descriptor()
}

func (SessionReduceRequest_WindowOperation_Event) Type() protoreflect.EnumType {
	return &file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_enumTypes[0]
}

func (x SessionReduceRequest_WindowOperation_Event) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SessionReduceRequest_WindowOperation_Event.Descriptor instead.
func (SessionReduceRequest_WindowOperation_Event) EnumDescriptor() ([]byte, []int) {
	return file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_rawDescGZIP(), []int{1, 0, 0}
}

// KeyedWindow represents a window with keys.
// since the client track the keys, we use keyed window.
type KeyedWindow struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Start *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=start,proto3" json:"start,omitempty"`
	End   *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=end,proto3" json:"end,omitempty"`
	Slot  string                 `protobuf:"bytes,3,opt,name=slot,proto3" json:"slot,omitempty"`
	Keys  []string               `protobuf:"bytes,4,rep,name=keys,proto3" json:"keys,omitempty"`
}

func (x *KeyedWindow) Reset() {
	*x = KeyedWindow{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeyedWindow) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyedWindow) ProtoMessage() {}

func (x *KeyedWindow) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyedWindow.ProtoReflect.Descriptor instead.
func (*KeyedWindow) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_rawDescGZIP(), []int{0}
}

func (x *KeyedWindow) GetStart() *timestamppb.Timestamp {
	if x != nil {
		return x.Start
	}
	return nil
}

func (x *KeyedWindow) GetEnd() *timestamppb.Timestamp {
	if x != nil {
		return x.End
	}
	return nil
}

func (x *KeyedWindow) GetSlot() string {
	if x != nil {
		return x.Slot
	}
	return ""
}

func (x *KeyedWindow) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

// *
// SessionReduceRequest represents a request element.
type SessionReduceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Payload   *SessionReduceRequest_Payload         `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	Operation *SessionReduceRequest_WindowOperation `protobuf:"bytes,2,opt,name=operation,proto3" json:"operation,omitempty"`
}

func (x *SessionReduceRequest) Reset() {
	*x = SessionReduceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SessionReduceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionReduceRequest) ProtoMessage() {}

func (x *SessionReduceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionReduceRequest.ProtoReflect.Descriptor instead.
func (*SessionReduceRequest) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_rawDescGZIP(), []int{1}
}

func (x *SessionReduceRequest) GetPayload() *SessionReduceRequest_Payload {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *SessionReduceRequest) GetOperation() *SessionReduceRequest_WindowOperation {
	if x != nil {
		return x.Operation
	}
	return nil
}

// *
// SessionReduceResponse represents a response element.
type SessionReduceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result *SessionReduceResponse_Result `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	// keyedWindow represents a window to which the result belongs.
	KeyedWindow *KeyedWindow `protobuf:"bytes,2,opt,name=keyedWindow,proto3" json:"keyedWindow,omitempty"`
	// EOF represents the end of the response for a window.
	EOF bool `protobuf:"varint,3,opt,name=EOF,proto3" json:"EOF,omitempty"`
}

func (x *SessionReduceResponse) Reset() {
	*x = SessionReduceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SessionReduceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionReduceResponse) ProtoMessage() {}

func (x *SessionReduceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionReduceResponse.ProtoReflect.Descriptor instead.
func (*SessionReduceResponse) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_rawDescGZIP(), []int{2}
}

func (x *SessionReduceResponse) GetResult() *SessionReduceResponse_Result {
	if x != nil {
		return x.Result
	}
	return nil
}

func (x *SessionReduceResponse) GetKeyedWindow() *KeyedWindow {
	if x != nil {
		return x.KeyedWindow
	}
	return nil
}

func (x *SessionReduceResponse) GetEOF() bool {
	if x != nil {
		return x.EOF
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
		mi := &file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadyResponse) ProtoMessage() {}

func (x *ReadyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_msgTypes[3]
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
	return file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_rawDescGZIP(), []int{3}
}

func (x *ReadyResponse) GetReady() bool {
	if x != nil {
		return x.Ready
	}
	return false
}

// WindowOperation represents a window operation.
// For Aligned window values can be one of OPEN, CLOSE, EXPAND, MERGE and APPEND.
type SessionReduceRequest_WindowOperation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Event        SessionReduceRequest_WindowOperation_Event `protobuf:"varint,1,opt,name=event,proto3,enum=sessionreduce.v1.SessionReduceRequest_WindowOperation_Event" json:"event,omitempty"`
	KeyedWindows []*KeyedWindow                             `protobuf:"bytes,2,rep,name=keyedWindows,proto3" json:"keyedWindows,omitempty"`
}

func (x *SessionReduceRequest_WindowOperation) Reset() {
	*x = SessionReduceRequest_WindowOperation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SessionReduceRequest_WindowOperation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionReduceRequest_WindowOperation) ProtoMessage() {}

func (x *SessionReduceRequest_WindowOperation) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionReduceRequest_WindowOperation.ProtoReflect.Descriptor instead.
func (*SessionReduceRequest_WindowOperation) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_rawDescGZIP(), []int{1, 0}
}

func (x *SessionReduceRequest_WindowOperation) GetEvent() SessionReduceRequest_WindowOperation_Event {
	if x != nil {
		return x.Event
	}
	return SessionReduceRequest_WindowOperation_OPEN
}

func (x *SessionReduceRequest_WindowOperation) GetKeyedWindows() []*KeyedWindow {
	if x != nil {
		return x.KeyedWindows
	}
	return nil
}

// Payload represents a payload element.
type SessionReduceRequest_Payload struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keys      []string               `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
	Value     []byte                 `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	EventTime *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=event_time,json=eventTime,proto3" json:"event_time,omitempty"`
	Watermark *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=watermark,proto3" json:"watermark,omitempty"`
	Headers   map[string]string      `protobuf:"bytes,5,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *SessionReduceRequest_Payload) Reset() {
	*x = SessionReduceRequest_Payload{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SessionReduceRequest_Payload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionReduceRequest_Payload) ProtoMessage() {}

func (x *SessionReduceRequest_Payload) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionReduceRequest_Payload.ProtoReflect.Descriptor instead.
func (*SessionReduceRequest_Payload) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_rawDescGZIP(), []int{1, 1}
}

func (x *SessionReduceRequest_Payload) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *SessionReduceRequest_Payload) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *SessionReduceRequest_Payload) GetEventTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EventTime
	}
	return nil
}

func (x *SessionReduceRequest_Payload) GetWatermark() *timestamppb.Timestamp {
	if x != nil {
		return x.Watermark
	}
	return nil
}

func (x *SessionReduceRequest_Payload) GetHeaders() map[string]string {
	if x != nil {
		return x.Headers
	}
	return nil
}

// Result represents a result element. It contains the result of the reduce function.
type SessionReduceResponse_Result struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keys  []string `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
	Value []byte   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Tags  []string `protobuf:"bytes,3,rep,name=tags,proto3" json:"tags,omitempty"`
}

func (x *SessionReduceResponse_Result) Reset() {
	*x = SessionReduceResponse_Result{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SessionReduceResponse_Result) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionReduceResponse_Result) ProtoMessage() {}

func (x *SessionReduceResponse_Result) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionReduceResponse_Result.ProtoReflect.Descriptor instead.
func (*SessionReduceResponse_Result) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_rawDescGZIP(), []int{2, 0}
}

func (x *SessionReduceResponse_Result) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *SessionReduceResponse_Result) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *SessionReduceResponse_Result) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

var File_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto protoreflect.FileDescriptor

var file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_rawDesc = []byte{
	0x0a, 0x33, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x2f, 0x76,
	0x31, 0x2f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x72, 0x65,
	0x64, 0x75, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x95, 0x01, 0x0a, 0x0b, 0x4b, 0x65, 0x79, 0x65, 0x64, 0x57,
	0x69, 0x6e, 0x64, 0x6f, 0x77, 0x12, 0x30, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12, 0x2c, 0x0a, 0x03, 0x65, 0x6e, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x03, 0x65, 0x6e, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6c, 0x6f, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x73, 0x6c, 0x6f, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79,
	0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x22, 0xe0, 0x05,
	0x0a, 0x14, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x48, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e,
	0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64,
	0x12, 0x54, 0x0a, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x36, 0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x72, 0x65, 0x64,
	0x75, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x64, 0x75, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x57, 0x69, 0x6e, 0x64,
	0x6f, 0x77, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x6f, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0xe9, 0x01, 0x0a, 0x0f, 0x57, 0x69, 0x6e, 0x64, 0x6f,
	0x77, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x52, 0x0a, 0x05, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x3c, 0x2e, 0x73, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x2e, 0x57, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x41,
	0x0a, 0x0c, 0x6b, 0x65, 0x79, 0x65, 0x64, 0x57, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x72, 0x65,
	0x64, 0x75, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4b, 0x65, 0x79, 0x65, 0x64, 0x57, 0x69, 0x6e,
	0x64, 0x6f, 0x77, 0x52, 0x0c, 0x6b, 0x65, 0x79, 0x65, 0x64, 0x57, 0x69, 0x6e, 0x64, 0x6f, 0x77,
	0x73, 0x22, 0x3f, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x08, 0x0a, 0x04, 0x4f, 0x50,
	0x45, 0x4e, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x43, 0x4c, 0x4f, 0x53, 0x45, 0x10, 0x01, 0x12,
	0x0a, 0x0a, 0x06, 0x45, 0x58, 0x50, 0x41, 0x4e, 0x44, 0x10, 0x02, 0x12, 0x09, 0x0a, 0x05, 0x4d,
	0x45, 0x52, 0x47, 0x45, 0x10, 0x03, 0x12, 0x0a, 0x0a, 0x06, 0x41, 0x50, 0x50, 0x45, 0x4e, 0x44,
	0x10, 0x04, 0x1a, 0xbb, 0x02, 0x0a, 0x07, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65,
	0x79, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x38, 0x0a, 0x09, 0x77, 0x61, 0x74, 0x65, 0x72, 0x6d, 0x61, 0x72, 0x6b,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x09, 0x77, 0x61, 0x74, 0x65, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x12, 0x55, 0x0a,
	0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x3b,
	0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x48,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x68, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x73, 0x1a, 0x3a, 0x0a, 0x0c, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x22, 0xfa, 0x01, 0x0a, 0x15, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x64, 0x75,
	0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x46, 0x0a, 0x06, 0x72, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x73, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x12, 0x3f, 0x0a, 0x0b, 0x6b, 0x65, 0x79, 0x65, 0x64, 0x57, 0x69, 0x6e, 0x64, 0x6f,
	0x77, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4b, 0x65, 0x79, 0x65, 0x64,
	0x57, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x52, 0x0b, 0x6b, 0x65, 0x79, 0x65, 0x64, 0x57, 0x69, 0x6e,
	0x64, 0x6f, 0x77, 0x12, 0x10, 0x0a, 0x03, 0x45, 0x4f, 0x46, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x03, 0x45, 0x4f, 0x46, 0x1a, 0x46, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6b,
	0x65, 0x79, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x61, 0x67,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x22, 0x25, 0x0a,
	0x0d, 0x52, 0x65, 0x61, 0x64, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x72, 0x65, 0x61, 0x64, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x72,
	0x65, 0x61, 0x64, 0x79, 0x32, 0xbb, 0x01, 0x0a, 0x0d, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x12, 0x66, 0x0a, 0x0f, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x46, 0x6e, 0x12, 0x26, 0x2e, 0x73, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x27, 0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x72, 0x65, 0x64, 0x75, 0x63,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x64, 0x75,
	0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x30, 0x01, 0x12, 0x42,
	0x0a, 0x07, 0x49, 0x73, 0x52, 0x65, 0x61, 0x64, 0x79, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x1a, 0x1f, 0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x72, 0x65, 0x64, 0x75, 0x63,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0x40, 0x5a, 0x3e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x6e, 0x75, 0x6d, 0x61, 0x70, 0x72, 0x6f, 0x6a, 0x2f, 0x6e, 0x75, 0x6d, 0x61, 0x66, 0x6c,
	0x6f, 0x77, 0x2d, 0x67, 0x6f, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x73, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_rawDescOnce sync.Once
	file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_rawDescData = file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_rawDesc
)

func file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_rawDescGZIP() []byte {
	file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_rawDescOnce.Do(func() {
		file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_rawDescData)
	})
	return file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_rawDescData
}

var file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_goTypes = []any{
	(SessionReduceRequest_WindowOperation_Event)(0), // 0: sessionreduce.v1.SessionReduceRequest.WindowOperation.Event
	(*KeyedWindow)(nil),                             // 1: sessionreduce.v1.KeyedWindow
	(*SessionReduceRequest)(nil),                    // 2: sessionreduce.v1.SessionReduceRequest
	(*SessionReduceResponse)(nil),                   // 3: sessionreduce.v1.SessionReduceResponse
	(*ReadyResponse)(nil),                           // 4: sessionreduce.v1.ReadyResponse
	(*SessionReduceRequest_WindowOperation)(nil),    // 5: sessionreduce.v1.SessionReduceRequest.WindowOperation
	(*SessionReduceRequest_Payload)(nil),            // 6: sessionreduce.v1.SessionReduceRequest.Payload
	nil,                                             // 7: sessionreduce.v1.SessionReduceRequest.Payload.HeadersEntry
	(*SessionReduceResponse_Result)(nil),            // 8: sessionreduce.v1.SessionReduceResponse.Result
	(*timestamppb.Timestamp)(nil),                   // 9: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),                           // 10: google.protobuf.Empty
}
var file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_depIdxs = []int32{
	9,  // 0: sessionreduce.v1.KeyedWindow.start:type_name -> google.protobuf.Timestamp
	9,  // 1: sessionreduce.v1.KeyedWindow.end:type_name -> google.protobuf.Timestamp
	6,  // 2: sessionreduce.v1.SessionReduceRequest.payload:type_name -> sessionreduce.v1.SessionReduceRequest.Payload
	5,  // 3: sessionreduce.v1.SessionReduceRequest.operation:type_name -> sessionreduce.v1.SessionReduceRequest.WindowOperation
	8,  // 4: sessionreduce.v1.SessionReduceResponse.result:type_name -> sessionreduce.v1.SessionReduceResponse.Result
	1,  // 5: sessionreduce.v1.SessionReduceResponse.keyedWindow:type_name -> sessionreduce.v1.KeyedWindow
	0,  // 6: sessionreduce.v1.SessionReduceRequest.WindowOperation.event:type_name -> sessionreduce.v1.SessionReduceRequest.WindowOperation.Event
	1,  // 7: sessionreduce.v1.SessionReduceRequest.WindowOperation.keyedWindows:type_name -> sessionreduce.v1.KeyedWindow
	9,  // 8: sessionreduce.v1.SessionReduceRequest.Payload.event_time:type_name -> google.protobuf.Timestamp
	9,  // 9: sessionreduce.v1.SessionReduceRequest.Payload.watermark:type_name -> google.protobuf.Timestamp
	7,  // 10: sessionreduce.v1.SessionReduceRequest.Payload.headers:type_name -> sessionreduce.v1.SessionReduceRequest.Payload.HeadersEntry
	2,  // 11: sessionreduce.v1.SessionReduce.SessionReduceFn:input_type -> sessionreduce.v1.SessionReduceRequest
	10, // 12: sessionreduce.v1.SessionReduce.IsReady:input_type -> google.protobuf.Empty
	3,  // 13: sessionreduce.v1.SessionReduce.SessionReduceFn:output_type -> sessionreduce.v1.SessionReduceResponse
	4,  // 14: sessionreduce.v1.SessionReduce.IsReady:output_type -> sessionreduce.v1.ReadyResponse
	13, // [13:15] is the sub-list for method output_type
	11, // [11:13] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_init() }
func file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_init() {
	if File_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*KeyedWindow); i {
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
		file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*SessionReduceRequest); i {
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
		file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*SessionReduceResponse); i {
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
		file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_msgTypes[3].Exporter = func(v any, i int) any {
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
		file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*SessionReduceRequest_WindowOperation); i {
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
		file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*SessionReduceRequest_Payload); i {
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
		file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*SessionReduceResponse_Result); i {
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
			RawDescriptor: file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_goTypes,
		DependencyIndexes: file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_depIdxs,
		EnumInfos:         file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_enumTypes,
		MessageInfos:      file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_msgTypes,
	}.Build()
	File_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto = out.File
	file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_rawDesc = nil
	file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_goTypes = nil
	file_pkg_apis_proto_sessionreduce_v1_sessionreduce_proto_depIdxs = nil
}
