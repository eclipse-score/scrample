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

use ad_demo::activities::application_config::{
    agent_assignments, agent_assignments_ids, topic_dependencies, BIND_ADDR, BIND_ADDR2,
    COM_BACKEND,
};
use core::time::Duration;
use feo::agent::com_init::initialize_com_secondary;
use feo::agent::relayed::secondary::{Secondary, SecondaryConfig};
use feo::agent::NodeAddress;
use feo::ids::ActivityId;
use feo::ids::AgentId;
use feo_log::{info, LevelFilter};
use std::collections::HashSet;

/// One Secondary agent for the ad-demo FEO application.
fn main() {
    feo_logger::init(LevelFilter::Debug, true, true);
    feo_tracing::init(feo_tracing::LevelFilter::TRACE);

    let secondary_agent_id = AgentId::new(101);

    let config = SecondaryConfig {
        id: secondary_agent_id,
        worker_assignments: agent_assignments().remove(&secondary_agent_id).unwrap(),
        timeout: Duration::from_secs(10),
        bind_address_senders: NodeAddress::Tcp(BIND_ADDR),
        bind_address_receivers: NodeAddress::Tcp(BIND_ADDR2),
    };

    let local_activities: HashSet<ActivityId> = agent_assignments_ids()
        .remove(&secondary_agent_id)
        .unwrap()
        .iter()
        .flat_map(|(_, activity_ids)| activity_ids.iter())
        .copied()
        .collect();

    // Initialize topics. Make it alive until application runs.
    let _topic_guards =
        initialize_com_secondary(COM_BACKEND, topic_dependencies(), &local_activities);

    info!("Starting secondary agent {}", secondary_agent_id);

    let secondary_agent = Secondary::new(config);
    secondary_agent.run();
}
