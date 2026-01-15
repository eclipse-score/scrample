<div align="center">

```text
|=================================|
|   ___  ___ ___  _ __ _____  __  |
|  / __|/ __/ _ \| '__/ _ \ \/ /  |
|  \__ \ (_| (_) | | |  __/>  <   |
|  |___/\___\___/|_|  \___/_/\_\  |
|                                 |
|=================================|
```

</div>

# scorex

`scorex` is a small CLI helper for generating S-CORE skeleton applications.

It is implemented in Go in [scorex/main.go](scorex/main.go) and uses Cobra for its CLI in
[scorex/cmd/root.go](scorex/cmd/root.go) and [scorex/cmd/init.go](scorex/cmd/init.go).

## Features

- Generate a new S-CORE Bazel project skeleton
- Pre-wire `MODULE.bazel`, `.bazelrc`, `.bazelversion`, `BUILD`, and `src/main.cpp`
- Use a central `known_good.json` to pin module versions and commits

The project layout and files are rendered from the templates in
[scorex/cmd/templates/application](scorex/cmd/templates/application).

## Installation

From the repository root:

```sh
cd scorex
go mod tidy
go build ./...
```

This creates a `scorex` binary in the `scorex/` directory.

## Usage

Show help:

```sh
./scorex --help
```

Generate a new S-CORE project (example):

```sh
./scorex init \
  --module score_baselibs \
  --module score_communication \
  --name my_score_app \
  --dir . \
  --bazel-version 8.3.0
```

This will create `./my_score_app` with:

- `MODULE.bazel`
- `.bazelrc`
- `.bazelversion`
- `BUILD`
- `src/BUILD`
- `src/main.cpp`

## Options

The `init` command (see [scorex/cmd/init.go](scorex/cmd/init.go)) supports:

- `--module` (repeatable): S-CORE modules to include, e.g. `score_communication`
- `--name`: Name of the generated project (default: `score_app`)
- `--dir`: Target directory where the project is created (default: current directory)
- `--known-good-url`: URL or file path to `known_good.json`
- `--bazel-version`: Bazel version written into `.bazelversion` (default: `8.3.0`)
