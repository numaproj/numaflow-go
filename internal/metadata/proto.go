package metadata

import (
	"github.com/numaproj/numaflow-go/pkg/apis/proto/common"
)

// UserMetadataFromProto converts the incoming proto metadata to UserMetadata.
func UserMetadataFromProto(proto *common.Metadata) *UserMetadata {
	md := NewUserMetadata()
	if proto == nil {
		return md
	}
	for group, kvGroup := range proto.GetUserMetadata() {
		if kvGroup != nil {
			md.data[group] = kvGroup.GetKeyValue()
		} else {
			md.data[group] = make(map[string][]byte)
		}
	}
	return md
}

// SystemMetadataFromProto converts the incoming proto metadata to SystemMetadata.
func SystemMetadataFromProto(proto *common.Metadata) *SystemMetadata {
	md := NewSystemMetadata()
	if proto == nil {
		return md
	}
	for group, kvGroup := range proto.GetSysMetadata() {
		if kvGroup != nil {
			md.data[group] = kvGroup.GetKeyValue()
		} else {
			md.data[group] = make(map[string][]byte)
		}
	}
	return md
}

// UserMetadataToProto converts UserMetadata to the outgoing proto metadata.
// It always returns a non-nil *common.Metadata (with empty maps if md is nil).
func UserMetadataToProto(md *UserMetadata) *common.Metadata {
	sys := make(map[string]*common.KeyValueGroup)
	user := make(map[string]*common.KeyValueGroup)
	if md != nil {
		for _, group := range md.Groups() {
			kv := make(map[string][]byte)
			for _, key := range md.Keys(group) {
				kv[key] = md.Value(group, key)
			}
			user[group] = &common.KeyValueGroup{KeyValue: kv}
		}
	}
	return &common.Metadata{
		SysMetadata:  sys,
		UserMetadata: user,
	}
}
