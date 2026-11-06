# Linux Cgroup Explorer

Inspect local cgroup memory, PID, and CPU budgets with bounded hierarchy traversal and structured JSON output.

[![Build and verify](https://github.com/netcodermike/linux-cgroup-explorer/actions/workflows/verify.yml/badge.svg)](https://github.com/netcodermike/linux-cgroup-explorer/actions/workflows/verify.yml)

**Language:** Go · **Platform:** Linux cgroup v2; synthetic fixtures also work on other Go platforms · **App type:** Linux container resource explorer

First-release learning implementation, with explicit scope below. This is one independently runnable app, not a placeholder in a combined repository.

## Features

- Current/configured memory and PID values.
- CPU quota in core equivalents and cumulative usage microseconds.
- Depth-limited traversal up to 2,000 groups.
- JSON nulls/read errors distinguish missing values from zero.

## Requirements

Go 1.22+ and a readable cgroup v2 hierarchy exposing cgroup.controllers. No external modules or writes to controllers.

## Quick start

```sh
git clone https://github.com/netcodermike/linux-cgroup-explorer.git
cd linux-cgroup-explorer
go build -o cgroup-explorer .
./cgroup-explorer --depth 2
./cgroup-explorer --root /sys/fs/cgroup --depth 0 --json
```

Commands run from the repository root. Lines marked for a particular operating system apply only there. Language tools are prerequisites, not installed by the app. Close desktop windows normally; terminal/server apps stop with Ctrl+C unless their usage notes say otherwise.

## Example workflow

Inspect depth 0 for the root, then choose a readable child path and compare its displayed values to the controller files. Use --json to see precise per-field errors and distinguish unlimited quotas from absent files.

## Architecture, configuration, and data

`main.go` reads memory.current, memory.max, pids.current, pids.max, cpu.max, and cpu.stat for each selected directory. CPU quota is quota/period; max becomes null in JSON. Missing fields are recorded per group instead of invented as zero. --root defaults to /sys/fs/cgroup, --depth accepts 0–8, and --json emits the Group schema. Root cgroup.controllers is required by the CLI, while unit tests use temporary synthetic files. WalkDir does not follow child directory symlinks. `main_test.go` checks quota parsing, missing fields, and depth boundaries.

## Verification

Build with the quick-start command where needed, then run:

```sh
go test ./...
```

Manual acceptance: Compare a known cgroup against its files, check max/unlimited quotas, inspect absent controller fields, and confirm depth bounds on a synthetic or live hierarchy.

See [validation record](docs/VALIDATION.md) for observed checks and their limits. CI results verify only what the workflow runs; they do not establish native GUI correctness or production readiness.

## Limitations and next steps

cgroup v2 only, no v1, container runtime discovery, history, percentage sampling, process attribution, or controller changes. Reported limits are local values, not effective limits after applying all ancestor constraints. Cumulative CPU usage is not a current utilization percentage. Root/controller fields may legitimately be absent. Directory traversal and file reads are not atomic. Permission-denied subtrees can be omitted.

## Troubleshooting

A missing cgroup.controllers indicates the wrong root, a v1 hierarchy, or restricted access. n/a means a field could not be read/parsed. Check --json for details. Select a smaller root/depth if the 2,000-group cap is exceeded. Unlimited/unknown text intentionally covers values the table cannot fully distinguish; JSON contains errors.

## Source map

- `go.mod`
- `main.go`
- `main_test.go`

## Development and license

Read [contributing](CONTRIBUTING.md) before changing the app, and [security/privacy](SECURITY.md) before supplying private data. Use synthetic fixtures for tests. No telemetry, paid API key, or account is required by the app itself. Dependencies, when present, are stated in Requirements.

[MIT license](LICENSE), copyright 2026 netcodermike.
