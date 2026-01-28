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

/// Messages
///
/// This module contains the definition of messages
/// to be used within this application.
use postcard::experimental::max_size::MaxSize;
use serde::{Deserialize, Serialize};

/// Camera image
///
/// Fake camera frame with number of detected people and cars,
/// and distance to the closest obstacle.
#[derive(Serialize, Deserialize, MaxSize, Debug, Default)]
#[repr(C)]
pub struct CameraImage {
    pub num_people: usize,
    pub num_cars: usize,
    pub obstacle_distance: f64,
}

/// Scene
///
/// Fake inferred scene with number of detected people and cars,
/// distance to the closest obstacle, and distances to left and right lane.
#[derive(Serialize, Deserialize, MaxSize, Debug, Default)]
#[repr(C)]
pub struct Scene {
    pub num_people: usize,
    pub num_cars: usize,
    pub obstacle_distance: f64,
    pub distance_left_lane: f64,
    pub distance_right_lane: f64,
}
