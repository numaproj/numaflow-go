package sinker

// SystemMetadata and UserMetadata are mappings of group name to key-value pairs
type SystemMetadata map[string]map[string][]byte
type UserMetadata map[string]map[string][]byte

type Metadata struct {
	previousVertex string
	systemMetadata SystemMetadata
	userMetadata   UserMetadata
}

func (md Metadata) PreviousVertex() string {
	return md.previousVertex
}

func (md Metadata) SystemMetadata() SystemMetadata {
	return md.systemMetadata
}

func (md Metadata) UserMetadata() UserMetadata {
	return md.userMetadata
}
