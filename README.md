# SCRAMPLE

[![Build Status](https://github.com/eclipse-score/scrample/actions/workflows/build.yml/badge.svg)](https://github.com/eclipse-score/scrample/actions/workflows/build.yml)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

**SCRAMPLE** (S-CORE + Sample) is a demonstration application showcasing the inter-process communication (IPC) capabilities of the [Eclipse S-CORE](https://projects.eclipse.org/projects/automotive.score) platform for Software Defined Vehicles (SDVs).

## Overview

This application demonstrates a producer-consumer pattern using S-CORE's middleware communication layer. It illustrates:

- **Event-based communication** using the S-CORE middleware (`score::mw::com`)
- **Shared memory IPC** for high-performance data transfer between processes
- **Type-safe serialization** of complex automotive data structures
- **Skeleton-Proxy pattern** following AUTOSAR Adaptive Platform concepts

The sample exchanges `MapApiLanesStamped` messages containing lane information data structures, simulating real-world automotive HD map data exchange scenarios.

## Architecture

SCRAMPLE consists of two operational modes:

### Skeleton (Publisher/Server)
- Offers a service instance identified by `score/MapApiLanesStamped`
- Generates and publishes lane data samples at a configurable cycle rate
- Populates complex nested data structures with randomized values
- Computes hash values for data integrity verification

### Proxy (Subscriber/Client)
- Discovers and connects to the skeleton service
- Subscribes to receive lane data events
- Validates received samples for ordering and data integrity
- Supports both event-driven callbacks and polling modes

## Prerequisites

- **Bazel** 8.3.0 or higher (see `.bazelversion`)
- **QNX SDP** (for cross-compilation to QNX targets)
- **C++17** compatible compiler
- **Rust toolchain** (automatically managed by Bazel for Rust tests)
- **Dependencies** (automatically managed via Bazel):
  - S-CORE Base Libraries (`score_baselibs`)
  - S-CORE Communication (`score_communication`)
  - Boost.Program_Options
  - GoogleTest (for C++ tests)
  - Rust toolchain (for Rust tests)

## Building

### Standard Build (Host Platform)
```bash
bazel build --config=host //src:scrample
```

### QNX Cross-Compilation
```bash
bazel build --config=x86_64-qnx //src:scrample
```

### Build All Tests
```bash
bazel test --config=host //tests/...
```

**Note:** Always use a build configuration (`--config=host` or `--config=x86_64-qnx`) to ensure proper dependency settings.

## Running

After building with `--config=host`, the binary will be in `bazel-bin/src/scrample`.

### Quick Start (Two Terminals)

To see the IPC communication in action, open two terminals:

**Terminal 1 - Start Skeleton (Publisher):**
```bash
./bazel-bin/src/scrample \
  --mode skeleton \
  --cycle-time 1000 \
  --num-cycles 10 \
  --service_instance_manifest src/etc/mw_com_config.json
```

**Terminal 2 - Start Proxy (Subscriber):**
```bash
./bazel-bin/src/scrample \
  --mode proxy \
  --cycle-time 500 \
  --num-cycles 20 \
  --service_instance_manifest src/etc/mw_com_config.json
```

You should see the proxy discover the skeleton service, subscribe, and receive `MapApiLanesStamped` samples. The proxy validates data integrity and ordering for each received sample.

### Start Skeleton (Publisher)
```bash
./bazel-bin/src/scrample \
  --mode skeleton \
  --cycle-time 1000 \
  --num-cycles 10 \
  --service_instance_manifest src/etc/mw_com_config.json
```

### Start Proxy (Subscriber)
```bash
./bazel-bin/src/scrample \
  --mode proxy \
  --cycle-time 500 \
  --num-cycles 20 \
  --service_instance_manifest src/etc/mw_com_config.json
```

### Command-Line Options

| Option | Description | Required |
|--------|-------------|----------|
| `--mode, -m` | Operation mode: `skeleton`/`send` or `proxy`/`recv` | Yes |
| `--cycle-time, -t` | Cycle time in milliseconds for sending/polling | Yes |
| `--num-cycles, -n` | Number of cycles to execute (0 = infinite) | Yes |
| `--service_instance_manifest, -s` | Path to communication config JSON | Optional |
| `--disable-hash-check, -d` | Skip sample hash validation in proxy mode | Optional |

## Configuration

The communication behavior is configured via `src/etc/mw_com_config.json`:

- Service type definitions and bindings
- Event definitions with IDs
- Instance-specific configuration (shared memory settings, subscriber limits)
- ASIL level and process ID restrictions

## Project Structure

```
scrample/
├── src/
│   ├── main.cpp                  # Entry point and CLI argument parsing
│   ├── sample_sender_receiver.cpp # Core skeleton/proxy logic
│   ├── datatype.h                # Data type definitions
│   ├── assert_handler.cpp        # Custom assertion handling
│   └── etc/
│       ├── mw_com_config.json    # Communication configuration
│       └── logging.json          # Logging configuration
├── tests/
│   ├── cpp/                      # C++ unit tests (GoogleTest)
│   └── rust/                     # Rust tests
├── docs/                         # Sphinx documentation
└── BUILD                         # Bazel build definitions
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

### Build Documentation
Build Sphinx documentation:
```bash
bazel build //:docs
```

**Note:** Formatting and documentation commands don't require `--config` flags.

## Testing

The project includes example tests demonstrating the testing infrastructure:

- **C++ Tests**: GoogleTest-based unit tests in `tests/cpp/`
- **Rust Tests**: Rust test framework in `tests/rust/`

Run all tests:
```bash
bazel test --config=host //tests/...
```

Run all tests and format checks:
```bash
bazel test --config=host //tests/... //:format.check
```

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

## Roadmap

Future extensions planned for SCRAMPLE:

- Additional S-CORE platform module demonstrations
- More complex communication patterns
- Performance benchmarking utilities
- Integration with other S-CORE components
