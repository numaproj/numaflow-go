package function

const (
	Protocol              = "unix"
	Addr                  = "/var/run/numaflow/function.sock"
	DatumKey              = "x-numaflow-datum-key"
	DefaultMaxMessageSize = 1024 * 1024 * 4
)
