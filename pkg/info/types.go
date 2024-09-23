package info

type Protocol string

const (
	UDS Protocol = "uds"
	TCP Protocol = "tcp"
)

type Language string

const (
	Go     Language = "go"
	Python Language = "python"
	Java   Language = "java"
	Rust   Language = "rust"
)

type MapMode string

const (
	UnaryMap  MapMode = "unary-map"
	StreamMap MapMode = "stream-map"
	BatchMap  MapMode = "batch-map"
)

// MapModeKey is the key used in the server info metadata map to indicate which map mode is enabled.
const MapModeKey = "MAP_MODE"

// MinimumNumaflowVersion is the minimum version of Numaflow required by the current SDK version
// To update this value, please follow the instructions for MINIMUM_NUMAFLOW_VERSION in
// https://github.com/numaproj/numaflow-rs/blob/main/src/shared.rs
const MinimumNumaflowVersion = "1.3.1-z"

// ServerInfo is the information about the server
type ServerInfo struct {
	Protocol               Protocol          `json:"protocol"`
	Language               Language          `json:"language"`
	MinimumNumaflowVersion string            `json:"minimum_numaflow_version"`
	Version                string            `json:"version"`
	Metadata               map[string]string `json:"metadata"`
}
