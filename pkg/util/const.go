package util

const (
	TCP          = "tcp"
	UDS          = "unix"
	FunctionAddr = "/var/run/numaflow/function.sock"
	TcpAddr      = ":55551"
	// DefaultMaxMessageSize overrides gRPC max message size configuration
	// https://github.com/grpc/grpc-go/blob/master/server.go#L58-L59
	//   - defaultServerMaxReceiveMessageSize
	//   - defaultServerMaxSendMessageSize
	DefaultMaxMessageSize = 1024 * 1024 * 64
	WinStartTime          = "x-numaflow-win-start-time"
	WinEndTime            = "x-numaflow-win-end-time"
	Delimiter             = ":"
	SinkAddr              = "/var/run/numaflow/udsink.sock"
)
