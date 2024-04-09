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

const MinimumClientVersion = ">= 1.2.0-0"

// ServerInfo is the information about the server
type ServerInfo struct {
	Protocol             Protocol          `json:"protocol"`
	Language             Language          `json:"language"`
	MinimumClientVersion string            `json:"minimumClientVersion"`
	Version              string            `json:"version"`
	Metadata             map[string]string `json:"metadata"`
}
