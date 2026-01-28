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

use crate::activities::common::*;

use crate::activities::camera_activity::Camera;
use crate::activities::mcap_activity::Mcap;
use crate::activities::render_activity::SceneRender;

use core::net::{IpAddr, Ipv4Addr, SocketAddr};
use feo::activity::{ActivityBuilder, ActivityIdAndBuilder};
use feo::ids::{ActivityId, AgentId, WorkerId};
use feo::topicspec::{Direction, TopicSpecification};
use feo_com::interface::ComBackend;
use std::collections::HashMap;
use std::path::{Path, PathBuf};

pub type WorkerAssignment = (WorkerId, Vec<(ActivityId, Box<dyn ActivityBuilder>)>);

pub type ActivityDependencies = HashMap<ActivityId, Vec<ActivityId>>;

pub const COM_BACKEND: ComBackend = ComBackend::Iox2;

pub const BIND_ADDR: SocketAddr = SocketAddr::new(IpAddr::V4(Ipv4Addr::LOCALHOST), 8081);
pub const BIND_ADDR2: SocketAddr = SocketAddr::new(IpAddr::V4(Ipv4Addr::LOCALHOST), 8082);

pub const TOPIC_INFERRED_SCENE: &str = "feo/com/vehicle/inferred/scene";
pub const TOPIC_CAMERA_FRONT: &str = "feo/com/vehicle/camera/front";

/// Just for demonstration purposes, currently we dont use recorders.
pub const MAX_ADDITIONAL_SUBSCRIBERS: usize = 2;

pub fn socket_paths() -> (PathBuf, PathBuf) {
    (
        Path::new("/tmp/feo_listener1.socket").to_owned(),
        Path::new("/tmp/feo_listener2.socket").to_owned(),
    )
}

pub fn agent_assignments() -> HashMap<AgentId, Vec<(WorkerId, Vec<ActivityIdAndBuilder>)>> {
    let worker_40: WorkerAssignment = (
        40.into(),
        vec![
            (
                0.into(),
                Box::new(|id| Camera::build(id, TOPIC_CAMERA_FRONT)),
            ),
            (
                1.into(),
                Box::new(|id| SceneRender::build(id, TOPIC_CAMERA_FRONT, TOPIC_INFERRED_SCENE)),
            ),
        ],
    );
    let worker_41: WorkerAssignment = (41.into(), vec![(2.into(), Box::new(|id| Mcap::build(id)))]);

    let assignment = [(100.into(), vec![worker_40]), (101.into(), vec![worker_41])]
        .into_iter()
        .collect();

    assignment
}

pub fn activity_dependencies() -> ActivityDependencies {
    let dependencies = [
        (0.into(), vec![]),
        (1.into(), vec![0.into()]),
        (2.into(), vec![]),
    ];

    dependencies.into()
}

pub fn topic_dependencies<'a>() -> Vec<TopicSpecification<'a>> {
    use Direction::*;

    vec![
        TopicSpecification::new::<CameraImage>(
            TOPIC_CAMERA_FRONT,
            vec![(0.into(), Outgoing), (1.into(), Incoming)],
        ),
        TopicSpecification::new::<Scene>(TOPIC_INFERRED_SCENE, vec![(1.into(), Outgoing)]),
    ]
}

pub fn worker_agent_map() -> HashMap<WorkerId, AgentId> {
    agent_assignments()
        .iter()
        .flat_map(|(agent_id, worker_activity_map)| {
            worker_activity_map
                .iter()
                .map(move |(worker_id, _)| (*worker_id, *agent_id))
        })
        .collect()
}

pub fn agent_assignments_ids() -> HashMap<AgentId, Vec<(WorkerId, Vec<ActivityId>)>> {
    agent_assignments()
        .into_iter()
        .map(|(agent_id, worker_activity_map)| {
            (
                agent_id,
                worker_activity_map
                    .into_iter()
                    .map(|(worker_id, activity_and_builder)| {
                        (
                            worker_id,
                            activity_and_builder
                                .into_iter()
                                .map(|(activity_id, _)| activity_id)
                                .collect(),
                        )
                    })
                    .collect(),
            )
        })
        .collect()
}
