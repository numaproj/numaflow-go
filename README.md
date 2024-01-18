# Numaflow Golang SDK

This SDK provides the interfaces to implement [Numaflow](https://github.com/numaproj/numaflow) User Defined Sources,
Source Transformer, Functions, Sinks or SideInputs in Golang.

- Implement [User Defined Sources](https://pkg.go.dev/github.com/numaproj/numaflow-go/pkg/sourcer)
- Implement [User Defined Source Transformers](https://pkg.go.dev/github.com/numaproj/numaflow-go/pkg/sourcetransformer)
- Implement User Defined Functions
  - [Map](https://pkg.go.dev/github.com/numaproj/numaflow-go/pkg/mapper)
  - [Reduce](https://pkg.go.dev/github.com/numaproj/numaflow-go/pkg/reducer)
- Implement [User Defined Sinks](https://pkg.go.dev/github.com/numaproj/numaflow-go/pkg/sinker)
- Implement [User Defined SideInputs](https://pkg.go.dev/github.com/numaproj/numaflow-go/pkg/sideinput)

## Development

### Useful Commands

`make test` - Run the tests.
`make proto`- Regenerate the protobuf files from the [proto files](https://github.com/numaproj/numaflow/tree/main/pkg/apis/proto) defined in [numaproj/numaflow](https://github.com/numaproj/numaflow) repository.
`make proto ORG=xxx PROJECT=xxx BRANCH=xxx` - Regenerate the protobuf files from specified github repository. Default values: `ORG=numaproj PROJECT=numaflow BRANCH=main`
