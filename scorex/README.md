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

The `scorex` CLI is distributed through multiple package managers for easy installation across different platforms.

### Installation Methods

#### macOS & Linux - Homebrew

```bash
# Add the tap (once a tap repository is created)
brew tap eclipse-score/tap

# Install scorex
brew install scorex
```

#### Windows - Scoop

```bash
# Add the bucket (once a bucket repository is created)
scoop bucket add eclipse-score https://github.com/eclipse-score/scoop-bucket

# Install scorex
scoop install scorex
```

#### Universal - Install Script

**macOS & Linux:**
```bash
curl -sSL https://raw.githubusercontent.com/eclipse-score/score_scrample/main/scorex/distribution/install.sh | sh
```

#### Manual Download

Download the appropriate binary for your platform from the [releases page](https://github.com/eclipse-score/score_scrample/releases):

- **Linux (x86_64)**: `scorex-VERSION-linux-x86_64.tar.gz`
- **macOS (Apple Silicon)**: `scorex-VERSION-macos-arm64.tar.gz`
- **macOS (Intel)**: `scorex-VERSION-macos-x86_64.tar.gz`
- **Windows (x86_64)**: `scorex-VERSION-windows-x86_64.zip`

Extract and move to a directory in your PATH.

### For Maintainers

#### Publishing a New Release

1. Create and push a new version tag:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. GitHub Actions will automatically:
   - Build binaries for all platforms
   - Create compressed archives
   - Generate checksums
   - Create a GitHub release

3. Update package manifests:

   **Homebrew Formula** (`distribution/homebrew/scorex.rb`):
   - Update version number
   - Update SHA256 checksums from `checksums.txt` in the release

   **Scoop Manifest** (`distribution/scoop/scorex.json`):
   - Update version number
   - Update SHA256 hash from `checksums.txt`

4. Commit and push updated manifests to respective repositories:
   - Homebrew: Create/update tap repository at `eclipse-score/homebrew-tap`
   - Scoop: Create/update bucket repository at `eclipse-score/scoop-bucket`

#### Setting Up Package Repositories

**Homebrew Tap:**
1. Create repository: `https://github.com/eclipse-score/homebrew-tap`
2. Add `distribution/homebrew/scorex.rb` to the repository root or `Formula/` directory
3. Users can then install with: `brew install eclipse-score/tap/scorex`

**Scoop Bucket:**
1. Create repository: `https://github.com/eclipse-score/scoop-bucket`
2. Add `distribution/scoop/scorex.json` to the `bucket/` directory
3. Users can then install with: `scoop bucket add eclipse-score <repo-url>` then `scoop install scorex`

#### Updating Checksums

After each release, download `checksums.txt` from the GitHub release and update:

```bash
# Example for version 1.0.0
curl -sL https://github.com/eclipse-score/score_scrample/releases/download/v1.0.0/checksums.txt

# Update the SHA256 values in:
# - distribution/homebrew/scorex.rb
# - distribution/scoop/scorex.json
```
