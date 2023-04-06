package function

import "os"

const (
	TCP      = "tcp"
	UDS      = "unix"
	UDS_ADDR = "/var/run/numaflow/function.sock"
	TCP_ADDR = ":55551"
	DatumKey = "x-numaflow-datum-key"
	// DefaultMaxMessageSize overrides gRPC max message size configuration
	// https://github.com/grpc/grpc-go/blob/master/server.go#L58-L59
	//   - defaultServerMaxReceiveMessageSize
	//   - defaultServerMaxSendMessageSize
	DefaultMaxMessageSize = 1024 * 1024 * 4
	WinStartTime          = "x-numaflow-win-start-time"
	WinEndTime            = "x-numaflow-win-end-time"
)

func IsMapMultiProcEnabled() bool {
	val, present := os.LookupEnv("MAP_MULTIPROC")
	if present && val == "true" {
		return true
	}
	return false
}
