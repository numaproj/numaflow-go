package function

const (
	DatumKey = "x-numaflow-datum-key"
	// DefaultMaxMessageSize overrides gRPC max message size configuration
	// https://github.com/grpc/grpc-go/blob/master/server.go#L58-L59
	//   - defaultServerMaxReceiveMessageSize
	//   - defaultServerMaxSendMessageSize
	DefaultMaxMessageSize = 1024 * 1024 * 4
	WinStartTime          = "x-numaflow-win-start-time"
	WinEndTime            = "x-numaflow-win-end-time"
)

var Protocol string
var Addr string
var MAP_MULTIPROC_SERV bool

//
//func init() {
//	if os.Getenv("MAP_MULTIPROC") == "true" {
//		Protocol = "tcp"
//		Addr = ":55551"
//		MAP_MULTIPROC_SERV = true
//	} else {
//		Protocol = "unix"
//		Addr = "/var/run/numaflow/function.sock"
//		MAP_MULTIPROC_SERV = false
//	}
//}
