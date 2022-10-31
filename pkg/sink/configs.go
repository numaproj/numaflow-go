package sink

const (
	Protocol              = "unix"
	Addr                  = "/var/run/numaflow/udsink.sock"
	DefaultMaxMessageSize = 1024 * 1024 * 4
)
