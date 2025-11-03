package sourcetransformer

// SystemMetadata is mapping of group name to key-value pairs
// SystemMetadata wraps system-generated metadata groups per message.
// It is read-only to UDFs
type SystemMetadata struct {
	data map[string]map[string][]byte
}

// NewSystemMetadata creates a new SystemMetadata.
// This is for internal and testing purposes only.
func NewSystemMetadata() *SystemMetadata {
	d := make(map[string]map[string][]byte)
	return &SystemMetadata{data: d}
}

// Groups returns the groups of the system metadata.
// If there are no groups, it returns an empty slice.
//
// Usage example:
//
//	```go
//	systemMetadata := datum.SystemMetadata()
//	groups := systemMetadata.Groups()
//	```
func (md *SystemMetadata) Groups() []string {
	if md == nil || md.data == nil {
		return []string{}
	}
	groups := make([]string, 0, len(md.data))
	for group := range md.data {
		groups = append(groups, group)
	}
	return groups
}

// Keys returns the keys of the system metadata for the given group.
// If the group is not present, it returns an empty slice.
//
// Usage example:
//
//	```go
//	systemMetadata := datum.SystemMetadata()
//	keys := systemMetadata.Keys("group-name")
//	```
func (md *SystemMetadata) Keys(group string) []string {
	if md == nil || md.data == nil {
		return []string{}
	}
	keys := make([]string, 0, len(md.data[group]))
	for key := range md.data[group] {
		keys = append(keys, key)
	}
	return keys
}

// Value returns the value of the system metadata for the given group and key.
// If the group or key is not present, it returns an empty slice.
//
// Usage example:
//
//	```go
//	systemMetadata := datum.SystemMetadata()
//	value := systemMetadata.Value("group-name", "key")
//	```
func (md *SystemMetadata) Value(group, key string) []byte {
	if md == nil || md.data == nil {
		return []byte{}
	}
	return md.data[group][key]
}

// UserMetadata wraps user-defined metadata groups per message.
type UserMetadata struct {
	data map[string]map[string][]byte
}

// NewUserMetadata creates a new UserMetadata.
func NewUserMetadata() *UserMetadata {
	d := make(map[string]map[string][]byte)
	return &UserMetadata{data: d}
}

// Groups returns the groups of the user metadata.
// If there are no groups, it returns an empty slice.
//
// Usage example:
//
//	```go
//	userMetadata := datum.UserMetadata()
//	groups := userMetadata.Groups()
//	```
func (md *UserMetadata) Groups() []string {
	if md == nil || md.data == nil {
		return []string{}
	}
	groups := make([]string, 0, len(md.data))
	for group := range md.data {
		groups = append(groups, group)
	}
	return groups
}

// Keys returns the keys of the user metadata for the given group.
// If the group is not present, it returns an empty slice.
//
// Usage example:
//
//	```go
//	userMetadata := datum.UserMetadata()
//	keys := userMetadata.Keys("group-name")
//	```
func (md *UserMetadata) Keys(group string) []string {
	if md == nil || md.data == nil {
		return []string{}
	}
	keys := make([]string, 0, len(md.data[group]))
	for key := range md.data[group] {
		keys = append(keys, key)
	}
	return keys
}

// Value returns the value of the user metadata for the given group and key.
// If the group or key is not present, it returns an empty slice.
//
// Usage example:
//
//	```go
//	userMetadata := datum.UserMetadata()
//	value := userMetadata.Value("group-name", "key")
//	```
func (md *UserMetadata) Value(group, key string) []byte {
	if md == nil || md.data == nil {
		return []byte{}
	}
	return md.data[group][key]
}

// CreateGroup creates a new group in the user metadata.
//
// Usage example:
//
//	```go
//	userMetadata := NewUserMetadata()
//	userMetadata.CreateGroup("group-name")
//	```
func (md *UserMetadata) CreateGroup(group string) {
	if md.data == nil {
		md.data = make(map[string]map[string][]byte)
	}
	if md.data[group] == nil {
		md.data[group] = make(map[string][]byte)
	}
}

// AddKV adds a key-value pair under the given group name to the user metadata.
// If the group is not present, it creates a new group.
// If the key is already present, it overwrites the value.
//
// Usage example:
//
//	```go
//	userMetadata := NewUserMetadata()
//	userMetadata.AddKV("group-name", "key", []byte("value"))
//	```
func (md *UserMetadata) AddKV(group, key string, value []byte) {
	if md.data == nil {
		md.data = make(map[string]map[string][]byte)
	}
	if md.data[group] == nil {
		md.data[group] = make(map[string][]byte)
	}
	md.data[group][key] = value
}

// RemoveKey removes a key from a group in the user metadata.
// If the key or group is not present, it's a no-op.
//
// Usage example:
//
//	```go
//	md := datum.Metadata()
//	userMetadata := md.UserMetadata()
//	userMetadata.RemoveKey("group-name", "key")
//	```
func (md *UserMetadata) RemoveKey(group, key string) {
	if md == nil || md.data == nil {
		return
	}
	delete(md.data[group], key)
}

// RemoveGroup removes a group from the user metadata.
// If the group is not present, it's a no-op.
//
// Usage example:
//
//	```go
//	md := datum.Metadata()
//	userMetadata := md.UserMetadata()
//	userMetadata.RemoveGroup("group-name")
//	```
func (md *UserMetadata) RemoveGroup(group string) {
	if md == nil || md.data == nil {
		return
	}
	delete(md.data, group)
}
