# Development

This document explains the development process for the Numaflow Go SDK.


### Testing

The SDK uses local references in the `go.mod` file, i.e. the `github.com/numaproj/numaflow-go` dependency is pointing to your local
`numaflow-go` directory. Thus, you can automatically test your SDK changes without the need to push to your forked repository or modify the `go.mod` file.
Simply make your changes, build and push the desired example image (with an appropriate tag, such as `test`), and you are ready to use it in any pipeline.

Each example directory has a Makefile which can be used to build, tag, and push images. In most cases, the `image-push` target should be used.
However, `buildx`, which is used to support multiple architectures, does not make a built image available locally. In the case that a developer 
wants this functionality, they can use the `image` target.

After making changes to the SDK, if you want to build all the example images at once, in the root directory you can run:
```shell
./update_examples.sh -bp -t <tag>
```
The default tag is `stable`, but it is recommended you specify your own for testing purposes, as the Github Actions CI uses the `stable` tag. Note: do not forget to clean up testing tags
in the registry, i.e. delete them, once you are done testing.

You can alternatively build a specific example image by running the following in the root directory and providing the path to the Dockerfile (relative to root):
```shell
./update_examples.sh -bpe <path> -t <tag>
```
This is essentially equivalent to running `make image-push TAG=<tag>` in the example directory itself.

### Deploying

1. Create a PR for your changes. Once merged, it will trigger a workflow to build and push the images for all the examples, 
with the tag `stable`. This consistent tag name is used so that the tags in the [E2E test pipelines](https://github.com/numaproj/numaflow/tree/main/test) do not need to be 
updated each time a new change is made. 
2. If the change that you have introduced is a breaking change, i.e. adds new features/logic, then after merge it is encouraged to update the dependency management files so that the version 
displayed in the `go.mod` file reflects the commit SHA of the merged changes. This can be done by getting the
commit SHA of the merged changes and using it with the update script. Alternatively, you can provide the specific version that you would like to update to, or even
pass in `latest` to fetch the latest version from the remote repository:
```shell
./update_examples.sh -u <SDK-version | commit-sha | latest>
```
After running the script, create another PR for these changes. Ideally, the update script should only be need to run when a new version is released, i.e. provide a version or `latest` to it,
or when a breaking change is merged before the next release, i.e. provide a commit SHA to it. If your merged change is a small chore, then it is unnecessary to run the update script as we want to
avoid flooding the commit history with dependency updates.

Updating the version may not seem necessary since we are using local references. However, the client prints
out information related to the server, which includes the SDK version, which is retrieved from the `go.mod` file.
After a change is merged, even though the images will be using the most up to date SDK
version, if the dependency management files are not updated, the logs will print the previous commit SHA as the SDK version.
Thus, in order for the correctness of the server information, consistency, and to avoid future confusion, it is recommended 
to update the `numaflow-go` dependency version across all the example directories, after a  large/breaking change has been made to the SDK.



