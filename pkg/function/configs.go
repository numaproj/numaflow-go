package function

const (
	Protocol = "unix"
	Addr     = "/var/run/numaflow/function.sock"
	// DefaultMaxMessageSize overrides gRPC max message size configuration
	// https://github.com/grpc/grpc-go/blob/master/server.go#L58-L59
	//   - defaultServerMaxReceiveMessageSize
	//   - defaultServerMaxSendMessageSize
	DefaultMaxMessageSize = 1024 * 1024 * 4
	WinStartTime          = "x-numaflow-win-start-time"
	WinEndTime            = "x-numaflow-win-end-time"
)
