package v1

//go:generate mockgen -destination sessionreducemock/sessionreducemock -package sessionreducemock github.com/numaproj/numaflow-go/pkg/apis/proto/sessionreduce/v1 SessionReduceClient,SessionReduce_SessionReduceFnClient
