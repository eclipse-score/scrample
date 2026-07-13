# SCRAMPLE

[![Build Status](https://github.com/eclipse-score/scrample/actions/workflows/build.yml/badge.svg)](https://github.com/eclipse-score/scrample/actions/workflows/build.yml)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

**SCRAMPLE** (S-CORE + Sample) is a demo application showcasing the logging capabilities of the [Eclipse S-CORE](https://projects.eclipse.org/projects/automotive.score) platform for Software Defined Vehicles (SDVs).

## Overview

SCRAMPLE provides two minimal "Hello World" applications — one in C++ and one in Rust — that both use the S-CORE `mw::log` logging library. They share a common logging configuration and produce identical output, demonstrating how to integrate `mw::log` in both languages.

## Prerequisites

- **Bazel** 8.3.0 or higher (see `.bazelversion`)
- **QNX SDP** (for cross-compilation to QNX targets)
- **C++17** compatible compiler
- **Rust toolchain** (automatically managed by Bazel via Ferrocene)
- **Dependencies** (automatically managed via Bazel):
  - S-CORE Base Libraries (`score_baselibs`)
  - S-CORE Logging (`score_logging`)
  - S-CORE Base Libraries Rust (`score_baselibs_rust`)

## Toolchains

SCRAMPLE uses the same toolchain setup as the other S-CORE modules.

### C++ Toolchains

Toolchains are provided by [`score_bazel_cpp_toolchains`](https://github.com/eclipse-score/bazel_cpp_toolchains) via the `gcc` Bzlmod extension:

| Toolchain | Target | Used for |
|---|---|---|
| `score_gcc_x86_64_toolchain` | x86_64 Linux (GCC 12.2.0) | `--config=host` builds |
| `score_qcc_x86_64_toolchain` | x86_64 QNX SDP 8.0.0 | `--config=x86_64-qnx` builds |

### Rust Toolchains

Rust toolchains are provided by [`score_toolchains_rust`](https://github.com/eclipse-score/toolchains_rust), which bundles pre-built [Ferrocene](https://ferrous-systems.com/ferrocene/) toolchains:

| Toolchain | Target | Used for |
|---|---|---|
| `ferrocene_x86_64_unknown_linux_gnu` | x86_64 Linux | Host builds + proc macro compilation |
| `ferrocene_x86_64_pc_nto_qnx800` | x86_64 QNX Neutrino 8.0.0 | `--config=x86_64-qnx` Rust builds |

The Ferrocene Linux toolchain is registered as a `common` toolchain so that proc macro crates (compiled on the host) and QNX target crates share the same compiler metadata format, which is required for cross-compilation compatibility.

### Key `.bazelrc` Flags

| Flag | Purpose |
|---|---|
| `--@score_baselibs//score/memory/shared/flags:use_typedshmd=False` | Disables proprietary Shared Memory Data Router (not available in public CI) |
| `--@score_baselibs//score/mw/log/flags:KRemote_Logging=False` | Disables remote logging backend (not needed for this demo) |

## Building

### Build the C++ app
```bash
bazel build --config=host //src_cpp:scrample_cpp
```

### Build the Rust app
```bash
bazel build --config=host //src_rust:scrample_rust
```

### Build both apps
```bash
bazel build --config=host //src_cpp:scrample_cpp //src_rust:scrample_rust
```

### QNX Cross-Compilation
```bash
bazel build --config=x86_64-qnx //src_cpp:scrample_cpp
bazel build --config=x86_64-qnx //src_rust:scrample_rust
```

**Note:** Always use a build configuration (`--config=host` or `--config=x86_64-qnx`) to ensure proper dependency settings.

## Running

After building with `--config=host`:

### C++ app
```bash
./bazel-bin/src_cpp/scrample_cpp
```

### Rust app
```bash
./bazel-bin/src_rust/scrample_rust
```

Both apps produce the same log output to the console:
```
2026/05/22 12:00:00.0000000 00000000 000 ECU1 SCRM scra log info verbose 1 Hello from SCRAMPLE!
```

## Project Structure

```
scrample/
├── src_cpp/
│   ├── main.cpp        # C++ app using score::mw::log
│   └── BUILD
├── src_rust/
│   ├── main.rs         # Rust app using score_log + score_log_bridge
│   └── BUILD
├── config/
│   └── logging.json    # Shared logging configuration (console, kInfo)
└── BUILD               # Bazel build definitions
```

## Development

### Code Formatting
Apply automatic formatting fixes:
```bash
bazel run //:format.fix
```

### Check Formatting
Check if code formatting is correct:
```bash
bazel test //:format.check
```

### Copyright Checking
Verify copyright headers are present:
```bash
bazel run //:copyright.check
```

**Note:** Formatting commands don't require `--config` flags.

## Contributing

SCRAMPLE is part of the Eclipse S-CORE project. Contributions are welcome!

1. Read the [Contributing Guide](CONTRIBUTION.md)
2. Sign the [Eclipse Contributor Agreement (ECA)](https://www.eclipse.org/legal/ECA.php)
3. Follow the [Developer Certificate of Origin (DCO)](https://www.eclipse.org/legal/dco/)
4. Submit pull requests via GitHub

For questions and discussions:
- Mailing list: [score-dev](https://accounts.eclipse.org/mailing-list/score-dev)
- Chat: [Eclipse S-CORE Matrix Room](https://chat.eclipse.org/#/room/#automotive.score:matrix.eclipse.org)

## License

This project is licensed under the [Apache License 2.0](LICENSE).

Copyright © 2025 Contributors to the Eclipse Foundation.

## Related Projects

- [Eclipse S-CORE](https://projects.eclipse.org/projects/automotive.score) - Main project
- [S-CORE Documentation](https://eclipse-score.github.io) - Full platform documentation
- [S-CORE Communication](https://github.com/eclipse-score/communication) - IPC middleware
- [S-CORE Base Libraries](https://github.com/eclipse-score/baselibs) - Core utilities

## Troubleshooting

### Runtime Warnings
When running the application, you may see:
```
mw::log initialization error: Error No logging configuration files could be found.
Fallback to console logging.
```
This is expected and harmless. The application falls back to console logging when the optional logging configuration isn't found at the expected system location.

### Build Warnings
You may see deprecation warnings during compilation related to:
- `string_view` null-termination checks
- `InstanceSpecifier::Create()` API deprecations

These are intentional warnings from the S-CORE libraries and do not prevent successful builds. They are addressed in the `.bazelrc` configuration with `-Wno-error=deprecated-declarations`.

### Build Configuration Required
Always use a build configuration (`--config=host` or `--config=x86_64-qnx`). Building without a config flag will fail with missing dependency errors because the required S-CORE library flags (like `tracing_library=stub`) won't be set.

### QNX Builds
QNX cross-compilation requires:
- QNX SDP installation and license
- Proper credential setup (see `.github/workflows/build.yml` for CI example)

The required `score_baselibs` feature flags (`use_typedshmd=False`, `KRemote_Logging=False`) are already configured in `.bazelrc` and apply automatically.

## Roadmap

Future extensions planned for SCRAMPLE:

- Additional S-CORE platform module demonstrations
- More complex communication patterns
- Performance benchmarking utilities
- Integration with other S-CORE components
