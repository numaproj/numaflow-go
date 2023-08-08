package source

const (
	TCP      = "tcp"
	UDS      = "unix"
	UdsAddr  = "/var/run/numaflow/source.sock"
	TCP_ADDR = ":55551"
	// DefaultMaxMessageSize overrides gRPC max message size configuration
	// https://github.com/grpc/grpc-go/blob/master/server.go#L58-L59
	//   - defaultServerMaxReceiveMessageSize
	//   - defaultServerMaxSendMessageSize
	DefaultMaxMessageSize = 1024 * 1024 * 64
	Delimiter             = ":"
)
