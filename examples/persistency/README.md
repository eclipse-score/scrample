
# Getting Started with Key-Value-Storage (persistency)

This guide will help you get started with the C++ and Rust implementations of the Key-Value-Storage (KVS) library, including how to use it and link it with Bazel.

## 1. Integrating with Bazel

### 1.1 Add this to your MODULE.bazel:
<details>
  <summary>MODULE.bazel</summary>

    module(name = "your_project_name")

    # Add the persistency dependency (replace version as needed)
    bazel_dep(name = "persistency", version = "0.2.0")

    # Add required toolchains and dependencies for C++ and Rust
    bazel_dep(name = "score_toolchains_gcc", version = "0.4", dev_dependency=True)
    gcc = use_extension("@score_toolchains_gcc//extentions:gcc.bzl", "gcc", dev_dependency=True)
    gcc.toolchain(
        url = "https://github.com/eclipse-score/toolchains_gcc_packages/releases/download/0.0.1/x86_64-unknown-linux-gnu_gcc12.tar.gz",
        sha256 = "457f5f20f57528033cb840d708b507050d711ae93e009388847e113b11bf3600",
        strip_prefix = "x86_64-unknown-linux-gnu",
    )
    use_repo(gcc, "gcc_toolchain", "gcc_toolchain_gcc")

    bazel_dep(name = "rules_rust", version = "0.61.0")
    crate = use_extension("@rules_rust//crate_universe:extensions.bzl", "crate")
    crate.from_specs(name = "crate_index")
    use_repo(crate, "crate_index")

    # Add any other dependencies required by persistency (see persistency's own MODULE.bazel for details)

</details>

### 1.2 Insert this into your .bazelrc:
<details>
  <summary>.bazelrc</summary>

  ```
  build --@score_baselibs//score/json:base_library=nlohmann
  build --@score_baselibs//score/mw/log/flags:KRemote_Logging=False

  ```
</details>

### 1.3 Run Bazel
If you start with a plain project add an empty file called `BUILD` into your project folder.

Now you can build the project with the command:
```sh
bazel build //...
```
(So far nothing happens, because no targets were defined.)

You can now continue in this guide to create a simple consumer-producer program or start on your own.



## 2. Using the C++ Implementation

### 2.1 Basic Usage

The C++ API is centered around `KvsBuilder` and `Kvs` classes. Here is a minimal example based on the test scenarios:

```cpp
#include "kvsbuilder.hpp"
#include <iostream>

int main() {
  // Use fully qualified names instead of 'using namespace'
  auto open_res = score::mw::per::kvs::KvsBuilder(score::mw::per::kvs::InstanceId(0))
            .need_defaults_flag(true)
            .need_kvs_flag(false)
            .dir(".")
            .build();
  if (!open_res) {
    std::cerr << "Failed to open KVS: " << static_cast<int>(open_res.error()) << std::endl;
    return 1;
  }
  score::mw::per::kvs::Kvs kvs = std::move(open_res.value());

  // Set a key-value pair
  kvs.set_value("username", score::mw::per::kvs::KvsValue("alice"));

  // Read a value (will fall back to default if not set)
  score::mw::per::kvs::Result<score::mw::per::kvs::KvsValue> get_res = kvs.get_value("username");
  if (get_res) {
    std::cout << "username: " << get_res.value().as_string() << std::endl;
  }

  // Read a default value (not set, but present in defaults)
  auto get_default = kvs.get_value("language");
  if (get_default) {
    std::cout << "language (default): " << get_default.value().as_string() << std::endl;
  }

  // Check if a key exists (only if written)
  if (kvs.key_exists("username").value_or(false)) {
    std::cout << "username exists!" << std::endl;
  }

  // Remove a key
  kvs.remove_key("username");

  // List all keys (only written keys, not defaults)
  auto keys_res = kvs.get_all_keys();
  if (keys_res) {
    std::cout << "All keys in KVS:" << std::endl;
    for (const auto& key : keys_res.value()) {
      std::cout << "  " << key << std::endl;
    }
  }

  // Flush changes to disk
  kvs.flush();

  return 0;
}
```

## 3. Using the Rust Implementation

### 3.1 Basic Usage

From `examples/basic.rs`:
```rust
use rust_kvs::prelude::*;
use std::collections::HashMap;
use tempfile::tempdir;

fn main() -> Result<(), ErrorCode> {
    let dir = tempdir()?;
    let dir_string = dir.path().to_string_lossy().to_string();
    let instance_id = InstanceId(0);

    // Build KVS instance
    let builder = KvsBuilder::new(instance_id)
        .dir(dir_string.clone())
        .kvs_load(KvsLoad::Optional);
    let kvs = builder.build()?;

    // Set values
    kvs.set_value("number", 123.0)?;
    kvs.set_value("bool", true)?;
    kvs.set_value("string", "First")?;

    // Get value
    let value = kvs.get_value("number")?;
    println!("number = {:?}", value);

    Ok(())
}
```

### 3.2 Snapshots Example

From `examples/snapshots.rs`:
```rust
let max_count = kvs.snapshot_max_count() as u32;
for index in 0..max_count {
    kvs.set_value("counter", index)?;
    kvs.flush()?;
    println!("Snapshot count: {:?}", kvs.snapshot_count());
}

// Restore a snapshot
kvs.snapshot_restore(SnapshotId(2))?;
```

### 3.3 Defaults Example

From `examples/defaults.rs`:
```rust
// Create defaults file and build KVS instance with defaults
create_defaults_file(dir.path().to_path_buf(), instance_id)?;
let builder = KvsBuilder::new(instance_id)
    .dir(dir_string)
    .defaults(KvsDefaults::Required);
let kvs = builder.build()?;

// Get default value
let k1_value = kvs.get_default_value("k1")?;
println!("k1 = {:?}", k1_value);
```
## 4. Default value file example
This file should be placed in the working directory:
```json
{
    "language": "en",
    "theme": "dark",
    "timeout": 30
}
```

**Important:**
- If you open the KVS with `.need_defaults_flag(true)`, the file must exist.
- The KVS will use these defaults for any key not explicitly set.
- You must also provide a CRC file (e.g., `defaults.json.crc`) alongside the defaults file. This CRC file is generated using the Adler-32 checksum algorithm, as implemented in the codebase. The CRC ensures the integrity of the defaults file at runtime.
## 5. More Examples
- See `src/cpp/tests/` for C++ test scenarios and usage patterns.
- See `src/rust/rust_kvs/examples/` for Rust usage patterns.
