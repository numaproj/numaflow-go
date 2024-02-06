package info

const (
	ServerInfoFilePath                  = "/var/run/numaflow/server-info"
	MapperServerInfoFilePath            = "/var/run/numaflow/mapper-server-info"
	MapstreamerServerInfoFilePath       = "/var/run/numaflow/mapstreamer-sever-info"
	ReducerServerInfoFilePath           = "/var/run/numaflow/reducer-server-info"
	ReducestreamerServerInfoFilePath    = "/var/run/numaflow/reducestreamer-server-info"
	SessionreducerServerInfoFilePath    = "/var/run/numaflow/sessionreducer-server-ifno"
	SideinputServerInfoFilePath         = "/var/run/numaflow/sideinput-server-info"
	SinkerServerInfoFilePath            = "/var/run/numaflow/sinker-server-info"
	SourcerServerInfoFilePath           = "/var/run/numaflow/sourcer-server-info"
	SourcetransformerServerInfoFilePath = "/var/run/numaflow/sourcetransformer-server-info"
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
	Protocol Protocol          `json:"protocol"`
	Language Language          `json:"language"`
	Version  string            `json:"version"`
	Metadata map[string]string `json:"metadata"`
}
