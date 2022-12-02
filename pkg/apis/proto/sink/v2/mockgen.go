package v2

//go:generate mockgen -destination sinkmock/sinkmock.go -package sinkmock github.com/numaproj/numaflow-go/pkg/apis/proto/sink/v2 UserDefinedSinkClient,UserDefinedSink_SinkFnClient
