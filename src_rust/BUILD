# *******************************************************************************
# Copyright (c) 2026 Contributors to the Eclipse Foundation
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
load("@rules_rust//rust:defs.bzl", "rust_binary")

RUSTC_FLAGS = select({
    "@platforms//os:qnx": [
        "-Clink-arg=-lc++",
        "-Clink-arg=-lm",
    ],
    "//conditions:default": [
        "-Clink-arg=-lstdc++",
        "-Clink-arg=-lm",
        "-Clink-arg=-lc",
    ],
})

rust_binary(
    name = "scrample_rust",
    srcs = ["main.rs"],
    data = ["//config:logging.json"],
    edition = "2021",
    rustc_flags = RUSTC_FLAGS,
    visibility = ["//visibility:public"],
    deps = [
        "@score_baselibs_rust//src/log/score_log",
        "@score_logging//score/mw/log/rust/score_log_bridge",
    ],
)
