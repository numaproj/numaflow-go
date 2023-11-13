// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: pkg/apis/proto/globalreduce/v1/globalreduce.proto

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

type GlobalReduceRequest_WindowOperation_Event int32

const (
	GlobalReduceRequest_WindowOperation_OPEN   GlobalReduceRequest_WindowOperation_Event = 0
	GlobalReduceRequest_WindowOperation_CLOSE  GlobalReduceRequest_WindowOperation_Event = 1
	GlobalReduceRequest_WindowOperation_EXPAND GlobalReduceRequest_WindowOperation_Event = 2
	GlobalReduceRequest_WindowOperation_MERGE  GlobalReduceRequest_WindowOperation_Event = 3
	GlobalReduceRequest_WindowOperation_APPEND GlobalReduceRequest_WindowOperation_Event = 4
)

// Enum value maps for GlobalReduceRequest_WindowOperation_Event.
var (
	GlobalReduceRequest_WindowOperation_Event_name = map[int32]string{
		0: "OPEN",
		1: "CLOSE",
		2: "EXPAND",
		3: "MERGE",
		4: "APPEND",
	}
	GlobalReduceRequest_WindowOperation_Event_value = map[string]int32{
		"OPEN":   0,
		"CLOSE":  1,
		"EXPAND": 2,
		"MERGE":  3,
		"APPEND": 4,
	}
)

func (x GlobalReduceRequest_WindowOperation_Event) Enum() *GlobalReduceRequest_WindowOperation_Event {
	p := new(GlobalReduceRequest_WindowOperation_Event)
	*p = x
	return p
}

func (x GlobalReduceRequest_WindowOperation_Event) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GlobalReduceRequest_WindowOperation_Event) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_enumTypes[0].Descriptor()
}

func (GlobalReduceRequest_WindowOperation_Event) Type() protoreflect.EnumType {
	return &file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_enumTypes[0]
}

func (x GlobalReduceRequest_WindowOperation_Event) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GlobalReduceRequest_WindowOperation_Event.Descriptor instead.
func (GlobalReduceRequest_WindowOperation_Event) EnumDescriptor() ([]byte, []int) {
	return file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_rawDescGZIP(), []int{1, 0, 0}
}

// Partition represents a window partition.
type Partition struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Start *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=start,proto3" json:"start,omitempty"`
	End   *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=end,proto3" json:"end,omitempty"`
	Slot  string                 `protobuf:"bytes,3,opt,name=slot,proto3" json:"slot,omitempty"`
}

func (x *Partition) Reset() {
	*x = Partition{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Partition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Partition) ProtoMessage() {}

func (x *Partition) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Partition.ProtoReflect.Descriptor instead.
func (*Partition) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_rawDescGZIP(), []int{0}
}

func (x *Partition) GetStart() *timestamppb.Timestamp {
	if x != nil {
		return x.Start
	}
	return nil
}

func (x *Partition) GetEnd() *timestamppb.Timestamp {
	if x != nil {
		return x.End
	}
	return nil
}

func (x *Partition) GetSlot() string {
	if x != nil {
		return x.Slot
	}
	return ""
}

// *
// GlobalReduceRequest represents a request element.
type GlobalReduceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Payload   *GlobalReduceRequest_Payload         `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	Operation *GlobalReduceRequest_WindowOperation `protobuf:"bytes,2,opt,name=operation,proto3" json:"operation,omitempty"`
}

func (x *GlobalReduceRequest) Reset() {
	*x = GlobalReduceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GlobalReduceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GlobalReduceRequest) ProtoMessage() {}

func (x *GlobalReduceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GlobalReduceRequest.ProtoReflect.Descriptor instead.
func (*GlobalReduceRequest) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_rawDescGZIP(), []int{1}
}

func (x *GlobalReduceRequest) GetPayload() *GlobalReduceRequest_Payload {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *GlobalReduceRequest) GetOperation() *GlobalReduceRequest_WindowOperation {
	if x != nil {
		return x.Operation
	}
	return nil
}

// *
// GlobalReduceResponse represents a response element.
type GlobalReduceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Results []*GlobalReduceResponse_Result `protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty"`
	// Partition represents a window partition to which the result belongs.
	Partition *Partition `protobuf:"bytes,2,opt,name=partition,proto3" json:"partition,omitempty"`
	// the key which was used for demultiplexing the request stream.
	CombinedKey string `protobuf:"bytes,3,opt,name=combinedKey,proto3" json:"combinedKey,omitempty"`
	// EventTime represents the event time of the result.
	EventTime *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=event_time,json=eventTime,proto3" json:"event_time,omitempty"`
}

