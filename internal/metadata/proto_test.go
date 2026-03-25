package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/numaproj/numaflow-go/pkg/apis/proto/common"
)

func TestUserMetadataFromProto_Nil(t *testing.T) {
	md := UserMetadataFromProto(nil)
	assert.NotNil(t, md)
	assert.Empty(t, md.Groups())
}

func TestUserMetadataFromProto_Empty(t *testing.T) {
	proto := &common.Metadata{}
	md := UserMetadataFromProto(proto)
	assert.NotNil(t, md)
	assert.Empty(t, md.Groups())
}

func TestUserMetadataFromProto_WithData(t *testing.T) {
	proto := &common.Metadata{
		UserMetadata: map[string]*common.KeyValueGroup{
			"group1": {KeyValue: map[string][]byte{"k1": []byte("v1"), "k2": []byte("v2")}},
			"group2": {KeyValue: map[string][]byte{"k3": []byte("v3")}},
		},
	}
	md := UserMetadataFromProto(proto)
	assert.Len(t, md.Groups(), 2)
	assert.Equal(t, []byte("v1"), md.Value("group1", "k1"))
	assert.Equal(t, []byte("v2"), md.Value("group1", "k2"))
	assert.Equal(t, []byte("v3"), md.Value("group2", "k3"))
}

func TestUserMetadataFromProto_NilKeyValueGroup(t *testing.T) {
	proto := &common.Metadata{
		UserMetadata: map[string]*common.KeyValueGroup{
			"group1": nil,
		},
	}
	md := UserMetadataFromProto(proto)
	assert.Contains(t, md.Groups(), "group1")
	assert.Empty(t, md.Keys("group1"))
}

func TestSystemMetadataFromProto_Nil(t *testing.T) {
	md := SystemMetadataFromProto(nil)
	assert.NotNil(t, md)
	assert.Empty(t, md.Groups())
}

func TestSystemMetadataFromProto_WithData(t *testing.T) {
	proto := &common.Metadata{
		SysMetadata: map[string]*common.KeyValueGroup{
			"sys-group": {KeyValue: map[string][]byte{"sk": []byte("sv")}},
		},
	}
	md := SystemMetadataFromProto(proto)
	assert.Len(t, md.Groups(), 1)
	assert.Equal(t, []byte("sv"), md.Value("sys-group", "sk"))
}

func TestSystemMetadataFromProto_NilKeyValueGroup(t *testing.T) {
	proto := &common.Metadata{
		SysMetadata: map[string]*common.KeyValueGroup{
			"sys-group": nil,
		},
	}
	md := SystemMetadataFromProto(proto)
	assert.Contains(t, md.Groups(), "sys-group")
	assert.Empty(t, md.Keys("sys-group"))
}

func TestUserMetadataToProto_Nil(t *testing.T) {
	proto := UserMetadataToProto(nil)
	assert.NotNil(t, proto)
	assert.NotNil(t, proto.SysMetadata)
	assert.NotNil(t, proto.UserMetadata)
	assert.Empty(t, proto.SysMetadata)
	assert.Empty(t, proto.UserMetadata)
}

func TestUserMetadataToProto_Empty(t *testing.T) {
	md := NewUserMetadata()
	proto := UserMetadataToProto(md)
	assert.NotNil(t, proto)
	assert.Empty(t, proto.UserMetadata)
}

func TestUserMetadataToProto_WithData(t *testing.T) {
	md := NewUserMetadata()
	md.AddKV("group1", "k1", []byte("v1"))
	md.AddKV("group1", "k2", []byte("v2"))
	md.AddKV("group2", "k3", []byte("v3"))

	proto := UserMetadataToProto(md)
	assert.Len(t, proto.UserMetadata, 2)
	assert.Equal(t, []byte("v1"), proto.UserMetadata["group1"].KeyValue["k1"])
	assert.Equal(t, []byte("v2"), proto.UserMetadata["group1"].KeyValue["k2"])
	assert.Equal(t, []byte("v3"), proto.UserMetadata["group2"].KeyValue["k3"])
	assert.Empty(t, proto.SysMetadata)
}

func TestRoundTrip_UserMetadata(t *testing.T) {
	original := NewUserMetadata()
	original.AddKV("headers", "content-type", []byte("application/json"))
	original.AddKV("headers", "x-custom", []byte("value"))
	original.AddKV("metrics", "count", []byte("42"))

	proto := UserMetadataToProto(original)
	restored := UserMetadataFromProto(proto)

	assert.ElementsMatch(t, original.Groups(), restored.Groups())
	for _, group := range original.Groups() {
		assert.ElementsMatch(t, original.Keys(group), restored.Keys(group))
		for _, key := range original.Keys(group) {
			assert.Equal(t, original.Value(group, key), restored.Value(group, key))
		}
	}
}
