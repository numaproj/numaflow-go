package sinker

// Since sink is the last vertex in the pipeline, only GET methods
// are available on SystemMetadata and UserMetadata.

// SystemMetadata and UserMetadata are mappings of group name to key-value pairs.
type SystemMetadata map[string]map[string][]byte
type UserMetadata map[string]map[string][]byte

// Metadata contains system and user metadata available at the sink.
type Metadata struct {
	systemMetadata SystemMetadata
	userMetadata   UserMetadata
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
