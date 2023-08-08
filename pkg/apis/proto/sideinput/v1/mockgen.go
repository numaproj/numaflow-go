package v1

//go:generate mockgen -destination sideinputmock/sideinputmock.go -package sideinputmock github.com/KeranYang/numaflow-go/pkg/apis/proto/sideinput/v1 UserDefinedSideInputClient
