package sinker

// Since sink is the last vertex in the pipeline, only GET methods
// are available on SystemMetadata and UserMetadata.

// SystemMetadata wraps system-generated metadata groups per message.
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
//	systemMetadata := datum.SystemMetadata()
//	groups := systemMetadata.Groups()
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
//	systemMetadata := datum.SystemMetadata()
//	keys := systemMetadata.Keys("group-name")
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
//	systemMetadata := datum.SystemMetadata()
//	value := systemMetadata.Value("group-name", "key")
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
//	userMetadata := datum.UserMetadata()
//	groups := userMetadata.Groups()
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
//	userMetadata := datum.UserMetadata()
//	keys := userMetadata.Keys("group-name")
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
//	userMetadata := datum.UserMetadata()
//	value := userMetadata.Value("group-name", "key")
func (md *UserMetadata) Value(group, key string) []byte {
	if md == nil || md.data == nil {
		return []byte{}
	}
	return md.data[group][key]
}
