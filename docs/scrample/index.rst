..
   # *******************************************************************************
   # Copyright (c) 2024 Contributors to the Eclipse Foundation
   #
   # See the NOTICE file(s) distributed with this work for additional
   # information regarding copyright ownership.
   #
   # This program and the accompanying materials are made available under the
   # terms of the Apache License Version 2.0 which is available at
   # https://www.apache.org/licenses/LICENSE-2.0
   #
   # SPDX-License-Identifier: Apache-2.0
   # *******************************************************************************

Scrample example
=============================

Welcome to the SCRAMPLE application!

Overview
--------

**SCRAMPLE** (S-CORE + Sample) is a minimal demonstration application for the
`Eclipse S-CORE <https://projects.eclipse.org/projects/automotive.score>`_ platform.
It provides two equivalent "Hello World" apps — one written in **C++** and one in **Rust** —
that both use the S-CORE ``mw::log`` logging library to emit a single log message to the console.

The goal of SCRAMPLE is to show how to set up a minimal Bazel workspace that integrates
S-CORE dependencies and produces a working logging application in either language.

Applications
------------

C++ application (``//src_cpp:scrample_cpp``)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The C++ app uses the ``score::mw::log`` API from ``@score_logging//score/mw/log``.
It locates the shared logging configuration at runtime via the ``MW_LOG_CONFIG_FILE``
environment variable, resolved from Bazel runfiles:

.. code-block:: cpp

   #include "score/mw/log/logging.h"

   int main()
   {
       // logging configuration is loaded from config/logging.json via MW_LOG_CONFIG_FILE
       score::mw::log::LogInfo("scrample") << "Hello from SCRAMPLE!";
       return 0;
   }

Rust application (``//src_rust:scrample_rust``)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The Rust app uses the ``score_log`` crate together with the ``score_log_bridge`` crate,
which bridges the Rust logging facade to the underlying ``mw::log`` C++ backend.
The path to the shared logging configuration is resolved from Bazel runfiles at startup:

.. code-block:: rust

   use score_log::info;
   use score_log_bridge::ScoreLogBridgeBuilder;

   fn main() {
       // logging configuration is loaded from config/logging.json via runfiles
       ScoreLogBridgeBuilder::new()
           .context("scrample")
           .config(config_path)
           .set_as_default_logger();
       info!("Hello from SCRAMPLE!");
   }

Shared Logging Configuration
-----------------------------

Both apps share the same ``config/logging.json`` file, which configures ``mw::log``
to write to the console at ``kInfo`` level or above:

.. code-block:: json

   {
       "appId": "SCRM",
       "appDesc": "SCRAMPLE App",
       "logMode": "kConsole",
       "logLevel": "kVerbose",
       "logLevelThresholdConsole": "kInfo"
   }

Building
--------

Build both apps for the host platform:

.. code-block:: bash

   bazel build --config=host //src_cpp:scrample_cpp //src_rust:scrample_rust

Cross-compile for QNX:

.. code-block:: bash

   bazel build --config=x86_64-qnx //src_cpp:scrample_cpp
   bazel build --config=x86_64-qnx //src_rust:scrample_rust

Running
-------

After a successful host build, run either app directly:

.. code-block:: bash

   ./bazel-bin/src_cpp/scrample_cpp
   ./bazel-bin/src_rust/scrample_rust

Both produce output in the following format:

.. code-block:: text

   2026/05/22 12:00:00.0000000 00000000 000 ECU1 SCRM scra log info verbose 1 Hello from SCRAMPLE!