package mapper

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
