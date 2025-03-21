package v1

//go:generate mockgen -destination accumulatoremock/accumulatormock.go -package accumulatormock github.com/numaproj/numaflow-go/pkg/apis/proto/accumulator/v1 AccumulatorClient,Accumulator_AccumulateFnClient
