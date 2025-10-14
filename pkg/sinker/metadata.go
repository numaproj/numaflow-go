package sinker

// Since sink is the last vertex in the pipeline, only GET methods
// are available on SystemMetadata and UserMetadata.

// Metadata contains system and user metadata available at the sink.
type Metadata struct {
	systemMetadata SystemMetadata
	userMetadata   UserMetadata
}

// SystemMetadata and UserMetadata are mappings of group name to key-value pairs.
type SystemMetadata struct {
	data map[string]map[string][]byte
}
type UserMetadata struct {
	data map[string]map[string][]byte
}

// newSystemMetadata wraps an existing map into SystemMetadata
// This is for internal use only.
func newSystemMetadata(d map[string]map[string][]byte) SystemMetadata {
	if d == nil {
		d = make(map[string]map[string][]byte)
	}
	return SystemMetadata{data: d}
}

// newUserMetadata wraps an existing map into UserMetadata
// This is for internal use only.
func newUserMetadata(d map[string]map[string][]byte) UserMetadata {
	if d == nil {
		d = make(map[string]map[string][]byte)
	}
	return UserMetadata{data: d}
}

// UserMetadata returns the user metadata map.
// Any changes to the user metadata will be ignored.
func (md Metadata) UserMetadata() UserMetadata {
	return md.userMetadata
}

// SystemMetadata returns the system metadata map.
// Any changes to the system metadata will be ignored.
func (md Metadata) SystemMetadata() SystemMetadata {
	return md.systemMetadata
}
