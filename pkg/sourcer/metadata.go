package sourcer

type SystemMetadata map[string]map[string][]byte
type UserMetadata map[string]map[string][]byte

type Metadata struct {
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
// Example: view := md.SysMetadataView()
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

// NewMetadata creates a new Metadata with user metadata initialized at origin (source).
// Only user metadata is meant to be set by the user.
// Example: userMetadata := make(UserMetadata{"group1": {"key1": "value1", "key2": "value2"}})
// md := NewMetadata(userMetadata)
// sourcer.NewMessage(value, offset, eventTime).WithMetadata(md)
func NewMetadata(userMetadata UserMetadata) Metadata {
	if userMetadata == nil {
		userMetadata = make(UserMetadata)
	}
	return Metadata{
		systemMetadata: make(SystemMetadata),
		userMetadata:   userMetadata,
	}
}

// UserMetadata returns the user metadata.
// Example: md.UserMetadata()
func (md Metadata) UserMetadata() UserMetadata {
	return md.userMetadata
}
