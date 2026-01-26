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

## Distribution

The `scorex` CLI can be distributed through multiple methods.

### Current Installation Methods

#### Manual Download (Available Now)

Download the appropriate binary for your platform from the [releases page](https://github.com/eclipse-score/score_scrample/releases):

- **Linux (x86_64)**: `scorex-VERSION-linux-x86_64.tar.gz`
- **macOS (Apple Silicon)**: `scorex-VERSION-macos-arm64.tar.gz`
- **macOS (Intel)**: `scorex-VERSION-macos-x86_64.tar.gz`
- **Windows (x86_64)**: `scorex-VERSION-windows-x86_64.zip`

**Installation steps:**

1. Download the appropriate archive for your platform
2. Extract it:
   ```bash
   # macOS/Linux
   tar -xzf scorex-VERSION-platform.tar.gz

   # Windows (PowerShell)
   Expand-Archive scorex-VERSION-windows-x86_64.zip
   ```
3. Move the binary to a directory in your PATH:
   ```bash
   # macOS/Linux
   sudo mv scorex-platform /usr/local/bin/scorex
   sudo chmod +x /usr/local/bin/scorex

   # Windows - move to a directory in your PATH or add the directory to PATH
   ```
4. Verify installation:
   ```bash
   scorex version
   ```

#### Universal Install Script (Available Now)

**macOS & Linux:**
```bash
curl -sSL https://raw.githubusercontent.com/eclipse-score/score_scrample/main/scorex/distribution/install.sh | sh
```

This script automatically:
- Detects your OS and architecture
- Downloads the correct binary
- Installs it to `/usr/local/bin/scorex`
- Makes it executable

### Package Manager Installation

#### macOS & Linux - Homebrew

```bash
# Install directly from this repository (no tap required)
brew install https://raw.githubusercontent.com/eclipse-score/score_scrample/main/scorex/distribution/homebrew/scorex.rb
```

**Note**: The formula checksums need to be updated manually after each release. See the maintainer section below.

#### Windows - Scoop

```bash
# Install directly from this repository (no bucket required)
scoop install https://raw.githubusercontent.com/eclipse-score/score_scrample/main/scorex/distribution/scoop/scorex.json
```

**Note**: The manifest checksums need to be updated manually after each release. See the maintainer section below.

### For Maintainers

#### Publishing a New Release

1. **Create and push a version tag:**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **GitHub Actions automatically:**
   - Builds binaries for all platforms (Linux x86_64, macOS ARM64, macOS Intel, Windows x86_64)
   - Creates compressed archives (.tar.gz for Unix, .zip for Windows)
   - Generates `checksums.txt` with SHA256 hashes
   - Creates a GitHub release with separate artifacts per platform
   - Uploads artifacts: `scorex-linux`, `scorex-macos`, `scorex-windows`

3. **Update package manager manifests:**
   - Download `checksums.txt` from the GitHub release
   - Update `distribution/homebrew/scorex.rb`:
     - Set `version` to the new version (without the `v` prefix)
     - Update the three `sha256` values (Linux, macOS ARM64, macOS Intel) from checksums.txt
   - Update `distribution/scoop/scorex.json`:
     - Set `version` to the new version (without the `v` prefix)
     - Update the `hash` value for Windows from checksums.txt
   - Commit and push these changes to the main repository

4. **Users can now install** via direct URLs, install script, or manual download
