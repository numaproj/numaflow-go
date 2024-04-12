# Release Guide

This document explains the release process for the Go SDK. You can find the most recent version under [Github Releases](https://github.com/numaproj/numaflow-go/releases).

### Before Release

If the version to be released has backwards incompatible changes, i.e. it does not support older versions of the Numaflow platform,
you must update the `MinimumNumaflowVersion` constant in the `pkg/info/types.go` file to the minimum Numaflow version that is supported by your new SDK version.
Ensure that this change is merged and included in the release.

### How to Release

This can be done via the Github UI. 
1. In the `Releases` section of the Go SDK repo, click `Draft a new release`
2. Create an appropriate tag for the version number following the [semantic versioning](https://semver.org/) specification and select it. The tag should be prefixed with a `'v'`
3. Make the title the same as the tag 
4. Click `Generate release notes` so that all the changes made since the last release are documented. If there are any major features or breaking
changes that you would like to highlight as part of the release, add those to the description as well 
5. Select `Set as the latest release` or `Set as a pre-release`, depending on your situation 
6. Finally, click `Publish release`, and your version tag will be the newest release on the repository

### After Release

1. Once a new version has been released, and its corresponding version tag exists on the remote repo, you want to update the `go.mod`
files to reflect this new version:
    ```shell
    ./hack/update_examples.sh -u <version-tag>
    ```
2. Create a PR for the changes that the script made. Once merged, it will trigger the `Docker Publish` workflow.
As a result, the correct SDK version will always be printed in the server information logs,
and the example images will always be using the latest changes (due to the local references)
3. Since the Go SDK is a dependency of Numaflow, once a new version has been released, the `go.mod` [file](https://github.com/numaproj/numaflow/blob/main/go.mod)
in the Numaflow repo should be updated
   - If you set the SDK release as latest, in your Numaflow repo run:
        ```shell
         go get github.com/numaproj/numaflow-go@latest
        ``` 
        Followed by:
        ```shell
        go mod tidy
        ```
        Create a PR for these changes
   - If you set the SDK release as pre-release, in your Numaflow repo run:
     ```shell
      go get github.com/numaproj/numaflow-go@<commit-sha>
      ``` 
     Where `commit-sha` is the latest commit SHA in the SDK repo. Then run:
     ```shell
     go mod tidy
     ```
     Create a PR for these changes