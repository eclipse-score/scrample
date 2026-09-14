// *******************************************************************************
// Copyright (c) 2026 Contributors to the Eclipse Foundation
//
// See the NOTICE file(s) distributed with this work for additional
// information regarding copyright ownership.
//
// This program and the accompanying materials are made available under the
// terms of the Apache License Version 2.0 which is available at
// <https://www.apache.org/licenses/LICENSE-2.0>
//
// SPDX-License-Identifier: Apache-2.0
// *******************************************************************************
use score_log::info;
use score_log_bridge::ScoreLogBridgeBuilder;
use std::path::PathBuf;

fn main() {
    let config_path = std::env::var("RUNFILES_DIR")
        .map(|d| PathBuf::from(d).join("_main/config/logging.json"))
        .unwrap_or_else(|_| {
            let exe = std::env::current_exe().unwrap();
            PathBuf::from(format!("{}.runfiles", exe.display())).join("_main/config/logging.json")
        });

    ScoreLogBridgeBuilder::new()
        .context("scrample")
        .config(config_path)
        .set_as_default_logger();

    info!("Hello from SCRAMPLE!");
}
