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

use anyhow::{Context, Result};
use camino::Utf8Path;
use mcap::MessageStream;
use memmap2::Mmap;
use serde_json;
use serde_json::Value;
use std::fs;

/// Mcap activity
///
/// This activity reads mcap msgs in Mcap file and publish it to tcp server,
/// to be forwarded to lichtblick via foxglove websocket server.
pub struct Mcap {
    activity_id: ActivityId,
    message_stream: MessageStream<'static>,
}

impl fmt::Debug for Mcap {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.debug_struct("Mcap")
            .field("activity_id", &self.activity_id)
            .field("message_stream", &"<MessageStream>")
            .finish()
    }
}

impl Mcap {
    fn map_mcap<P: AsRef<Utf8Path>>(p: P) -> Result<Mmap> {
        let path = p.as_ref();
        debug!("Opening Mcap file at path: {}", path);
        let file = fs::File::open(path).context("Couldn't open MCAP file at")?;
        unsafe { Mmap::map(&file) }.context("Couldn't map MCAP file")
    }

    pub fn build(activity_id: ActivityId) -> Box<dyn Activity> {
        let mcap_to_map = Self::map_mcap("feo/ad-demo/src/assets/gps_route.mcap").expect("Could not open MCAP file");
        let static_slice: &'static [u8] = Box::leak(mcap_to_map.to_vec().into_boxed_slice());
        let message_stream = MessageStream::new(static_slice).expect("Failed to create MessageStream");
        Box::new(Self {
            activity_id,
            message_stream,
        })
    }

    fn get_single_msg(&mut self) -> Result<Option<Value>> {
        match self.message_stream.next() {
            Some(Ok(message)) => {
                let data_json = serde_json::from_slice(&message.data).context("Failed to convert msg data as JSON")?;

                debug!("single mcap message data: {}", data_json);

                Ok(Some(data_json))
            },
            Some(Err(e)) => {
                debug!("Error reading MCAP message: {}", e);
                Ok(None)
            },
            None => {
                debug!("No more messages in MCAP file");
                Ok(None)
            },
        }
    }

    fn send_tcp_msg(&mut self, msg: &str) {
        use std::io::Write;
        use std::net::{Shutdown, TcpStream};
        if let Ok(mut stream) = TcpStream::connect("127.0.0.1:9001") {
            debug!("Connected to the tcp server!");

            let _ = stream.write(msg.as_bytes());

            debug!("Message sent to the TCP server!");

            let _ = stream.shutdown(Shutdown::Both);
        } else {
            debug!("Couldn't connect to tcp server...");
        }
    }
}

impl Activity for Mcap {
    fn id(&self) -> ActivityId {
        self.activity_id
    }

    #[instrument(name = "Mcap startup")]
    fn startup(&mut self) {
        debug!("Mcap Startup");
    }

    #[instrument(name = "Mcap")]
    fn step(&mut self) {
        debug!("Stepping Mcap");

        if let Ok(Some(mcap_msg_json)) = self.get_single_msg() {
            let compact_json = serde_json::to_string(&mcap_msg_json).expect("failed to stringify Json");

            debug!("Read Mcap message: {compact_json:?}");
            self.send_tcp_msg(&compact_json);
        } else {
            debug!("No message in mcap left.");
            self.send_tcp_msg("No message in mcap left");
        }
    }

    #[instrument(name = "Mcap shutdown")]
    fn shutdown(&mut self) {
        debug!("Mcap Shutdown");
    }
}