func (x *GlobalReduceResponse) Reset() {
	*x = GlobalReduceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GlobalReduceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GlobalReduceResponse) ProtoMessage() {}

func (x *GlobalReduceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GlobalReduceResponse.ProtoReflect.Descriptor instead.
func (*GlobalReduceResponse) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_rawDescGZIP(), []int{2}
}

func (x *GlobalReduceResponse) GetResults() []*GlobalReduceResponse_Result {
	if x != nil {
		return x.Results
	}
	return nil
}

func (x *GlobalReduceResponse) GetPartition() *Partition {
	if x != nil {
		return x.Partition
	}
	return nil
}

func (x *GlobalReduceResponse) GetCombinedKey() string {
	if x != nil {
		return x.CombinedKey
	}
	return ""
}

func (x *GlobalReduceResponse) GetEventTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EventTime
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
		mi := &file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadyResponse) ProtoMessage() {}

func (x *ReadyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_msgTypes[3]
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
	return file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_rawDescGZIP(), []int{3}
}

func (x *ReadyResponse) GetReady() bool {
	if x != nil {
		return x.Ready
	}
	return false
}

// WindowOperation represents a window operation.
// it can be one of OPEN, CLOSE, EXPAND, MERGE, APPEND.
type GlobalReduceRequest_WindowOperation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Event      GlobalReduceRequest_WindowOperation_Event `protobuf:"varint,1,opt,name=event,proto3,enum=globalreduce.v1.GlobalReduceRequest_WindowOperation_Event" json:"event,omitempty"`
	Partitions []*Partition                              `protobuf:"bytes,2,rep,name=partitions,proto3" json:"partitions,omitempty"`
}

func (x *GlobalReduceRequest_WindowOperation) Reset() {
	*x = GlobalReduceRequest_WindowOperation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GlobalReduceRequest_WindowOperation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GlobalReduceRequest_WindowOperation) ProtoMessage() {}

func (x *GlobalReduceRequest_WindowOperation) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GlobalReduceRequest_WindowOperation.ProtoReflect.Descriptor instead.
func (*GlobalReduceRequest_WindowOperation) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_rawDescGZIP(), []int{1, 0}
}

func (x *GlobalReduceRequest_WindowOperation) GetEvent() GlobalReduceRequest_WindowOperation_Event {
	if x != nil {
		return x.Event
	}
	return GlobalReduceRequest_WindowOperation_OPEN
}

func (x *GlobalReduceRequest_WindowOperation) GetPartitions() []*Partition {
	if x != nil {
		return x.Partitions
	}
	return nil
}

// Payload represents a payload element.
type GlobalReduceRequest_Payload struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keys      []string               `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
	Value     []byte                 `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	EventTime *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=event_time,json=eventTime,proto3" json:"event_time,omitempty"`
	Watermark *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=watermark,proto3" json:"watermark,omitempty"`
}

func (x *GlobalReduceRequest_Payload) Reset() {
	*x = GlobalReduceRequest_Payload{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GlobalReduceRequest_Payload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GlobalReduceRequest_Payload) ProtoMessage() {}

func (x *GlobalReduceRequest_Payload) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GlobalReduceRequest_Payload.ProtoReflect.Descriptor instead.
func (*GlobalReduceRequest_Payload) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_rawDescGZIP(), []int{1, 1}
}

func (x *GlobalReduceRequest_Payload) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *GlobalReduceRequest_Payload) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *GlobalReduceRequest_Payload) GetEventTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EventTime
	}
	return nil
}

func (x *GlobalReduceRequest_Payload) GetWatermark() *timestamppb.Timestamp {
	if x != nil {
		return x.Watermark
	}
	return nil
}

