### Steps to run the foxglove websocket server script
1. **For first run only**:
    1. `bazelisk run //feo/ad-demo/lichtblick-com:requirements.update`
        1. It will use [requirements file](requirements.in) to create [requirements lock file](requirements_lock.txt) to install required packages in the pip environment.
        1. you can always test if the requirements are up to date for your copy
            1. `bazelisk test //feo/ad-demo/lichtblick-com:requirements_test`
1. **In one terminal**: run the FoxGloveServer & TCP server on localhost
    1. `bazelisk run //feo/ad-demo/lichtblick-com:foxglove_ws_server`
        1. By default this runs the web socket server and tcp server on localhost
        1. To run the server on some other machine or remote host, use the ip of the remote machine:
            1. `bazelisk run //feo/ad-demo/lichtblick-com:foxglove_ws_server <ip.of.remote.machine>`
        1. This will run the python binary target (the servers) on a bazel environment which isolates it from the sytem - similar to python virtual environment.
1. **In other terminals**: [Start the primary and secondary agents in other terminals to run the FEO application](../README.md#running).
    1. This message will then be forwarded via the foxglove web socket server to lichtblick with topic name **/gps/fix** in the **messaging schema format [locationfix](https://lichtblick-suite.github.io/docs/docs/visualization/message-schemas/location-fix)**.
