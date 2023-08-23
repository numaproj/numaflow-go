package shared

const (
	TCP                   = "tcp"
	UDS                   = "unix"
	MapAddr               = "/var/run/numaflow/map.sock"
	MapStreamAddr         = "/var/run/numaflow/mapstream.sock"
	ReduceAddr            = "/var/run/numaflow/reduce.sock"
	SourceAddr            = "/var/run/numaflow/source.sock"
	SinkAddr              = "/var/run/numaflow/sink.sock"
	SourceTransformerAddr = "/var/run/numaflow/sourcetransform.sock"
	SideInputAddr         = "/var/run/numaflow/sideinput.sock"
	TcpAddr               = ":55551"
	// DefaultMaxMessageSize overrides gRPC max message size configuration
	// https://github.com/grpc/grpc-go/blob/master/server.go#L58-L59
	//   - defaultServerMaxReceiveMessageSize
	//   - defaultServerMaxSendMessageSize
	DefaultMaxMessageSize = 1024 * 1024 * 64
	WinStartTime          = "x-numaflow-win-start-time"
	WinEndTime            = "x-numaflow-win-end-time"
	Delimiter             = ":"
)
