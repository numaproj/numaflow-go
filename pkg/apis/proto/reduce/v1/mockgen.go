package v1

//go:generate mockgen -destination reducemock/reducemock.go -package funcmock github.com/numaproj/numaflow-go/pkg/apis/proto/reduce/v1 ReduceClient,Reduce_ReduceFnClient
