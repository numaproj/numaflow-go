// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.29.3
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

// Status is the status of the response.
type Status int32

const (
	Status_SUCCESS  Status = 0
	Status_FAILURE  Status = 1
	Status_FALLBACK Status = 2
	Status_SERVE    Status = 3
)

// Enum value maps for Status.
var (
	Status_name = map[int32]string{
		0: "SUCCESS",
		1: "FAILURE",
		2: "FALLBACK",
		3: "SERVE",
	}
	Status_value = map[string]int32{
		"SUCCESS":  0,
		"FAILURE":  1,
		"FALLBACK": 2,
		"SERVE":    3,
	}
)

func (x Status) Enum() *Status {
	p := new(Status)
	*p = x
	return p
}

func (x Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Status) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_apis_proto_sink_v1_sink_proto_enumTypes[0].Descriptor()
}

func (Status) Type() protoreflect.EnumType {
	return &file_pkg_apis_proto_sink_v1_sink_proto_enumTypes[0]
}

func (x Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Status.Descriptor instead.
func (Status) EnumDescriptor() ([]byte, []int) {
	return file_pkg_apis_proto_sink_v1_sink_proto_rawDescGZIP(), []int{0}
}

// *
// SinkRequest represents a request element.
type SinkRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Required field indicating the request.
	Request *SinkRequest_Request `protobuf:"bytes,1,opt,name=request,proto3" json:"request,omitempty"`
	// Required field indicating the status of the request.
	// If eot is set to true, it indicates the end of transmission.
	Status *TransmissionStatus `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	// optional field indicating the handshake message.
	Handshake *Handshake `protobuf:"bytes,3,opt,name=handshake,proto3,oneof" json:"handshake,omitempty"`
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

func (x *SinkRequest) GetRequest() *SinkRequest_Request {
	if x != nil {
		return x.Request
	}
	return nil
}

func (x *SinkRequest) GetStatus() *TransmissionStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

func (x *SinkRequest) GetHandshake() *Handshake {
	if x != nil {
		return x.Handshake
	}
	return nil
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
		mi := &file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Handshake) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Handshake) ProtoMessage() {}

func (x *Handshake) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Handshake.ProtoReflect.Descriptor instead.
func (*Handshake) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_sink_v1_sink_proto_rawDescGZIP(), []int{1}
}

func (x *Handshake) GetSot() bool {
	if x != nil {
		return x.Sot
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
		mi := &file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadyResponse) ProtoMessage() {}

func (x *ReadyResponse) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use ReadyResponse.ProtoReflect.Descriptor instead.
func (*ReadyResponse) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_sink_v1_sink_proto_rawDescGZIP(), []int{2}
}

func (x *ReadyResponse) GetReady() bool {
	if x != nil {
		return x.Ready
	}
	return false
}

// *
// TransmissionStatus is the status of the transmission.
type TransmissionStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Eot bool `protobuf:"varint,1,opt,name=eot,proto3" json:"eot,omitempty"`
}

func (x *TransmissionStatus) Reset() {
	*x = TransmissionStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransmissionStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransmissionStatus) ProtoMessage() {}

func (x *TransmissionStatus) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use TransmissionStatus.ProtoReflect.Descriptor instead.
func (*TransmissionStatus) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_sink_v1_sink_proto_rawDescGZIP(), []int{3}
}

func (x *TransmissionStatus) GetEot() bool {
	if x != nil {
		return x.Eot
	}
	return false
}

// *
// SinkResponse is the individual response of each message written to the sink.
type SinkResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Results   []*SinkResponse_Result `protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty"`
	Handshake *Handshake             `protobuf:"bytes,2,opt,name=handshake,proto3,oneof" json:"handshake,omitempty"`
	Status    *TransmissionStatus    `protobuf:"bytes,3,opt,name=status,proto3,oneof" json:"status,omitempty"`
}

func (x *SinkResponse) Reset() {
	*x = SinkResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SinkResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SinkResponse) ProtoMessage() {}

