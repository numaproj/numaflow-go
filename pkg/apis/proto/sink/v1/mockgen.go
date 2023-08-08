package v1

//go:generate mockgen -destination sinkmock/sinkmock.go -package sinkmock github.com/KeranYang/numaflow-go/pkg/apis/proto/sink/v1 UserDefinedSinkClient,UserDefinedSink_SinkFnClient
