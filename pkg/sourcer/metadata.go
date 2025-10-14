package sourcer

// Metadata provides per-message metadata passed between vertices.
// Source is the origin or the first vertex in the pipeline.
// Here, first time the user metadata can be set by the user.
// A vertex could create one or more set of key-value pairs per group-name.

// User metadata format: map[groupName]map[key][]byte,
// where per group name, there can be multiple key-value pairs.

// Metadata wraps user metadata at the source. Only user metadata is set here.
type Metadata struct {
	userMetadata UserMetadata
}

// UserMetadata wraps user-defined metadata groups per message.
type UserMetadata struct {
	data map[string]map[string][]byte
}

// NewUserMetadata creates an empty UserMetadata.
func NewUserMetadata() UserMetadata {
	return UserMetadata{data: make(map[string]map[string][]byte)}
}

// NewUserMetadataWithData wraps an existing map into UserMetadata.
// If d is nil, an empty map is created.
func NewUserMetadataWithData(d map[string]map[string][]byte) UserMetadata {
	if d == nil {
		d = make(map[string]map[string][]byte)
	}
	return UserMetadata{data: d}
}

// Data returns the underlying groups map for iteration.
func (md UserMetadata) Data() map[string]map[string][]byte {
	return md.data
}

// NewMetadata creates a new source Metadata with provided user metadata.
// If nil is passed, an empty user metadata map is created.
// Example:
//
//	user := make(UserMetadata)
//	user["kafka-headers"] = map[string][]byte{"k1": []byte("v1")}
//	sourcer.NewMessage(value, offset, eventTime).WithMetadata(sourcer.NewMetadata(user))
func NewMetadata(userMetadata UserMetadata) Metadata {
	if userMetadata.data == nil {
		userMetadata = NewUserMetadata()
	}
	return Metadata{
		userMetadata: userMetadata,
	}
}

// UserMetadata returns the user metadata map.
// Each group is a key->[]byte map.
// Example: md.UserMetadata()
func (md Metadata) UserMetadata() UserMetadata {
	return md.userMetadata
}
