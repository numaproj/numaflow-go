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

### Proto File Changes

The source proto files are located at [numaflow](https://github.com/numaproj/numaflow/tree/main/pkg/apis/proto) repository, so any proto file change should start from numaflow repository change, and Go should always be the 1st SDK to update.

This means, usually there will be two PRs to update gRPC contracts, one in `numaflow` for proto file changes, the other one in `numaflow-go` to update the detailed implementation. There is one trick to unblock the code change in `numaflow-go` to run codegen for proto files, before the proto file change PR is merged in `numaflow` repository:

```shell
make proto ORG=xxx BRANCH=xxx
```

This will run proto gen with the source files defined in your own repository and branch.

For example:

```shell
# Try to fetch proto files defined in branch "proto-chg" of https://github.com/anyone/numaflow
make proto ORG=anyone BRANCH=proto-chg

# Try to fetch proto files defined in branch "my-change" of https://github.com/numaproj/numaflow
make proto BRANCH=my-change
```
