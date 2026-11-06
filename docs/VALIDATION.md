# Validation

Authoring date: November 6, 2026. Host: Windows. Target: Linux cgroup v2; synthetic fixtures also work on other Go platforms.

## Local result

Go/Linux unavailable on authoring host. CI checks fixture parsing, traversal, and live root reads.

## Reproducible checks

The GitHub Actions workflow builds/runs the following on `ubuntu-latest`:

```sh
gofmt -w main.go main_test.go
go vet ./...
go test -race ./...
go build -o cgroup-explorer .
./cgroup-explorer --depth 0 --json
```

## Manual acceptance

Compare a known cgroup against its files, check max/unlimited quotas, inspect absent controller fields, and confirm depth bounds on a synthetic or live hierarchy.

A workflow file is not evidence of a passing run. Observed CI run links and source commits are recorded after publication. Native GUI, sound, sleep/resume, and platform permissions require direct target-platform acceptance even when model tests pass.

## Observed CI result

Passed: [target-platform build/check](https://github.com/netcodermike/linux-cgroup-explorer/actions/runs/42204070153). Source commit: `4567a0bfd8019b2b1ce9dca1e345eaf782bd528f`. Later documentation-only changes do not alter the tested implementation. The reproducible commands above identify test scope; native GUI/audio/permissions are not implied by compilation.
