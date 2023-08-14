package sinkfn

//go:generate mockgen -destination sinkmock/sinkmock.go -package sinkmock github.com/numaproj/numaflow-go/pkg/apis/proto/sinkfn SinkClient,Sink_SinkFnClient
