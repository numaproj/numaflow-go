package sourcer

type SystemMetadata map[string]map[string][]byte
type UserMetadata map[string]map[string][]byte

type Metadata struct {
	systemMetadata SystemMetadata
	userMetadata   UserMetadata
}

// NewMetadata creates a new Metadata with user metadata initialized for origin (source).
// Only user metadata is meant to be set by user code; system metadata will be populated by the platform.
func NewMetadata(userMetadata UserMetadata) Metadata {
	if userMetadata == nil {
		userMetadata = make(UserMetadata)
	}
	return Metadata{
		systemMetadata: make(SystemMetadata),
		userMetadata:   userMetadata,
	}
}

func (md Metadata) SystemMetadata() SystemMetadata {
	return md.systemMetadata
}

func (md Metadata) UserMetadata() UserMetadata {
	return md.userMetadata
}
