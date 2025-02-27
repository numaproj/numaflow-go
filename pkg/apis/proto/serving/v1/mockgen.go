package v1

//go:generate mockgen -destination servingmock/serving.go -package servingmock github.com/numaproj/numaflow-go/pkg/apis/proto/serving/v1 ServingStoreClient
