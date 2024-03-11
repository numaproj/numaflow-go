# Development

This document explains the development process for the Numaflow Go SDK.


### Testing

The SDK uses local references in the `go.mod` file, i.e. the `github.com/numaproj/numaflow-go` dependency is pointing to your local
`numaflow-go` directory. Thus, you can automatically test your SDK changes without the need to push to your forked repository or modify the `go.mod` file.
Simply make your changes, build and push the desired example image, and you are ready to use it in any pipeline.

### Deploying

1. Create a PR for your changes. Once merged, it will trigger a workflow to build and push the images for all the examples, 
with the tag `stable`. This consistent tag name is used so that the tags in the [E2E test pipelines](https://github.com/numaproj/numaflow/tree/main/test) do not need to be 
updated each time a new change is made. 
2. After the changes have been merged it is encouraged to update the dependency management files so that the version 
displayed in the `go.mod` file reflects the commit SHA of the merged changes. This can be done simply by getting the
commit SHA of the merged changes and using it with the update script: 
```shell
./update_examples -c <commit-sha>
```
After running the script, create another PR for these changes.

Updating the version may not seem necessary since we are using local references. However, the client prints
out information related to the server, which includes the SDK version, which is retrieved from the `go.mod` file.
After a change is merged, even though the images will be using the most up to date SDK
version, if the dependency management files are not updated, the logs will print the previous commit SHA as the SDK version.
Thus, in order for the correctness of the server information, consistency, and to avoid future confusion, it is recommended 
to update the `numaflow-go` dependency version across all the example directories, after a change has been made to the SDK.



