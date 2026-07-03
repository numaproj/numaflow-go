package sinker

import "github.com/numaproj/numaflow-go/internal/nackoptions"

// NackOptions carries per-message redelivery options for a nacked message.
// See ResponseNack.
type NackOptions = nackoptions.NackOptions
