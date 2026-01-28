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

/// Camera activity
///
/// This activity emulates a camera generating a [CameraImage].
#[derive(Debug)]
pub struct Camera {
    activity_id: ActivityId,
    output_image: Box<dyn ActivityOutput<CameraImage>>,
}

impl Camera {
    pub fn build(activity_id: ActivityId, image_topic: &str) -> Box<dyn Activity> {
        Box::new(Self {
            activity_id,
            output_image: activity_output(image_topic),
        })
    }

    fn get_image(&mut self) -> CameraImage {
        let mut rnd_generator = rand::rng();

        CameraImage {
            num_people: rnd_generator.random_range(0..20),
            num_cars: rnd_generator.random_range(0..7),
            obstacle_distance: rnd_generator.random_range(20.0..50.0),
        }
    }
}

impl Activity for Camera {
    fn id(&self) -> ActivityId {
        self.activity_id
    }

    #[instrument(name = "Camera startup")]
    fn startup(&mut self) {}

    #[instrument(name = "Camera")]
    fn step(&mut self) {
        debug!("Stepping Camera");

        if let Ok(camera) = self.output_image.write_uninit() {
            let image = self.get_image();
            debug!("Sending image: {image:?}");
            let camera = camera.write_payload(image);
            camera.send().unwrap();
        }
    }

    #[instrument(name = "Camera shutdown")]
    fn shutdown(&mut self) {}
}
