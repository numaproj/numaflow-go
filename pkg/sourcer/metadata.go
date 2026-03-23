package sourcer

import "github.com/numaproj/numaflow-go/internal/metadata"

// UserMetadata wraps user-defined metadata groups per message.
// Source is the origin or the first vertex in the pipeline.
// Here, first time the user metadata can be set by the user.
type UserMetadata = metadata.UserMetadata

// NewUserMetadata creates a new UserMetadata.
var NewUserMetadata = metadata.NewUserMetadata
