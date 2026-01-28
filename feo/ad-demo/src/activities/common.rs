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

/// common imports for activities and config
pub(in crate::activities) use crate::activities::messages::{CameraImage, Scene};

pub(in crate::activities) use core::fmt;
pub(in crate::activities) use rand::Rng;

pub(in crate::activities) use feo::activity::Activity;
pub(in crate::activities) use feo::ids::ActivityId;
pub(in crate::activities) use feo_com::interface::{ActivityInput, ActivityOutput};
pub(in crate::activities) use feo_log::debug;
pub(in crate::activities) use feo_tracing::instrument;

// imports specific to this file
use feo_com::iox2::{Iox2Input, Iox2Output};

/// Create an activity input.
pub(in crate::activities) fn activity_input<T>(topic: &str) -> Box<dyn ActivityInput<T>>
where
    T: fmt::Debug + 'static,
{
    return Box::new(Iox2Input::new(topic));
}

/// Create an activity output.
pub(in crate::activities) fn activity_output<T>(topic: &str) -> Box<dyn ActivityOutput<T>>
where
    T: fmt::Debug + 'static,
{
    return Box::new(Iox2Output::new(topic));
}
