# AD (Autonomous Driving) Demo using FEO Framework

### Activities overview in ad-demo FEO (Fixed Execution Order) application
This [FEO](https://eclipse-score.github.io/score/main/features/frameworks/feo/index.html#) application contains three activities to demonstrate activity and topic creation, define their dependencies and how they communicate:
1. **Camera** - Simulates a camera sensor generating person detection data with fake object count, car count, obstacle distance and outputs as **CameraImage** topic message.
2. **SceneRender** - This activity simulates processing of camera images and extracting fake lane distance information to create enhanced scene data on top of Camera data. This activity commuincates with Camera activity to receive CameraImage topic and outputs **Scene** topic message.
3. **Mcap** - This activity reads GPS route data from an MCAP file and publishes to a TCP server for visualization in Lichtblick via Foxglove WebSocket Server. This activity is standalone and neither does commuincate with any other activity nor outputs any message.
This application tries to demostrate the basic usage of FEO framework.

### First time setup
1. Follow the steps mentioned here: https://code.visualstudio.com/docs/languages/rust
    1. It gives all the required steps to get started with rust on linux
    1. Install [rustup](https://www.rust-lang.org/tools/install) in your system
        1. These are the set of available [components](https://rust-lang.github.io/rustup/concepts/components.html) once you install rustup in your system:
    1. Install [rust analyzer extension](https://marketplace.visualstudio.com/items?itemName=rust-lang.rust-analyzer) for Rust lint and debug and support
1. So that rust analyzer works with the rust intelligence since you have a bazel project with minimal cargo.toml. You need to have a `rust_project.json` so that your rust code intelligence works properly. Run the following, from project root, to generate it:
    1. `bazelisk run @rules_rust//tools/rust_analyzer:gen_rust_project -- //feo/ad-demo/...`
    1. _Tip_ : Run it every time you would update the project or reorganize it.

### Running ad-demo application
1. Run **agent_primary** (one terminal) :
    1. `bazelisk run //feo/ad-demo:agent_primary 2000`
        1. The FEO application cycles every 2000 milliseconds.
        1. This will run the Camera and Environment Renderer activity
2. Run **agent_secondary** (another terminal) :
    1. `bazelisk run //feo/ad-demo:agent_secondary`
    1. This will run Mcap activity
    1. Now, you can see in the logs the interaction between Camera and SceneRender activity, whereas MCap activity is independent to any of them.
        1. Camera and SceneRender logs are visible in previous terminal whereas MCap logs are visible in this terminal.
3. Run [foxglove websocket server script (in yet another terminal)](lichtblick-com/README.md).
    1. `bazelisk run //feo/ad-demo/lichtblick-com:foxglove_ws_server`
        1. Pleae check out the [Readme](lichtblick-com/README.md) for your first run.
1. Now open [Lichtblick](https://github.com/Lichtblick-Suite/lichtblick/releases) and
    1. Go to "Open connection" → "Foxglove WebSocket"
    1. Enter: `ws://localhost:8765`
        1. usually this is default and already exists
        1. if hosted on remote machine Enter: `ws://<ip.of.remote.machine>:8765`
    1. Now open [map panel](https://lichtblick-suite.github.io/docs/docs/visualization/message-schemas/location-fix) and listen to topic name **/gps/fix**.
        1. you can also open [raw message panel](https://lichtblick-suite.github.io/docs/docs/visualization/panels/raw-messages-panel) and listen to topic name **/gps/fix** - to see the lat long.
    1. **You should be able to see that the location (latitude, longitude) message from Mcap activity in FEO application is visible in map panel coming via foxglove webserver.**
