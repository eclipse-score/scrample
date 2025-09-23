
# C++ & Rust Bazel Template Repository

This repository serves as a **template** for setting up **C++ and Rust projects** using **Bazel**.
It provides a **standardized project structure**, ensuring best practices for:

- **Build configuration** with Bazel.
- **Testing** (unit and integration tests).
- **Documentation** setup.
- **CI/CD workflows**.
- **Development environment** configuration.

---

## 📂 Project Structure

| File/Folder                         | Description                                       |
| ----------------------------------- | ------------------------------------------------- |
| `README.md`                         | Short description & build instructions            |
| `src/`                              | Source files for the module                       |
| `tests/`                            | Unit tests (UT) and integration tests (IT)        |
| `examples/`                         | Example files used for guidance                   |
| `docs/`                             | Documentation (Doxygen for C++ / mdBook for Rust) |
| `.github/workflows/`                | CI/CD pipelines                                   |
| `.vscode/`                          | Recommended VS Code settings                      |
| `.bazelrc`, `MODULE.bazel`, `BUILD` | Bazel configuration & settings                    |
| `project_config.bzl`                | Project-specific metadata for Bazel macros        |
| `LICENSE.md`                        | Licensing information                             |
| `CONTRIBUTION.md`                   | Contribution guidelines                           |

---

## 🚀 Getting Started

### 1️⃣ Clone the Repository

```sh
git clone https://github.com/eclipse-score/YOUR_PROJECT.git
cd YOUR_PROJECT
```

### 2️⃣ Build the Examples of module

> DISCLAIMER: Depending what module implements, it's possible that different
> configuration flags needs to be set on command line.

To build all targets of the module the following command can be used:

```sh
bazel build //src/...
```

This command will instruct Bazel to build all targets that are under Bazel
package `src/`. The ideal solution is to provide single target that builds
artifacts, for example:

```sh
bazel build //src/<module_name>:release_artifacts
```

where `:release_artifacts` is filegroup target that collects all release
artifacts of the module.

> NOTE: This is just proposal, the final decision is on module maintainer how
> the module code needs to be built.

### 3️⃣ Run Tests

```sh
bazel test //tests/...
```

---

## 🛠 Tools & Linters

The template integrates **tools and linters** from **centralized repositories** to ensure consistency across projects.

- **C++:** `clang-tidy`, `cppcheck`, `Google Test`
- **Rust:** `clippy`, `rustfmt`, `Rust Unit Tests`
- **CI/CD:** GitHub Actions for automated builds and tests

---

## 📖 Documentation

- A **centralized docs structure** is planned.

---

## ⚙️ `project_config.bzl`

This file defines project-specific metadata used by Bazel macros, such as `dash_license_checker`.

### 📌 Purpose

It provides structured configuration that helps determine behavior such as:

- Source language type (used to determine license check file format)
- Safety level or other compliance info (e.g. ASIL level)

### 📄 Example Content

```python
PROJECT_CONFIG = {
    "asil_level": "QM",  # or "ASIL-A", "ASIL-B", etc.
    "source_code": ["cpp", "rust"]  # Languages used in the module
}
```

### 🔧 Use Case

When used with macros like `dash_license_checker`, it allows dynamic selection of file types
 (e.g., `cargo`, `requirements`) based on the languages declared in `source_code`.
