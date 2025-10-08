package sinker

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

// UserMetadataView is a read-only interface for inspecting user metadata.
// User Metadata is not allowed to be updated as it is the last stage of the pipeline.
type UserMetadataView interface {
	// Get returns the value of the given group and key of the user metadata.
	Get(group, key string) []byte
	// Keys returns the keys of the given group of the user metadata.
	Keys(group string) []string
	// Groups returns the groups of the user metadata.
	Groups() []string
}

// userMetadataView wraps the underlying user metadata.
type userMetadataView struct {
	data UserMetadata
}

// UserMetadataView returns a read-only view over the user metadata.
func (md Metadata) UserMetadataView() UserMetadataView {
	return userMetadataView{data: md.userMetadata}
}

// Get returns the value of the given group and key of the user metadata.
// Example: view := md.UserMetadataView()
// view.Get(group, key)
func (v userMetadataView) Get(group, key string) []byte {
	return v.data[group][key]
}

// Keys returns the keys of the given group of the user metadata.
// Example: view := md.UserMetadataView()
// view.Keys(group)
func (v userMetadataView) Keys(group string) []string {
	if v.data[group] == nil {
		return nil
	}
	keys := make([]string, 0, len(v.data[group]))
	for key := range v.data[group] {
		keys = append(keys, key)
	}
	return keys
}

// Groups returns the groups of the user metadata.
// Example: view := md.UserMetadataView()
// view.Groups()
func (v userMetadataView) Groups() []string {
	groups := make([]string, 0, len(v.data))
	for group := range v.data {
		groups = append(groups, group)
	}
	return groups
}

func (md Metadata) PreviousVertex() string {
	return md.previousVertex
}
