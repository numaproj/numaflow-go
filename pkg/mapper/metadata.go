package mapper

import "strconv"

// SystemMetadata and UserMetadata are mappings of group name to key-value pairs
type SystemMetadata map[string]map[string][]byte
type UserMetadata map[string]map[string][]byte

type Metadata struct {
	previousVertex string
	systemMetadata SystemMetadata
	userMetadata   UserMetadata
}

// SystemMetadataView is a read-only interface for inspecting system metadata.
type SystemMetadataView interface {
	// Get returns the value of the given group and key of the system metadata.
	Get(group, key string) []byte
	// Keys returns the keys of the given group of the system metadata.
	Keys(group string) []string
	// Groups returns the groups of the system metadata.
	Groups() []string
}

// systemMetadataView wraps the underlying system metadata.
type systemMetadataView struct {
	data SystemMetadata
}

// SystemMetadataView returns a read-only view over the system metadata.
func (md Metadata) SystemMetadataView() SystemMetadataView {
	return systemMetadataView{data: md.systemMetadata}
}

// Get returns the value of the given group and key of the system metadata.
// Example: view := md.SystemMetadataView()
// view.Get(group, key)
func (v systemMetadataView) Get(group, key string) []byte {
	return v.data[group][key]
}

// Keys returns the keys of the given group of the system metadata.
// Example: view := md.SystemMetadataView()
// view.Keys(group)
func (v systemMetadataView) Keys(group string) []string {
	if v.data[group] == nil {
		return nil
	}
	keys := make([]string, 0, len(v.data[group]))
	for key := range v.data[group] {
		keys = append(keys, key)
	}
	return keys
}

// Groups returns the groups of the system metadata.
// Example: view := md.SystemMetadataView()
// view.Groups()
func (v systemMetadataView) Groups() []string {
	groups := make([]string, 0, len(v.data))
	for group := range v.data {
		groups = append(groups, group)
	}
	return groups
}

// PreviousVertex returns the previous vertex.
// Example: md.PreviousVertex()
func (md Metadata) PreviousVertex() string {
	return md.previousVertex
}

// UserMetadata returns the user metadata.
// Example: md.UserMetadata()
func (md Metadata) UserMetadata() UserMetadata {
	return md.userMetadata
}

// SetKVGroup sets a group of key-value pairs in the user metadata.
// Example: md := datum.Metadata()
// md.SetKVGroup(group, kv)
func (md *Metadata) SetKVGroup(group string, kv map[string][]byte) {
	if md.userMetadata == nil {
		md.userMetadata = make(UserMetadata)
	}
	if md.userMetadata[group] == nil {
		md.userMetadata[group] = make(map[string][]byte)
	}
	md.userMetadata[group] = kv
}

// AppendKV appends a key-value pair to the user metadata.
// Example: md := datum.Metadata()
// md.AppendKV("group1", "key1", []byte("value1"))
func (md *Metadata) AppendKV(group, key string, value []byte) {
	if md.userMetadata == nil {
		md.userMetadata = make(UserMetadata)
	}
	if md.userMetadata[group] == nil {
		md.userMetadata[group] = make(map[string][]byte)
	}
	md.userMetadata[group][key] = value
}

// AppendKVString appends a key-value pair of string type to the user metadata.
// Example: md := datum.Metadata()
// md.AppendKVString("group1", "key1", "value1")
func (md *Metadata) AppendKVString(group, key, value string) {
	if md.userMetadata == nil {
		md.userMetadata = make(UserMetadata)
	}
	if md.userMetadata[group] == nil {
		md.userMetadata[group] = make(map[string][]byte)
	}
	md.userMetadata[group][key] = []byte(value)
}

// AppendKVInt appends a key-value pair of int type to the user metadata.
// Example: md := datum.Metadata()
// md.AppendKVInt("group1", "key1", 2)
func (md *Metadata) AppendKVInt(group, key string, value int) {
	if md.userMetadata == nil {
		md.userMetadata = make(UserMetadata)
	}
	if md.userMetadata[group] == nil {
		md.userMetadata[group] = make(map[string][]byte)
	}
	md.userMetadata[group][key] = []byte(strconv.Itoa(value))
}

// RemoveKey removes a key from a group in the user metadata.
// Example: md := datum.Metadata()
// md.RemoveKey("group1", "key1")
func (md *Metadata) RemoveKey(group, key string) {
	// md should never be nil
	if md == nil {
		return
	}
	delete(md.userMetadata[group], key)
}

// RemoveGroup removes a group from the user metadata.
// Example: md := datum.Metadata()
// md.RemoveGroup("group1")
func (md *Metadata) RemoveGroup(group string) {
	// md should never be nil
	if md == nil {
		return
	}
	delete(md.userMetadata, group)
}
