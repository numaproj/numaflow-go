// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v4.25.1
// source: pkg/apis/proto/accumulator/v1/accumulator.proto

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

type AccumulatorRequest_WindowOperation_Event int32

const (
	AccumulatorRequest_WindowOperation_OPEN   AccumulatorRequest_WindowOperation_Event = 0
	AccumulatorRequest_WindowOperation_CLOSE  AccumulatorRequest_WindowOperation_Event = 1
	AccumulatorRequest_WindowOperation_APPEND AccumulatorRequest_WindowOperation_Event = 2
)

// Enum value maps for AccumulatorRequest_WindowOperation_Event.
var (
	AccumulatorRequest_WindowOperation_Event_name = map[int32]string{
		0: "OPEN",
		1: "CLOSE",
		2: "APPEND",
	}
	AccumulatorRequest_WindowOperation_Event_value = map[string]int32{
		"OPEN":   0,
		"CLOSE":  1,
		"APPEND": 2,
	}
)

func (x AccumulatorRequest_WindowOperation_Event) Enum() *AccumulatorRequest_WindowOperation_Event {
	p := new(AccumulatorRequest_WindowOperation_Event)
	*p = x
	return p
}

func (x AccumulatorRequest_WindowOperation_Event) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AccumulatorRequest_WindowOperation_Event) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_apis_proto_accumulator_v1_accumulator_proto_enumTypes[0].Descriptor()
}

func (AccumulatorRequest_WindowOperation_Event) Type() protoreflect.EnumType {
	return &file_pkg_apis_proto_accumulator_v1_accumulator_proto_enumTypes[0]
}

func (x AccumulatorRequest_WindowOperation_Event) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AccumulatorRequest_WindowOperation_Event.Descriptor instead.
func (AccumulatorRequest_WindowOperation_Event) EnumDescriptor() ([]byte, []int) {
	return file_pkg_apis_proto_accumulator_v1_accumulator_proto_rawDescGZIP(), []int{1, 0, 0}
}

// Payload represents a payload element.
type Payload struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keys      []string               `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
	Value     []byte                 `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	EventTime *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=event_time,json=eventTime,proto3" json:"event_time,omitempty"`
	Watermark *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=watermark,proto3" json:"watermark,omitempty"`
	Headers   map[string]string      `protobuf:"bytes,5,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Payload) Reset() {
	*x = Payload{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_accumulator_v1_accumulator_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Payload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Payload) ProtoMessage() {}

func (x *Payload) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_accumulator_v1_accumulator_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Payload.ProtoReflect.Descriptor instead.
func (*Payload) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_accumulator_v1_accumulator_proto_rawDescGZIP(), []int{0}
}

func (x *Payload) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *Payload) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *Payload) GetEventTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EventTime
	}
	return nil
}

func (x *Payload) GetWatermark() *timestamppb.Timestamp {
	if x != nil {
		return x.Watermark
	}
	return nil
}

func (x *Payload) GetHeaders() map[string]string {
	if x != nil {
		return x.Headers
	}
	return nil
}

// *
// AccumulatorRequest represents a request element.
type AccumulatorRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Payload *Payload `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	// This ID is used to uniquely identify a map request
	Id        string                              `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Operation *AccumulatorRequest_WindowOperation `protobuf:"bytes,3,opt,name=operation,proto3" json:"operation,omitempty"`
	Handshake *Handshake                          `protobuf:"bytes,4,opt,name=handshake,proto3,oneof" json:"handshake,omitempty"`
}

func (x *AccumulatorRequest) Reset() {
	*x = AccumulatorRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_accumulator_v1_accumulator_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccumulatorRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccumulatorRequest) ProtoMessage() {}

func (x *AccumulatorRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_accumulator_v1_accumulator_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccumulatorRequest.ProtoReflect.Descriptor instead.
func (*AccumulatorRequest) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_accumulator_v1_accumulator_proto_rawDescGZIP(), []int{1}
}

func (x *AccumulatorRequest) GetPayload() *Payload {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *AccumulatorRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *AccumulatorRequest) GetOperation() *AccumulatorRequest_WindowOperation {
	if x != nil {
		return x.Operation
	}
	return nil
}

func (x *AccumulatorRequest) GetHandshake() *Handshake {
	if x != nil {
		return x.Handshake
	}
	return nil
}

// *
// Window represents a window.
type Window struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Start *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=start,proto3" json:"start,omitempty"`
	End   *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=end,proto3" json:"end,omitempty"`
	Slot  string                 `protobuf:"bytes,3,opt,name=slot,proto3" json:"slot,omitempty"`
}

