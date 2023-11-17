package v1

//go:generate mockgen -destination globalreducemock/globalreducemock.go -package globalreducemock github.com/numaproj/numaflow-go/pkg/apis/proto/globalreduce/v1 GlobalReduceClient,GlobalReduce_GlobalReduceFnClient
