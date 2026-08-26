// Package numaflow defines the root Server interface shared by Numaflow SDK components.
package numaflow

import "context"

// Server is the interface for all numaflow servers.
type Server interface {
	// Start starts the server.
	Start(ctx context.Context) error
}
