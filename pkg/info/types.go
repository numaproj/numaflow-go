package info

type Protocol string

const (
	UDS Protocol = "uds"
	TCP Protocol = "tcp"
)

type Language string

const (
	Go Language = "go"
)

type ContainerType string

// the string content matches the corresponding server info file name.
// DO NOT change it unless the server info file name is changed.
const (
	Sourcer           ContainerType = "sourcer"
	Sourcetransformer ContainerType = "sourcetransformer"
	Sinker            ContainerType = "sinker"
	Mapper            ContainerType = "mapper"
	Reducer           ContainerType = "reducer"
	Reducestreamer    ContainerType = "reducestreamer"
	Sessionreducer    ContainerType = "sessionreducer"
	Sideinput         ContainerType = "sideinput"
	Fbsinker          ContainerType = "fb-sinker"
	Serving           ContainerType = "serving"
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
var MinimumNumaflowVersion = map[ContainerType]string{
	Sourcer:           "1.4.0-z",
	Sourcetransformer: "1.4.0-z",
	Sinker:            "1.4.0-z",
	Mapper:            "1.4.0-z",
	Reducestreamer:    "1.4.0-z",
	Reducer:           "1.4.0-z",
	Sessionreducer:    "1.4.0-z",
	Sideinput:         "1.4.0-z",
	Fbsinker:          "1.4.0-z",
	Serving:           "1.5.0-z",
}

// ServerInfo is the information about the server
type ServerInfo struct {
	Protocol               Protocol          `json:"protocol"`
	Language               Language          `json:"language"`
	MinimumNumaflowVersion string            `json:"minimum_numaflow_version"`
	Version                string            `json:"version"`
	Metadata               map[string]string `json:"metadata"`
}
