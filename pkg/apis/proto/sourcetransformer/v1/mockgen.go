package v1

//go:generate mockgen -destination funcmock/funcmock.go -package funcmock github.com/numaproj/numaflow-go/pkg/apis/proto/sourcetransformer/v1 SourceTransformerClient
