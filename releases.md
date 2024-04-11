# Release Guide

This document explains the release process for the Go SDK. You can find the most recent version under [Github Releases](https://github.com/numaproj/numaflow-go/releases).

### Before Release

If the version to be released has backwards incompatible changes, i.e. it does not support older versions of the Numaflow platform,
you must update the `MinimumNumaflowVersion` constant in the `pkg/info/types.go` file to the minimum Numaflow version that is supported by your new SDK version.
Ensure that this change is merged and included in the release.

### How to Release

This can be done via the Github UI. In the `Releases` section of the Go SDK repo, click `Draft a new release`. Create an appropriate tag for the version number, using the [semantic versioning](https://semver.org/) specification, and select it. Make 
the title the same as the tag. Click `Generate release notes` so that all the changes made since the last release are documented. If there are any major features or breaking
changes that you would like to highlight as part of the release, add those to the description as well. Then set the release as either pre-release or latest, depending
on your situation. Finally, click `Publish release`, and your version tag will be the newest release on the repository.

### After Release

Once a new version has been released, and its corresponding version tag exists on the remote repo, you want to update the `go.mod`
files to reflect this new version:
```shell
./hack/update_examples.sh -u <version-tag>
  ```
After running the above, create a PR for the changes that the script made. Once merged, it will trigger the Docker Publish workflow.
As a result, the correct SDK version will always be printed in the server information logs,
and the example images will always be using the latest changes (due to the local references).