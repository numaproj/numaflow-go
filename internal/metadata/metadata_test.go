package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemMetadata_New(t *testing.T) {
	md := NewSystemMetadata()
	assert.NotNil(t, md)
	assert.Empty(t, md.Groups())
}

func TestSystemMetadata_NilSafety(t *testing.T) {
	var md *SystemMetadata
	assert.Empty(t, md.Groups())
	assert.Empty(t, md.Keys("any"))
	assert.Empty(t, md.Value("any", "key"))
}

func TestSystemMetadata_GroupsKeysValue(t *testing.T) {
	md := NewSystemMetadata()
	md.data["group1"] = map[string][]byte{
		"key1": []byte("val1"),
		"key2": []byte("val2"),
	}
	md.data["group2"] = map[string][]byte{
		"key3": []byte("val3"),
	}

	groups := md.Groups()
	assert.Len(t, groups, 2)
	assert.ElementsMatch(t, []string{"group1", "group2"}, groups)

	keys := md.Keys("group1")
	assert.Len(t, keys, 2)
	assert.ElementsMatch(t, []string{"key1", "key2"}, keys)

	assert.Equal(t, []byte("val1"), md.Value("group1", "key1"))
	assert.Empty(t, md.Keys("nonexistent"))
	assert.Empty(t, md.Value("group1", "nonexistent"))
}

func TestUserMetadata_New(t *testing.T) {
	md := NewUserMetadata()
	assert.NotNil(t, md)
	assert.Empty(t, md.Groups())
}

func TestUserMetadata_NilSafety(t *testing.T) {
	var md *UserMetadata
	assert.Empty(t, md.Groups())
	assert.Empty(t, md.Keys("any"))
	assert.Empty(t, md.Value("any", "key"))
	// Write methods on nil should not panic
	md.RemoveKey("any", "key")
	md.RemoveGroup("any")
}

func TestUserMetadata_GroupsKeysValue(t *testing.T) {
	md := NewUserMetadata()
	md.AddKV("group1", "key1", []byte("val1"))
	md.AddKV("group1", "key2", []byte("val2"))
	md.AddKV("group2", "key3", []byte("val3"))

	groups := md.Groups()
	assert.Len(t, groups, 2)
	assert.ElementsMatch(t, []string{"group1", "group2"}, groups)

	keys := md.Keys("group1")
	assert.Len(t, keys, 2)
	assert.ElementsMatch(t, []string{"key1", "key2"}, keys)

	assert.Equal(t, []byte("val1"), md.Value("group1", "key1"))
}

func TestUserMetadata_CreateGroup(t *testing.T) {
	md := NewUserMetadata()
	md.CreateGroup("mygroup")
	assert.Contains(t, md.Groups(), "mygroup")
	assert.Empty(t, md.Keys("mygroup"))
}

func TestUserMetadata_AddKV(t *testing.T) {
	md := NewUserMetadata()
	md.AddKV("g", "k", []byte("v"))
	assert.Equal(t, []byte("v"), md.Value("g", "k"))

	// Overwrite
	md.AddKV("g", "k", []byte("v2"))
	assert.Equal(t, []byte("v2"), md.Value("g", "k"))
}

func TestUserMetadata_RemoveKey(t *testing.T) {
	md := NewUserMetadata()
	md.AddKV("g", "k1", []byte("v1"))
	md.AddKV("g", "k2", []byte("v2"))
	md.RemoveKey("g", "k1")
	assert.Empty(t, md.Value("g", "k1"))
	assert.Equal(t, []byte("v2"), md.Value("g", "k2"))

	// Remove from nonexistent group is a no-op
	md.RemoveKey("nonexistent", "k1")
}

func TestUserMetadata_RemoveGroup(t *testing.T) {
	md := NewUserMetadata()
	md.AddKV("g1", "k", []byte("v"))
	md.AddKV("g2", "k", []byte("v"))
	md.RemoveGroup("g1")
	assert.Len(t, md.Groups(), 1)
	assert.Contains(t, md.Groups(), "g2")

	// Remove nonexistent group is a no-op
	md.RemoveGroup("nonexistent")
}
