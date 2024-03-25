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
2. Create a PR. Once your PR has been merged, a Github Actions workflow (Docker Publish) will be triggered, to build, tag (with `stable`), and push
   all example images. This ensures that all example images are using the most up-to-date version of the SDK, i.e. the one including your
   changes.

### After release

Once a new version has been released, and its corresponding version tag exists on the remote repo, you want to update the `go.mod` 
files to reflect this new version:
```shell
./hack/update_examples.sh -u <version>
  ```
After running the above, create a PR for the changes that the script made. Once merged, it will trigger the Docker Publish workflow.
As a result, the correct SDK version will always be printed in the server information logs, 
and the example images will always be using the latest changes (due to the local references).

### Adding a New Example

If you add a new example, in order for it to be used by the Docker Publish workflow, add the path to 
the example to the `dockerfile_paths` matrix in `build-push.yaml`.