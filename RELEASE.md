# Release

Releases are automated from semver git tags.

To make a release, create and push a semver tag from the commit you want to release:

```sh
git tag v1.2.3
git push origin v1.2.3
```

The `Publish tag` GitHub Actions workflow runs on `v*` tags and validates the tag format before publishing. Use `vMAJOR.MINOR.PATCH`, such as `v1.2.3`, or a prerelease tag like `v1.2.3-rc.1`.

The workflow publishes these multi-arch images:

```text
docker.io/kuberhealthy/test-check:<tag>
ghcr.io/kuberhealthy/test-check:<tag>
```

The image manifest includes:

- `linux/amd64`
- `linux/arm64`

After the image is pushed, the workflow updates the example healthcheck YAML on `main` to use the released Docker Hub image tag.

The workflow then creates or updates a GitHub release with the same semver as the tag. The release notes link to the GitHub package, and the release includes a `release-images.txt` asset listing the image and supported platforms.

Docker image tags do not support `+`, so do not use semver build metadata in release tags.
