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
import asyncio
import json
import time
import argparse
from foxglove_websocket import run_cancellable
from foxglove_websocket.server import FoxgloveServer, FoxgloveServerListener
from foxglove_websocket.types import ChannelId

LOCALHOST = "127.0.0.1"
TCP_PORT = 9001

"""
TCP listener that receives Json messages and forwards them to Foxglove channel
"""


async def tcp_listener(foxglove_server: FoxgloveServer, channel_id: ChannelId):
    async def handle_client(reader, writer):
        address = writer.get_extra_info("peername")
        print(f"TCP connection opened at : {address}")

        try:
            while True:
                data = await reader.readline()
                if not data:
                    break

                tcp_message = data.decode("utf-8").strip()
                if tcp_message:
                    message_fxg = json.loads(tcp_message)
                    # covariance to just make location more visible in Lichtblick
                    message_fxg.update(
                        {
                            "position_covariance": [
                                700,
                                0.0,
                                0.0,
                                0.0,
                                700,
                                0.0,
                                0.0,
                                0.0,
                                700,
                            ],
                            "position_covariance_type": 2,
                        },
                    )

                    print(f"Received TCP message: {message_fxg}")

                    await foxglove_server.send_message(
                        channel_id,
                        time.time_ns(),
                        json.dumps(message_fxg).encode("utf8"),
                    )
        except Exception as e:
            print(f"Error in TCP client: {e}")
        finally:
            writer.close()
            await writer.wait_closed()
            print(f"TCP connection closed at : {address}")

    tcp_server = await asyncio.start_server(handle_client, LOCALHOST, TCP_PORT)

    tcp_address = tcp_server.sockets[0].getsockname()
    print(f"TCP server listening on {tcp_address}")

    async with tcp_server:
        await tcp_server.serve_forever()


"""
Main function to start Foxglove server and TCP listener
"""


async def main(host):
    class Listener(FoxgloveServerListener):
        async def on_subscribe(self, server: FoxgloveServer, channel_id: ChannelId):
            print("First client subscribed to", channel_id)

        async def on_unsubscribe(self, server: FoxgloveServer, channel_id: ChannelId):
            print("Last client unsubscribed from", channel_id)

    async with FoxgloveServer(
        host,
        8765,
        "foxglove server",
        supported_encodings=["json"],
    ) as server:
        server.set_listener(Listener())
        channel_id = await server.add_channel(
            {
                "topic": "/gps/fix",
                "encoding": "json",
                "schemaName": "foxglove.LocationFix",
                "schema": json.dumps(
                    {
                        "type": "object",
                        "properties": {
                            "altitude": {"type": "number"},
                            "latitude": {"type": "number"},
                            "longitude": {"type": "number"},
                            "position_covariance": {
                                "type": "array",
                                "items": {"type": "number"},
                                "minItems": 9,
                                "maxItems": 9,
                            },
                            "position_covariance_type": {"type": "integer"},
                        },
                        "required": ["latitude", "longitude"],
                    }
                ),
                "schemaEncoding": "jsonschema",
            }
        )

        tcp_task = asyncio.create_task(tcp_listener(server, channel_id))

        try:
            await asyncio.Event().wait()
        finally:
            tcp_task.cancel()
            try:
                await tcp_task
            except asyncio.CancelledError:
                pass
            print("TCP server closed.")


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "host",
        type=str,
        nargs="?",
        default=LOCALHOST,
        help="Host the server in this machine. If not provided, defaults to localhost",
    )
    args = parser.parse_args()
    run_cancellable(main(args.host))
