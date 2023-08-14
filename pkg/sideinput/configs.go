package sideinput

const (
	Protocol = "unix"
	Addr     = "/var/run/numaflow/udsideinput.sock"
	// DefaultMaxMessageSize overrides gRPC max message size configuration
	// https://github.com/grpc/grpc-go/blob/master/server.go#L58-L59
	//   - defaultServerMaxReceiveMessageSize
	//   - defaultServerMaxSendMessageSize
	DefaultMaxMessageSize = 1024 * 1024 * 64
	DefaultTimeout        = 120
)
