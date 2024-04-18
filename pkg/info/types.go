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
)

// MinimumNumaflowVersion specifies the minimum Numaflow version required by the current SDK version
const MinimumNumaflowVersion = "1.2.0-rc4"

// ServerInfo is the information about the server
type ServerInfo struct {
	Protocol               Protocol          `json:"protocol"`
	Language               Language          `json:"language"`
	MinimumNumaflowVersion string            `json:"minimum_numaflow_version"`
	Version                string            `json:"version"`
	Metadata               map[string]string `json:"metadata"`
}
