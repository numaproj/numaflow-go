package sink

const (
	Protocol = "unix"
	Addr     = "/var/run/numaflow/udsink.sock"
	// DefaultMaxMessageSize overrides gRPC max message size configuration
	// https://github.com/grpc/grpc-go/blob/master/server.go#L58-L59
	//   - defaultServerMaxReceiveMessageSize
	//	 - defaultServerMaxSendMessageSize
	DefaultMaxMessageSize = 1024 * 1024 * 4
)
