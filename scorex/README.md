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

### Installation Methods

#### Install Script (Recommended)

**macOS & Linux:**
```bash
curl -sSL https://raw.githubusercontent.com/eclipse-score/score_scrample/main/scorex/distribution/install.sh | sh
```

**Note**: Requires a published release. If no releases exist yet, use manual installation below.

This script automatically:
- Detects your OS and architecture
- Downloads the correct binary from the latest release
- Installs it to `$HOME/.local/bin/scorex` (customizable via `SCOREX_INSTALL_DIR`)
- Makes it executable
- Removes macOS quarantine attribute automatically

**Add to PATH** (if not already):
```bash
export PATH="$HOME/.local/bin:$PATH"
```

Add this to your `~/.zshrc` or `~/.bashrc` to make it permanent.

#### Manual Download

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
4. **macOS only**: Remove the quarantine attribute (required for unsigned binaries):
   ```bash
   sudo xattr -d com.apple.quarantine /usr/local/bin/scorex
   ```
   
   Alternatively, on first run, right-click the binary in Finder and select "Open" to bypass Gatekeeper.

5. Verify installation:
   ```bash
   scorex version
   ```

**Note for macOS users**: The binaries are currently unsigned. You may see a security warning. Use the `xattr` command above or right-click > Open to bypass Gatekeeper.

### For Maintainers

#### Testing the Release Flow in PRs

The release workflow runs on PRs and creates artifacts. To test the binaries:

1. **Go to the PR's Actions tab** and find the latest "Release scorex CLI" workflow run
2. **Download the artifact** for your platform:
   - `scorex-linux` - Contains Linux binary
   - `scorex-macos` - Contains macOS binaries (ARM64 + Intel)
   - `scorex-windows` - Contains Windows binary

3. **Extract and test locally**:
   ```bash
   # Download artifact from GitHub Actions UI
   unzip scorex-macos.zip
   
   # Extract the tar.gz
   tar -xzf scorex-pr-*-macos-arm64.tar.gz
   
   # Make executable and remove quarantine (macOS)
   chmod +x scorex-macos-arm64
   xattr -d com.apple.quarantine scorex-macos-arm64 2>/dev/null || true
   
   # Test it
   ./scorex-macos-arm64 --help
   ./scorex-macos-arm64 version
   
   # Test creating a project
   ./scorex-macos-arm64 init --name test_app --dir /tmp/test_scorex
   ```

**Note**: The install script cannot be tested in PRs since it requires a published GitHub release.

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
   - Creates a GitHub release with all artifacts
   - For PRs and manual runs: Uploads separate artifacts per platform (`scorex-linux`, `scorex-macos`, `scorex-windows`)

3. **Users can install** via the install script or manual download from the releases page
