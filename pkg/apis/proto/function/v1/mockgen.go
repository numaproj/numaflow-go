package v1

//go:generate mockgen -destination funcmock/funcmock.go -package funcmock github.com/KeranYang/numaflow-go/pkg/apis/proto/function/v1 UserDefinedFunctionClient,UserDefinedFunction_ReduceFnClient,UserDefinedFunction_MapStreamFnClient
