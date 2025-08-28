package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHelperFunction(t *testing.T) {
	// Test creating new metadata
	md := New()
	if md == nil {
		t.Fatal("Expected non-nil metadata")
	}

	// Test that it's properly initialized
	if md.MessageMetadata == nil {
		t.Fatal("Expected MessageMetadata to be initialized")
	}

	if md.UserMetadata == nil {
		t.Fatal("Expected UserMetadata to be initialized")
	}

	if md.SysMetadata == nil {
		t.Fatal("Expected SysMetadata to be initialized")
	}
}

func TestSetOnNil(t *testing.T) {
	// Test setting metadata on nil receiver (should panic - this is expected Go behavior)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when calling Set on nil receiver")
		}
	}()

	var md *Metadata // nil
	md.Set("test", "key", []byte("value"))
}

func TestSetOnValid(t *testing.T) {
	// Test setting metadata on valid receiver
	md := New()
	md.Set("app", "version", []byte("v1.0"))

	// Verify it was set
	value, exists := md.Get("app", "version")
	if !exists {
		t.Error("Expected metadata to exist after setting")
	}

	if string(value) != "v1.0" {
		t.Errorf("Expected value 'v1.0', got '%s'", string(value))
	}
}

func TestSetCaseInsensitive(t *testing.T) {
	meta := New()

	// Set metadata with mixed case
	meta.Set("App", "Version", []byte("v1.0"))
	meta.Set("APP", "VERSION", []byte("v2.0"))
	meta.Set("app", "version", []byte("v3.0"))

	// All should retrieve the same value (last one set)
	value, exists := meta.Get("App", "Version")
	assert.True(t, exists)
	assert.Equal(t, []byte("v3.0"), value)

	value, exists = meta.Get("APP", "VERSION")
	assert.True(t, exists)
	assert.Equal(t, []byte("v3.0"), value)

	value, exists = meta.Get("app", "version")
	assert.True(t, exists)
	assert.Equal(t, []byte("v3.0"), value)

	// Check that only one group and key exist (lowercase)
	assert.Len(t, meta.UserMetadata, 1)
	assert.Contains(t, meta.UserMetadata, "app")
	assert.Len(t, meta.UserMetadata["app"].KeyValue, 1)
	assert.Contains(t, meta.UserMetadata["app"].KeyValue, "version")
}

func TestGetGroup(t *testing.T) {
	meta := New()

	// Set some metadata
	meta.Set("App", "Version", []byte("v1.0"))
	meta.Set("App", "Build", []byte("123"))

	// Get all metadata for the group
	groupData := meta.GetGroup("App")
	assert.NotNil(t, groupData)
	assert.Len(t, groupData, 2)

	// Check specific values
	assert.Equal(t, []byte("v1.0"), groupData["version"])
	assert.Equal(t, []byte("123"), groupData["build"])

	// Test with non-existent group
	emptyGroup := meta.GetGroup("NonExistent")
	assert.Nil(t, emptyGroup)
}
