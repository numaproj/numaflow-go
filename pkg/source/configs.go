package source

const (
	Protocol = "unix"
	Addr     = "/var/run/numaflow/source.sock"
	DatumKey = "x-numaflow-datum-key"
	// DefaultMaxMessageSize overrides gRPC max message size configuration
	// https://github.com/grpc/grpc-go/blob/master/server.go#L58-L59
	//   - defaultServerMaxReceiveMessageSize
	//   - defaultServerMaxSendMessageSize
	DefaultMaxMessageSize = 1024 * 1024 * 4
)
