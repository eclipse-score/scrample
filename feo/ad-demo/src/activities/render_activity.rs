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
use crate::activities::messages::{CameraImage, Scene};

/// Scene render activity
///
/// This activity emulates rendering a scene from a [CameraImage] which can be rendered to a display.
#[derive(Debug)]
pub struct SceneRender {
    activity_id: ActivityId,
    input_image: Box<dyn ActivityInput<CameraImage>>,
    output_scene: Box<dyn ActivityOutput<Scene>>,
}

impl SceneRender {
    pub fn build(
        activity_id: ActivityId,
        image_topic: &str,
        scene_topic: &str,
    ) -> Box<dyn Activity> {
        Box::new(Self {
            activity_id,
            input_image: activity_input(image_topic),
            output_scene: activity_output(scene_topic),
        })
    }
}

impl Activity for SceneRender {
    fn id(&self) -> ActivityId {
        self.activity_id
    }

    #[instrument(name = "SceneRender startup")]
    fn startup(&mut self) {}

    #[instrument(name = "SceneRender")]
    fn step(&mut self) {
        debug!("Stepping SceneRender");

        let received_image = self.input_image.read();
        let output_scene = self.output_scene.write_uninit();

        if let (Ok(received_image), Ok(output_scene)) = (received_image, output_scene) {
            debug!("Received image: {:?}", *received_image);

            let mut rnd_generator = rand::rng();

            let op_scene = Scene {
                num_people: received_image.num_people,
                num_cars: received_image.num_cars,
                obstacle_distance: received_image.obstacle_distance,
                distance_left_lane: rnd_generator.random_range(10.0..50.0),
                distance_right_lane: rnd_generator.random_range(10.0..50.0),
            };

            debug!("Sending scene: {op_scene:?}");

            let send_output_scene = output_scene.write_payload(op_scene);

            send_output_scene.send().unwrap();

            debug!("Rendering scene");
        }
    }

    #[instrument(name = "SceneRender shutdown")]
    fn shutdown(&mut self) {}
}
