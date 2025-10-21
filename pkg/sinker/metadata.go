package sinker

// Since sink is the last vertex in the pipeline, only GET methods
// are available on SystemMetadata and UserMetadata.

// SystemMetadata wraps system-generated metadata groups per message.
type SystemMetadata struct {
	data map[string]map[string][]byte
}

// NewSystemMetadata wraps an existing map into SystemMetadata
// This is for internal and testing purposes only.
func NewSystemMetadata(d map[string]map[string][]byte) *SystemMetadata {
	if d == nil {
		d = make(map[string]map[string][]byte)
	}
	return &SystemMetadata{data: d}
}

// Groups returns the groups of the system metadata.
// If there are no groups, it returns an empty slice.
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

// NewUserMetadata wraps an existing map into UserMetadata
// This is for internal and testing purposes only.
func NewUserMetadata(d map[string]map[string][]byte) *UserMetadata {
	if d == nil {
		d = make(map[string]map[string][]byte)
	}
	return &UserMetadata{data: d}
}

// Groups returns the groups of the user metadata.
// If there are no groups, it returns an empty slice.
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
func (md *UserMetadata) Value(group, key string) []byte {
	if md == nil || md.data == nil {
		return []byte{}
	}
	return md.data[group][key]
}
