package metadata

import (
	"strings"

	mappb "github.com/numaproj/numaflow-go/pkg/apis/proto/map/v1"
)

// Metadata wraps the protobuf MessageMetadata
type Metadata struct {
	*mappb.MessageMetadata
}

// New creates a new empty Metadata instance
//
// Users should use the New() helper function only when datum.Metadata() returns nil
// This ensures that any existing metadata from previous vertex is not overriddden
// All metadata keys and groups are stored in lowercase for consistency.
func New() *Metadata {
	return &Metadata{
		MessageMetadata: &mappb.MessageMetadata{
			SysMetadata:  make(map[string]*mappb.KeyValueGroup),
			UserMetadata: make(map[string]*mappb.KeyValueGroup),
		},
	}
}

// Set sets user metadata for a specific group and key
//
// Panics if called on nil Metadata receiver
// For a group,key pair, if the key already exists, it will be overwritten
// A new key value pair is appended to the group if the key does not exist.
// Groups and keys are converted to lowercase for consistency.
func (m *Metadata) Set(group, key string, value []byte) {
	if m == nil || m.MessageMetadata == nil {
		panic("Set called on nil Metadata receiver")
	}

	groupLower := strings.ToLower(group)
	keyLower := strings.ToLower(key)

	// Initialize group if it doesn't exist
	if m.UserMetadata[groupLower] == nil {
		m.UserMetadata[groupLower] = &mappb.KeyValueGroup{
			KeyValue: make(map[string]*mappb.KeyValue),
		}
	}

	// Set the key-value pair
	m.UserMetadata[groupLower].KeyValue[keyLower] = &mappb.KeyValue{
		Key:   keyLower,
		Value: value,
	}
}

// Get retrieves a specific user metadata value by group and key
// Groups and keys are converted to lowercase for consistency.
func (m *Metadata) Get(group, key string) ([]byte, bool) {
	if m.MessageMetadata == nil || m.UserMetadata == nil {
		return nil, false
	}

	// Convert to lowercase for consistency
	groupLower := strings.ToLower(group)
	keyLower := strings.ToLower(key)

	if groupData, exists := m.UserMetadata[groupLower]; exists && groupData != nil {
		if kv, exists := groupData.KeyValue[keyLower]; exists && kv != nil {
			return kv.Value, true
		}
	}
	return nil, false
}

// GetGroup retrieves all key-value pairs for a specific group of UserMetadata
// Groups are converted to lowercase for consistency.
func (m *Metadata) GetGroup(group string) map[string][]byte {
	if m.MessageMetadata == nil || m.UserMetadata == nil {
		return nil
	}

	// Convert to lowercase for consistency
	groupLower := strings.ToLower(group)

	if groupData, exists := m.UserMetadata[groupLower]; exists && groupData != nil {
		result := make(map[string][]byte)
		for key, kv := range groupData.KeyValue {
			if kv != nil {
				result[key] = kv.Value
			}
		}
		return result
	}
	return nil
}

// GetSys retrieves a specific system metadata value by group and key (read-only)
// Groups and keys are converted to lowercase for consistency.
func (m *Metadata) GetSys(group, key string) ([]byte, bool) {
	if m.MessageMetadata == nil || m.SysMetadata == nil {
		return nil, false
	}

	// Convert to lowercase for consistency
	groupLower := strings.ToLower(group)
	keyLower := strings.ToLower(key)

	if groupData, exists := m.SysMetadata[groupLower]; exists && groupData != nil {
		if kv, exists := groupData.KeyValue[keyLower]; exists && kv != nil {
			return kv.Value, true
		}
	}
	return nil, false
}