func (x *SinkResponse) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use SinkResponse.ProtoReflect.Descriptor instead.
func (*SinkResponse) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_sink_v1_sink_proto_rawDescGZIP(), []int{4}
}

func (x *SinkResponse) GetResults() []*SinkResponse_Result {
	if x != nil {
		return x.Results
	}
	return nil
}

func (x *SinkResponse) GetHandshake() *Handshake {
	if x != nil {
		return x.Handshake
	}
	return nil
}

func (x *SinkResponse) GetStatus() *TransmissionStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

type SinkRequest_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keys      []string               `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
	Value     []byte                 `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	EventTime *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=event_time,json=eventTime,proto3" json:"event_time,omitempty"`
	Watermark *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=watermark,proto3" json:"watermark,omitempty"`
	Id        string                 `protobuf:"bytes,5,opt,name=id,proto3" json:"id,omitempty"`
	Headers   map[string]string      `protobuf:"bytes,6,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *SinkRequest_Request) Reset() {
	*x = SinkRequest_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SinkRequest_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SinkRequest_Request) ProtoMessage() {}

func (x *SinkRequest_Request) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use SinkRequest_Request.ProtoReflect.Descriptor instead.
func (*SinkRequest_Request) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_sink_v1_sink_proto_rawDescGZIP(), []int{0, 0}
}

func (x *SinkRequest_Request) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *SinkRequest_Request) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *SinkRequest_Request) GetEventTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EventTime
	}
	return nil
}

func (x *SinkRequest_Request) GetWatermark() *timestamppb.Timestamp {
	if x != nil {
		return x.Watermark
	}
	return nil
}

func (x *SinkRequest_Request) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *SinkRequest_Request) GetHeaders() map[string]string {
	if x != nil {
		return x.Headers
	}
	return nil
}

type SinkResponse_Result struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id is the ID of the message, can be used to uniquely identify the message.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// status denotes the status of persisting to sink. It can be SUCCESS, FAILURE, or FALLBACK.
	Status Status `protobuf:"varint,2,opt,name=status,proto3,enum=sink.v1.Status" json:"status,omitempty"`
	// err_msg is the error message, set it if success is set to false.
	ErrMsg        string `protobuf:"bytes,3,opt,name=err_msg,json=errMsg,proto3" json:"err_msg,omitempty"`
	ServeResponse []byte `protobuf:"bytes,4,opt,name=serve_response,json=serveResponse,proto3,oneof" json:"serve_response,omitempty"`
}

func (x *SinkResponse_Result) Reset() {
	*x = SinkResponse_Result{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SinkResponse_Result) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SinkResponse_Result) ProtoMessage() {}

func (x *SinkResponse_Result) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SinkResponse_Result.ProtoReflect.Descriptor instead.
func (*SinkResponse_Result) Descriptor() ([]byte, []int) {
	return file_pkg_apis_proto_sink_v1_sink_proto_rawDescGZIP(), []int{4, 0}
}

func (x *SinkResponse_Result) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *SinkResponse_Result) GetStatus() Status {
	if x != nil {
		return x.Status
	}
	return Status_SUCCESS
}

func (x *SinkResponse_Result) GetErrMsg() string {
	if x != nil {
		return x.ErrMsg
	}
	return ""
}

func (x *SinkResponse_Result) GetServeResponse() []byte {
	if x != nil {
		return x.ServeResponse
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
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xfb, 0x03, 0x0a, 0x0b, 0x53,
	0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x36, 0x0a, 0x07, 0x72, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x73, 0x69,
	0x6e, 0x6b, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x33, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x35, 0x0a, 0x09, 0x68, 0x61, 0x6e, 0x64, 0x73,
	0x68, 0x61, 0x6b, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x73, 0x69, 0x6e,
	0x6b, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x48, 0x00,
	0x52, 0x09, 0x68, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x88, 0x01, 0x01, 0x1a, 0xb9,
	0x02, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65,
	0x79, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x69,
	0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12,
	0x38, 0x0a, 0x09, 0x77, 0x61, 0x74, 0x65, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09,
	0x77, 0x61, 0x74, 0x65, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x43, 0x0a, 0x07, 0x68, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x73, 0x69, 0x6e,
	0x6b, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x1a, 0x3a,
	0x0a, 0x0c, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x68,
	0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x22, 0x1d, 0x0a, 0x09, 0x48, 0x61, 0x6e, 0x64,
	0x73, 0x68, 0x61, 0x6b, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x6f, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x03, 0x73, 0x6f, 0x74, 0x22, 0x25, 0x0a, 0x0d, 0x52, 0x65, 0x61, 0x64, 0x79,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x65, 0x61, 0x64,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x72, 0x65, 0x61, 0x64, 0x79, 0x22, 0x26,
	0x0a, 0x12, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x6f, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x03, 0x65, 0x6f, 0x74, 0x22, 0xec, 0x02, 0x0a, 0x0c, 0x53, 0x69, 0x6e, 0x6b, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x73, 0x69, 0x6e, 0x6b, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x12,
	0x35, 0x0a, 0x09, 0x68, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x12, 0x2e, 0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x61, 0x6e,
	0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x48, 0x00, 0x52, 0x09, 0x68, 0x61, 0x6e, 0x64, 0x73, 0x68,
	0x61, 0x6b, 0x65, 0x88, 0x01, 0x01, 0x12, 0x38, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x76, 0x31,
	0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x48, 0x01, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x88, 0x01, 0x01,
	0x1a, 0x99, 0x01, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x27, 0x0a, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x73, 0x69,
	0x6e, 0x6b, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x17, 0x0a, 0x07, 0x65, 0x72, 0x72, 0x5f, 0x6d, 0x73, 0x67, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x65, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x12, 0x2a, 0x0a,
	0x0e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x0d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x88, 0x01, 0x01, 0x42, 0x11, 0x0a, 0x0f, 0x5f, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0c, 0x0a, 0x0a,
	0x5f, 0x68, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x2a, 0x3b, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x0b, 0x0a, 0x07, 0x53, 0x55, 0x43, 0x43, 0x45, 0x53, 0x53, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07,
	0x46, 0x41, 0x49, 0x4c, 0x55, 0x52, 0x45, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x46, 0x41, 0x4c,
	0x4c, 0x42, 0x41, 0x43, 0x4b, 0x10, 0x02, 0x12, 0x09, 0x0a, 0x05, 0x53, 0x45, 0x52, 0x56, 0x45,
	0x10, 0x03, 0x32, 0x7c, 0x0a, 0x04, 0x53, 0x69, 0x6e, 0x6b, 0x12, 0x39, 0x0a, 0x06, 0x53, 0x69,
	0x6e, 0x6b, 0x46, 0x6e, 0x12, 0x14, 0x2e, 0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x76, 0x31, 0x2e, 0x53,
	0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x73, 0x69, 0x6e,
	0x6b, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x28, 0x01, 0x30, 0x01, 0x12, 0x39, 0x0a, 0x07, 0x49, 0x73, 0x52, 0x65, 0x61, 0x64, 0x79,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x73, 0x69, 0x6e, 0x6b, 0x2e,
	0x76, 0x31, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x42, 0x38, 0x5a, 0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e,
	0x75, 0x6d, 0x61, 0x70, 0x72, 0x6f, 0x6a, 0x2f, 0x6e, 0x75, 0x6d, 0x61, 0x66, 0x6c, 0x6f, 0x77,
	0x2d, 0x67, 0x6f, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x73, 0x69, 0x6e, 0x6b, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
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

var file_pkg_apis_proto_sink_v1_sink_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pkg_apis_proto_sink_v1_sink_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_pkg_apis_proto_sink_v1_sink_proto_goTypes = []any{
	(Status)(0),                   // 0: sink.v1.Status
	(*SinkRequest)(nil),           // 1: sink.v1.SinkRequest
	(*Handshake)(nil),             // 2: sink.v1.Handshake
	(*ReadyResponse)(nil),         // 3: sink.v1.ReadyResponse
	(*TransmissionStatus)(nil),    // 4: sink.v1.TransmissionStatus
	(*SinkResponse)(nil),          // 5: sink.v1.SinkResponse
	(*SinkRequest_Request)(nil),   // 6: sink.v1.SinkRequest.Request
	nil,                           // 7: sink.v1.SinkRequest.Request.HeadersEntry
	(*SinkResponse_Result)(nil),   // 8: sink.v1.SinkResponse.Result
	(*timestamppb.Timestamp)(nil), // 9: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),         // 10: google.protobuf.Empty
}
var file_pkg_apis_proto_sink_v1_sink_proto_depIdxs = []int32{
	6,  // 0: sink.v1.SinkRequest.request:type_name -> sink.v1.SinkRequest.Request
	4,  // 1: sink.v1.SinkRequest.status:type_name -> sink.v1.TransmissionStatus
	2,  // 2: sink.v1.SinkRequest.handshake:type_name -> sink.v1.Handshake
	8,  // 3: sink.v1.SinkResponse.results:type_name -> sink.v1.SinkResponse.Result
	2,  // 4: sink.v1.SinkResponse.handshake:type_name -> sink.v1.Handshake
	4,  // 5: sink.v1.SinkResponse.status:type_name -> sink.v1.TransmissionStatus
	9,  // 6: sink.v1.SinkRequest.Request.event_time:type_name -> google.protobuf.Timestamp
	9,  // 7: sink.v1.SinkRequest.Request.watermark:type_name -> google.protobuf.Timestamp
	7,  // 8: sink.v1.SinkRequest.Request.headers:type_name -> sink.v1.SinkRequest.Request.HeadersEntry
	0,  // 9: sink.v1.SinkResponse.Result.status:type_name -> sink.v1.Status
	1,  // 10: sink.v1.Sink.SinkFn:input_type -> sink.v1.SinkRequest
	10, // 11: sink.v1.Sink.IsReady:input_type -> google.protobuf.Empty
	5,  // 12: sink.v1.Sink.SinkFn:output_type -> sink.v1.SinkResponse
	3,  // 13: sink.v1.Sink.IsReady:output_type -> sink.v1.ReadyResponse
	12, // [12:14] is the sub-list for method output_type
	10, // [10:12] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_pkg_apis_proto_sink_v1_sink_proto_init() }
func file_pkg_apis_proto_sink_v1_sink_proto_init() {
	if File_pkg_apis_proto_sink_v1_sink_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[0].Exporter = func(v any, i int) any {
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
		file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[1].Exporter = func(v any, i int) any {
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
		file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[2].Exporter = func(v any, i int) any {
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
		file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*TransmissionStatus); i {
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
		file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[4].Exporter = func(v any, i int) any {
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
		file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*SinkRequest_Request); i {
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
		file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*SinkResponse_Result); i {
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
	file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[0].OneofWrappers = []any{}
	file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[4].OneofWrappers = []any{}
	file_pkg_apis_proto_sink_v1_sink_proto_msgTypes[7].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_apis_proto_sink_v1_sink_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_apis_proto_sink_v1_sink_proto_goTypes,
		DependencyIndexes: file_pkg_apis_proto_sink_v1_sink_proto_depIdxs,
		EnumInfos:         file_pkg_apis_proto_sink_v1_sink_proto_enumTypes,
		MessageInfos:      file_pkg_apis_proto_sink_v1_sink_proto_msgTypes,
	}.Build()
	File_pkg_apis_proto_sink_v1_sink_proto = out.File
	file_pkg_apis_proto_sink_v1_sink_proto_rawDesc = nil
	file_pkg_apis_proto_sink_v1_sink_proto_goTypes = nil
	file_pkg_apis_proto_sink_v1_sink_proto_depIdxs = nil
}
