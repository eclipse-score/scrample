# *******************************************************************************
# Copyright (c) 2025 Contributors to the Eclipse Foundation
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
cc_binary(
    name = "scrample",
    srcs = [
        "assert_handler.cpp",
        "assert_handler.h",
        "main.cpp",
    ],
    visibility = ["//visibility:public"],
    deps = [
        ":sample_sender_receiver",
        "@boost.program_options",
        "@score_baselibs//score/language/futurecpp",
        "@score_logging//score/mw/log",
        "@score_communication//score/mw/com",
    ],
)

cc_library(
    name = "sample_sender_receiver",
    srcs = [
        "sample_sender_receiver.cpp",
    ],
    hdrs = [
        "sample_sender_receiver.h",
    ],
    deps = [
        ":datatype",
        "@score_baselibs//score/mw/log:frontend",
        "@score_communication//score/mw/com",
    ],
)

cc_library(
    name = "datatype",
    srcs = [
        "datatype.cpp",
    ],
    hdrs = [
        "datatype.h",
    ],
    deps = [
        "@score_baselibs//score/language/futurecpp",
        "@score_communication//score/mw/com",
    ],
)
