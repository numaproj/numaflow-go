package info

const (
	SocketAddress = "/var/run/numaflow/info.sock"
)

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

// ServerInfo is the information about the server
type ServerInfo struct {
	Protocol  Protocol          `json:"protocol"`
	Language  Language          `json:"language"`
	Version   string            `json:"version"`
	Metaddata map[string]string `json:"metadata"`
}