func (x *Window) Reset() {
	*x = Window{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_accumulator_v1_accumulator_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Window) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Window) ProtoMessage() {}

func (x *Window) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_accumulator_v1_accumulator_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Window.ProtoReflect.Descriptor instead.
func (*Window) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_accumulator_v1_accumulator_proto_rawDescGZIP(), []int{2}
}

func (x *Window) GetStart() *timestamppb.Timestamp {
	if x != nil {
		return x.Start
	}
	return nil
}

func (x *Window) GetEnd() *timestamppb.Timestamp {
	if x != nil {
		return x.End
	}
	return nil
}

func (x *Window) GetSlot() string {
	if x != nil {
		return x.Slot
	}
	return ""
}

// *
// AccumulatorResponse represents a response element.
type AccumulatorResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Payload *Payload `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	// This ID is used to uniquely identify a map request
	Id string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	// window represents a window to which the result belongs.
	Window    *Window    `protobuf:"bytes,3,opt,name=window,proto3" json:"window,omitempty"`
	Handshake *Handshake `protobuf:"bytes,4,opt,name=handshake,proto3,oneof" json:"handshake,omitempty"`
	// EOF represents the end of the response for a window.
	EOF bool `protobuf:"varint,5,opt,name=EOF,proto3" json:"EOF,omitempty"`
}

func (x *AccumulatorResponse) Reset() {
	*x = AccumulatorResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_accumulator_v1_accumulator_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccumulatorResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccumulatorResponse) ProtoMessage() {}

func (x *AccumulatorResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_accumulator_v1_accumulator_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccumulatorResponse.ProtoReflect.Descriptor instead.
func (*AccumulatorResponse) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_accumulator_v1_accumulator_proto_rawDescGZIP(), []int{3}
}

func (x *AccumulatorResponse) GetPayload() *Payload {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *AccumulatorResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *AccumulatorResponse) GetWindow() *Window {
	if x != nil {
		return x.Window
	}
	return nil
}

func (x *AccumulatorResponse) GetHandshake() *Handshake {
	if x != nil {
		return x.Handshake
	}
	return nil
}

func (x *AccumulatorResponse) GetEOF() bool {
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
		mi := &file_pkg_apis_proto_accumulator_v1_accumulator_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadyResponse) ProtoMessage() {}

func (x *ReadyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_accumulator_v1_accumulator_proto_msgTypes[4]
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
	return file_pkg_apis_proto_accumulator_v1_accumulator_proto_rawDescGZIP(), []int{4}
}

func (x *ReadyResponse) GetReady() bool {
	if x != nil {
		return x.Ready
	}
	return false
}

// Handshake message between client and server to indicate the start of transmission.
type Handshake struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Required field indicating the start of transmission.
	Sot bool `protobuf:"varint,1,opt,name=sot,proto3" json:"sot,omitempty"`
}

func (x *Handshake) Reset() {
	*x = Handshake{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_accumulator_v1_accumulator_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Handshake) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Handshake) ProtoMessage() {}

func (x *Handshake) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_accumulator_v1_accumulator_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Handshake.ProtoReflect.Descriptor instead.
func (*Handshake) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_accumulator_v1_accumulator_proto_rawDescGZIP(), []int{5}
}

func (x *Handshake) GetSot() bool {
	if x != nil {
		return x.Sot
	}
	return false
}

// WindowOperation represents a window operation.
// For Aligned windows, OPEN, APPEND and CLOSE events are sent.
type AccumulatorRequest_WindowOperation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Event   AccumulatorRequest_WindowOperation_Event `protobuf:"varint,1,opt,name=event,proto3,enum=accumulator.v1.AccumulatorRequest_WindowOperation_Event" json:"event,omitempty"`
	Windows []*Window                                `protobuf:"bytes,2,rep,name=windows,proto3" json:"windows,omitempty"`
}

func (x *AccumulatorRequest_WindowOperation) Reset() {
	*x = AccumulatorRequest_WindowOperation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_accumulator_v1_accumulator_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccumulatorRequest_WindowOperation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccumulatorRequest_WindowOperation) ProtoMessage() {}

func (x *AccumulatorRequest_WindowOperation) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_accumulator_v1_accumulator_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccumulatorRequest_WindowOperation.ProtoReflect.Descriptor instead.
func (*AccumulatorRequest_WindowOperation) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_accumulator_v1_accumulator_proto_rawDescGZIP(), []int{1, 0}
}

func (x *AccumulatorRequest_WindowOperation) GetEvent() AccumulatorRequest_WindowOperation_Event {
	if x != nil {
		return x.Event
	}
	return AccumulatorRequest_WindowOperation_OPEN
}

func (x *AccumulatorRequest_WindowOperation) GetWindows() []*Window {
	if x != nil {
		return x.Windows
	}
	return nil
}

var File_pkg_apis_proto_accumulator_v1_accumulator_proto protoreflect.FileDescriptor

var file_pkg_apis_proto_accumulator_v1_accumulator_proto_rawDesc = []byte{
	0x0a, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x61, 0x63, 0x63, 0x75, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2f, 0x76, 0x31, 0x2f,
	0x61, 0x63, 0x63, 0x75, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0e, 0x61, 0x63, 0x63, 0x75, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x76,
	0x31, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xa4, 0x02, 0x0a, 0x07, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6b,
	0x65, 0x79, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x38, 0x0a, 0x09, 0x77, 0x61, 0x74, 0x65, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x77, 0x61, 0x74, 0x65, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x12, 0x3e, 0x0a, 0x07, 0x68, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x61, 0x63,
	0x63, 0x75, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x1a, 0x3a, 0x0a, 0x0c, 0x48, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xb5, 0x03, 0x0a, 0x12, 0x41, 0x63, 0x63, 0x75, 0x6d,
	0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x31, 0x0a,
	0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17,
	0x2e, 0x61, 0x63, 0x63, 0x75, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e,
	0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x50, 0x0a, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x32, 0x2e, 0x61, 0x63, 0x63, 0x75, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x6f,
	0x72, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x63, 0x63, 0x75, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x57, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x4f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x3c, 0x0a, 0x09, 0x68, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x61, 0x63, 0x63, 0x75, 0x6d, 0x75, 0x6c, 0x61,
	0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65,
	0x48, 0x00, 0x52, 0x09, 0x68, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x88, 0x01, 0x01,
	0x1a, 0xbd, 0x01, 0x0a, 0x0f, 0x57, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x4f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x4e, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x38, 0x2e, 0x61, 0x63, 0x63, 0x75, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x6f,
	0x72, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x63, 0x63, 0x75, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x57, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x4f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x05, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x12, 0x30, 0x0a, 0x07, 0x77, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x61, 0x63, 0x63, 0x75, 0x6d, 0x75, 0x6c, 0x61,
	0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x57, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x52, 0x07, 0x77,
	0x69, 0x6e, 0x64, 0x6f, 0x77, 0x73, 0x22, 0x28, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12,
	0x08, 0x0a, 0x04, 0x4f, 0x50, 0x45, 0x4e, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x43, 0x4c, 0x4f,
	0x53, 0x45, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x41, 0x50, 0x50, 0x45, 0x4e, 0x44, 0x10, 0x02,
	0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x68, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x22, 0x7c,
	0x0a, 0x06, 0x57, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x12, 0x30, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12, 0x2c, 0x0a, 0x03, 0x65, 0x6e,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x03, 0x65, 0x6e, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6c, 0x6f, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73, 0x6c, 0x6f, 0x74, 0x22, 0xe6, 0x01, 0x0a,
	0x13, 0x41, 0x63, 0x63, 0x75, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x31, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x61, 0x63, 0x63, 0x75, 0x6d, 0x75, 0x6c, 0x61,
	0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x07,
	0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2e, 0x0a, 0x06, 0x77, 0x69, 0x6e, 0x64, 0x6f,
	0x77, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x61, 0x63, 0x63, 0x75, 0x6d, 0x75,
	0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x57, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x52,
	0x06, 0x77, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x12, 0x3c, 0x0a, 0x09, 0x68, 0x61, 0x6e, 0x64, 0x73,
	0x68, 0x61, 0x6b, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x61, 0x63, 0x63,
	0x75, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x61, 0x6e, 0x64,
	0x73, 0x68, 0x61, 0x6b, 0x65, 0x48, 0x00, 0x52, 0x09, 0x68, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61,
	0x6b, 0x65, 0x88, 0x01, 0x01, 0x12, 0x10, 0x0a, 0x03, 0x45, 0x4f, 0x46, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x03, 0x45, 0x4f, 0x46, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x68, 0x61, 0x6e, 0x64,
	0x73, 0x68, 0x61, 0x6b, 0x65, 0x22, 0x25, 0x0a, 0x0d, 0x52, 0x65, 0x61, 0x64, 0x79, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x65, 0x61, 0x64, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x72, 0x65, 0x61, 0x64, 0x79, 0x22, 0x1d, 0x0a, 0x09,
	0x48, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x6f, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x03, 0x73, 0x6f, 0x74, 0x32, 0xac, 0x01, 0x0a, 0x0b,
	0x41, 0x63, 0x63, 0x75, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x5b, 0x0a, 0x0c, 0x41,
	0x63, 0x63, 0x75, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x65, 0x46, 0x6e, 0x12, 0x22, 0x2e, 0x61, 0x63,
	0x63, 0x75, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x63, 0x63,
	0x75, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x23, 0x2e, 0x61, 0x63, 0x63, 0x75, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31,
	0x2e, 0x41, 0x63, 0x63, 0x75, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x30, 0x01, 0x12, 0x40, 0x0a, 0x07, 0x49, 0x73, 0x52, 0x65,
	0x61, 0x64, 0x79, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1d, 0x2e, 0x61, 0x63,
	0x63, 0x75, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x61,
	0x64, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x3f, 0x5a, 0x3d, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x75, 0x6d, 0x61, 0x70, 0x72, 0x6f,
	0x6a, 0x2f, 0x6e, 0x75, 0x6d, 0x61, 0x66, 0x6c, 0x6f, 0x77, 0x2d, 0x67, 0x6f, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x63, 0x63,
	0x75, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_pkg_apis_proto_accumulator_v1_accumulator_proto_rawDescOnce sync.Once
	file_pkg_apis_proto_accumulator_v1_accumulator_proto_rawDescData = file_pkg_apis_proto_accumulator_v1_accumulator_proto_rawDesc
)

func file_pkg_apis_proto_accumulator_v1_accumulator_proto_rawDescGZIP() []byte {
	file_pkg_apis_proto_accumulator_v1_accumulator_proto_rawDescOnce.Do(func() {
		file_pkg_apis_proto_accumulator_v1_accumulator_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_apis_proto_accumulator_v1_accumulator_proto_rawDescData)
	})
	return file_pkg_apis_proto_accumulator_v1_accumulator_proto_rawDescData
}

var file_pkg_apis_proto_accumulator_v1_accumulator_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pkg_apis_proto_accumulator_v1_accumulator_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_pkg_apis_proto_accumulator_v1_accumulator_proto_goTypes = []any{
	(AccumulatorRequest_WindowOperation_Event)(0), // 0: accumulator.v1.AccumulatorRequest.WindowOperation.Event
	(*Payload)(nil),             // 1: accumulator.v1.Payload
	(*AccumulatorRequest)(nil),  // 2: accumulator.v1.AccumulatorRequest
	(*Window)(nil),              // 3: accumulator.v1.Window
	(*AccumulatorResponse)(nil), // 4: accumulator.v1.AccumulatorResponse
	(*ReadyResponse)(nil),       // 5: accumulator.v1.ReadyResponse
	(*Handshake)(nil),           // 6: accumulator.v1.Handshake
	nil,                         // 7: accumulator.v1.Payload.HeadersEntry
	(*AccumulatorRequest_WindowOperation)(nil), // 8: accumulator.v1.AccumulatorRequest.WindowOperation
	(*timestamppb.Timestamp)(nil),              // 9: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),                      // 10: google.protobuf.Empty
}
var file_pkg_apis_proto_accumulator_v1_accumulator_proto_depIdxs = []int32{
	9,  // 0: accumulator.v1.Payload.event_time:type_name -> google.protobuf.Timestamp
	9,  // 1: accumulator.v1.Payload.watermark:type_name -> google.protobuf.Timestamp
	7,  // 2: accumulator.v1.Payload.headers:type_name -> accumulator.v1.Payload.HeadersEntry
	1,  // 3: accumulator.v1.AccumulatorRequest.payload:type_name -> accumulator.v1.Payload
	8,  // 4: accumulator.v1.AccumulatorRequest.operation:type_name -> accumulator.v1.AccumulatorRequest.WindowOperation
	6,  // 5: accumulator.v1.AccumulatorRequest.handshake:type_name -> accumulator.v1.Handshake
	9,  // 6: accumulator.v1.Window.start:type_name -> google.protobuf.Timestamp
	9,  // 7: accumulator.v1.Window.end:type_name -> google.protobuf.Timestamp
	1,  // 8: accumulator.v1.AccumulatorResponse.payload:type_name -> accumulator.v1.Payload
	3,  // 9: accumulator.v1.AccumulatorResponse.window:type_name -> accumulator.v1.Window
	6,  // 10: accumulator.v1.AccumulatorResponse.handshake:type_name -> accumulator.v1.Handshake
	0,  // 11: accumulator.v1.AccumulatorRequest.WindowOperation.event:type_name -> accumulator.v1.AccumulatorRequest.WindowOperation.Event
	3,  // 12: accumulator.v1.AccumulatorRequest.WindowOperation.windows:type_name -> accumulator.v1.Window
	2,  // 13: accumulator.v1.Accumulator.AccumulateFn:input_type -> accumulator.v1.AccumulatorRequest
	10, // 14: accumulator.v1.Accumulator.IsReady:input_type -> google.protobuf.Empty
	4,  // 15: accumulator.v1.Accumulator.AccumulateFn:output_type -> accumulator.v1.AccumulatorResponse
	5,  // 16: accumulator.v1.Accumulator.IsReady:output_type -> accumulator.v1.ReadyResponse
	15, // [15:17] is the sub-list for method output_type
	13, // [13:15] is the sub-list for method input_type
	13, // [13:13] is the sub-list for extension type_name
	13, // [13:13] is the sub-list for extension extendee
	0,  // [0:13] is the sub-list for field type_name
}

func init() { file_pkg_apis_proto_accumulator_v1_accumulator_proto_init() }
func file_pkg_apis_proto_accumulator_v1_accumulator_proto_init() {
	if File_pkg_apis_proto_accumulator_v1_accumulator_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_apis_proto_accumulator_v1_accumulator_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Payload); i {
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
		file_pkg_apis_proto_accumulator_v1_accumulator_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*AccumulatorRequest); i {
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
		file_pkg_apis_proto_accumulator_v1_accumulator_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*Window); i {
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
		file_pkg_apis_proto_accumulator_v1_accumulator_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*AccumulatorResponse); i {
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
		file_pkg_apis_proto_accumulator_v1_accumulator_proto_msgTypes[4].Exporter = func(v any, i int) any {
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
		file_pkg_apis_proto_accumulator_v1_accumulator_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*Handshake); i {
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
		file_pkg_apis_proto_accumulator_v1_accumulator_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*AccumulatorRequest_WindowOperation); i {
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
	file_pkg_apis_proto_accumulator_v1_accumulator_proto_msgTypes[1].OneofWrappers = []any{}
	file_pkg_apis_proto_accumulator_v1_accumulator_proto_msgTypes[3].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_apis_proto_accumulator_v1_accumulator_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_apis_proto_accumulator_v1_accumulator_proto_goTypes,
		DependencyIndexes: file_pkg_apis_proto_accumulator_v1_accumulator_proto_depIdxs,
		EnumInfos:         file_pkg_apis_proto_accumulator_v1_accumulator_proto_enumTypes,
		MessageInfos:      file_pkg_apis_proto_accumulator_v1_accumulator_proto_msgTypes,
	}.Build()
	File_pkg_apis_proto_accumulator_v1_accumulator_proto = out.File
	file_pkg_apis_proto_accumulator_v1_accumulator_proto_rawDesc = nil
	file_pkg_apis_proto_accumulator_v1_accumulator_proto_goTypes = nil
	file_pkg_apis_proto_accumulator_v1_accumulator_proto_depIdxs = nil
}
