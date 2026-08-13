package batchmapper

import "github.com/numaproj/numaflow-go/internal/nackoptions"

// NackOptions carries per-message redelivery options for a nacked message.
// See MessageToNack.
type NackOptions = nackoptions.NackOptions
