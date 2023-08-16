package v1

//go:generate mockgen -destination transformermock/transformermock.go -package transformermock github.com/numaproj/numaflow-go/pkg/apis/proto/sourcetransformer/v1 SourceTransformerClient
