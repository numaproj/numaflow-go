package sinker

// Since sink is the last vertex in the pipeline, user metadata
// should be treated as read-only. Write methods are available
// on the type but should not be used in sink implementations.

import "github.com/numaproj/numaflow-go/internal/metadata"

// SystemMetadata wraps system-generated metadata groups per message.
// It is read-only to UDFs.
type SystemMetadata = metadata.SystemMetadata

// UserMetadata wraps user-defined metadata groups per message.
type UserMetadata = metadata.UserMetadata

// NewSystemMetadata creates a new SystemMetadata.
// This is for internal and testing purposes only.
var NewSystemMetadata = metadata.NewSystemMetadata

// NewUserMetadata creates a new UserMetadata.
var NewUserMetadata = metadata.NewUserMetadata
