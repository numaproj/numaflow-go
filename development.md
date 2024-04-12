# Developer Guide

This document explains the development process for the Numaflow Go SDK.

### Testing

The SDK uses local references in the `go.mod` file, i.e. the `github.com/numaproj/numaflow-go` dependency is pointing to your local
`numaflow-go` directory. Thus, you can automatically test your SDK changes without the need to push to your forked repository or modify the `go.mod` file.
Simply make your changes, build the desired example image, and you are ready to use it in any pipeline.

Each example directory has a Makefile which can be used to build, tag, and push images.
If you want to build and tag the image and then immediately push it to quay.io, use the `image-push` target.
If you want to build and tag a local image without pushing, use the `image` target.

If you want to build and push all the example images at once, you can run:
```shell
./hack/update_examples.sh -bp -t <tag>
```
The default tag is `stable`, but it is recommended you specify your own for testing purposes, as the Github Actions CI uses the `stable` tag.
This consistent tag name is used so that the tags in the [E2E test pipelines](https://github.com/numaproj/numaflow/tree/main/test) do not need to be updated each time an SDK change is made.

You can alternatively build and push a specific example image by running the following:
```shell
./hack/update_examples.sh -bpe <path-to-dockerfile> -t <tag>
```
This is essentially equivalent to running `make image-push TAG=<tag>` in the example directory itself.

Note: before running the script, ensure that through the CLI, you are logged into quay.io.

### Deploying

After confirming that your changes pass local testing:
1. Clean up testing artifacts
2. Create a PR. Once your PR has been merged, a Github Actions workflow (`Docker Publish`) will be triggered, to build, tag (with `stable`), and push
all example images. This ensures that all example images are using the most up-to-date version of the SDK, i.e. the one including your changes
3. If your SDK changes included modifications to any files in `pkg/info` or `pkg/apis/proto`, it is necessary 
to update the `go.mod` [file](https://github.com/numaproj/numaflow/blob/main/go.mod) in the Numaflow repo. This is because `numaflow-go` is a dependency of the Numaflow platform, and the files
in these directories are imported and used by Numaflow. Thus, get the commit SHA
of the merged PR from the previous step, and in the Numaflow repo run:
    ```shell
    go get github.com/numaproj/numaflow-go@<commit-sha>
   ```
   Followed by:
   ```shell
   go mod tidy
   ```
   Create a PR for these changes

### Adding a New Example

If you add a new example, in order for it to be used by the `Docker Publish` workflow, add its path
to the `dockerfile_paths` matrix in `.github/workflows/build-push.yaml`.