// Result represents a result element. It contains the result of the reduce function.
type GlobalReduceResponse_Result struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keys  []string `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
	Value []byte   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Tags  []string `protobuf:"bytes,3,rep,name=tags,proto3" json:"tags,omitempty"`
}

func (x *GlobalReduceResponse_Result) Reset() {
	*x = GlobalReduceResponse_Result{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GlobalReduceResponse_Result) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GlobalReduceResponse_Result) ProtoMessage() {}

func (x *GlobalReduceResponse_Result) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GlobalReduceResponse_Result.ProtoReflect.Descriptor instead.
func (*GlobalReduceResponse_Result) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_rawDescGZIP(), []int{2, 0}
}

func (x *GlobalReduceResponse_Result) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *GlobalReduceResponse_Result) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *GlobalReduceResponse_Result) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

var File_pkg_apis_proto_globalreduce_v1_globalreduce_proto protoreflect.FileDescriptor

var file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_rawDesc = []byte{
	0x0a, 0x31, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x2f, 0x76, 0x31,
	0x2f, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x72, 0x65, 0x64, 0x75, 0x63,
	0x65, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x7f, 0x0a, 0x09, 0x50, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x30, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x12, 0x2c, 0x0a, 0x03, 0x65, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x03, 0x65, 0x6e, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x73, 0x6c, 0x6f, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73,
	0x6c, 0x6f, 0x74, 0x22, 0xbf, 0x04, 0x0a, 0x13, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x52, 0x65,
	0x64, 0x75, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x46, 0x0a, 0x07, 0x70,
	0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x67,
	0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47,
	0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x2e, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c,
	0x6f, 0x61, 0x64, 0x12, 0x52, 0x0a, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x34, 0x2e, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x72,
	0x65, 0x64, 0x75, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x52,
	0x65, 0x64, 0x75, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x57, 0x69, 0x6e,
	0x64, 0x6f, 0x77, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x6f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0xe0, 0x01, 0x0a, 0x0f, 0x57, 0x69, 0x6e, 0x64,
	0x6f, 0x77, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x50, 0x0a, 0x05, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x3a, 0x2e, 0x67, 0x6c, 0x6f,
	0x62, 0x61, 0x6c, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x6c, 0x6f,
	0x62, 0x61, 0x6c, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x2e, 0x57, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x3a, 0x0a,
	0x0a, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x70,
	0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x3f, 0x0a, 0x05, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x12, 0x08, 0x0a, 0x04, 0x4f, 0x50, 0x45, 0x4e, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05,
	0x43, 0x4c, 0x4f, 0x53, 0x45, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x45, 0x58, 0x50, 0x41, 0x4e,
	0x44, 0x10, 0x02, 0x12, 0x09, 0x0a, 0x05, 0x4d, 0x45, 0x52, 0x47, 0x45, 0x10, 0x03, 0x12, 0x0a,
	0x0a, 0x06, 0x41, 0x50, 0x50, 0x45, 0x4e, 0x44, 0x10, 0x04, 0x1a, 0xa8, 0x01, 0x0a, 0x07, 0x50,
	0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x12, 0x39, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x38, 0x0a, 0x09, 0x77,
	0x61, 0x74, 0x65, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x77, 0x61, 0x74, 0x65,
	0x72, 0x6d, 0x61, 0x72, 0x6b, 0x22, 0xbd, 0x02, 0x0a, 0x14, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c,
	0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x46,
	0x0a, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x2c, 0x2e, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x07, 0x72,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x12, 0x38, 0x0a, 0x09, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6c, 0x6f, 0x62,
	0x61, 0x6c, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x72, 0x74,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x20, 0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x62, 0x69, 0x6e, 0x65, 0x64, 0x4b, 0x65, 0x79, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x6f, 0x6d, 0x62, 0x69, 0x6e, 0x65, 0x64, 0x4b,
	0x65, 0x79, 0x12, 0x39, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x1a, 0x46, 0x0a,
	0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x04, 0x74, 0x61, 0x67, 0x73, 0x22, 0x25, 0x0a, 0x0d, 0x52, 0x65, 0x61, 0x64, 0x79, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x65, 0x61, 0x64, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x72, 0x65, 0x61, 0x64, 0x79, 0x32, 0xb4, 0x01, 0x0a,
	0x0c, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x12, 0x61, 0x0a,
	0x0e, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x46, 0x6e, 0x12,
	0x24, 0x2e, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x72, 0x65,
	0x64, 0x75, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x52, 0x65,
	0x64, 0x75, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x30, 0x01,
	0x12, 0x41, 0x0a, 0x07, 0x49, 0x73, 0x52, 0x65, 0x61, 0x64, 0x79, 0x12, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x1a, 0x1e, 0x2e, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x72, 0x65, 0x64, 0x75,
	0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x40, 0x5a, 0x3e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x6e, 0x75, 0x6d, 0x61, 0x70, 0x72, 0x6f, 0x6a, 0x2f, 0x6e, 0x75, 0x6d, 0x61, 0x66,
	0x6c, 0x6f, 0x77, 0x2d, 0x67, 0x6f, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x73, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_rawDescOnce sync.Once
	file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_rawDescData = file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_rawDesc
)

func file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_rawDescGZIP() []byte {
	file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_rawDescOnce.Do(func() {
		file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_rawDescData)
	})
	return file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_rawDescData
}

var file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_goTypes = []interface{}{
	(GlobalReduceRequest_WindowOperation_Event)(0), // 0: globalreduce.v1.GlobalReduceRequest.WindowOperation.Operation
	(*Partition)(nil),                           // 1: globalreduce.v1.Partition
	(*GlobalReduceRequest)(nil),                 // 2: globalreduce.v1.GlobalReduceRequest
	(*GlobalReduceResponse)(nil),                // 3: globalreduce.v1.GlobalReduceResponse
	(*ReadyResponse)(nil),                       // 4: globalreduce.v1.ReadyResponse
	(*GlobalReduceRequest_WindowOperation)(nil), // 5: globalreduce.v1.GlobalReduceRequest.WindowOperation
	(*GlobalReduceRequest_Payload)(nil),         // 6: globalreduce.v1.GlobalReduceRequest.Payload
	(*GlobalReduceResponse_Result)(nil),         // 7: globalreduce.v1.GlobalReduceResponse.Result
	(*timestamppb.Timestamp)(nil),               // 8: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),                       // 9: google.protobuf.Empty
}
var file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_depIdxs = []int32{
	8,  // 0: globalreduce.v1.Partition.start:type_name -> google.protobuf.Timestamp
	8,  // 1: globalreduce.v1.Partition.end:type_name -> google.protobuf.Timestamp
	6,  // 2: globalreduce.v1.GlobalReduceRequest.payload:type_name -> globalreduce.v1.GlobalReduceRequest.Payload
	5,  // 3: globalreduce.v1.GlobalReduceRequest.operation:type_name -> globalreduce.v1.GlobalReduceRequest.WindowOperation
	7,  // 4: globalreduce.v1.GlobalReduceResponse.results:type_name -> globalreduce.v1.GlobalReduceResponse.Result
	1,  // 5: globalreduce.v1.GlobalReduceResponse.partition:type_name -> globalreduce.v1.Partition
	8,  // 6: globalreduce.v1.GlobalReduceResponse.event_time:type_name -> google.protobuf.Timestamp
	0,  // 7: globalreduce.v1.GlobalReduceRequest.WindowOperation.event:type_name -> globalreduce.v1.GlobalReduceRequest.WindowOperation.Operation
	1,  // 8: globalreduce.v1.GlobalReduceRequest.WindowOperation.partitions:type_name -> globalreduce.v1.Partition
	8,  // 9: globalreduce.v1.GlobalReduceRequest.Payload.event_time:type_name -> google.protobuf.Timestamp
	8,  // 10: globalreduce.v1.GlobalReduceRequest.Payload.watermark:type_name -> google.protobuf.Timestamp
	2,  // 11: globalreduce.v1.GlobalReduce.GlobalReduceFn:input_type -> globalreduce.v1.GlobalReduceRequest
	9,  // 12: globalreduce.v1.GlobalReduce.IsReady:input_type -> google.protobuf.Empty
	3,  // 13: globalreduce.v1.GlobalReduce.GlobalReduceFn:output_type -> globalreduce.v1.GlobalReduceResponse
	4,  // 14: globalreduce.v1.GlobalReduce.IsReady:output_type -> globalreduce.v1.ReadyResponse
	13, // [13:15] is the sub-list for method output_type
	11, // [11:13] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_init() }
func file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_init() {
	if File_pkg_apis_proto_globalreduce_v1_globalreduce_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Partition); i {
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
		file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GlobalReduceRequest); i {
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
		file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GlobalReduceResponse); i {
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
		file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
		file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GlobalReduceRequest_WindowOperation); i {
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
		file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GlobalReduceRequest_Payload); i {
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
		file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GlobalReduceResponse_Result); i {
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
			RawDescriptor: file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_goTypes,
		DependencyIndexes: file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_depIdxs,
		EnumInfos:         file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_enumTypes,
		MessageInfos:      file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_msgTypes,
	}.Build()
	File_pkg_apis_proto_globalreduce_v1_globalreduce_proto = out.File
	file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_rawDesc = nil
	file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_goTypes = nil
	file_pkg_apis_proto_globalreduce_v1_globalreduce_proto_depIdxs = nil
}
