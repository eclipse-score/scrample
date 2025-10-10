# Getting Started
This guide shows you how to setup and use the IPC module of score.

To use the communication module in your project you need to follow these steps:

## 1. Setup

### 1.1 Add this to your MODULE.bazel:

<details>
  <summary>MODULE.bazel</summary>

  ```starlark
  module(name = "use_com_test")

  bazel_dep(name = "score_toolchains_gcc", version = "0.4", dev_dependency=True)

  gcc = use_extension("@score_toolchains_gcc//extentions:gcc.bzl", "gcc", dev_dependency=True)
  gcc.toolchain(
      url = "https://github.com/eclipse-score/toolchains_gcc_packages/releases/download/0.0.1/x86_64-unknown-linux-gnu_gcc12.tar.gz",
      sha256 = "457f5f20f57528033cb840d708b507050d711ae93e009388847e113b11bf3600",
      strip_prefix = "x86_64-unknown-linux-gnu",
  )

  use_repo(gcc, "gcc_toolchain", "gcc_toolchain_gcc")

  bazel_dep(name = "rules_rust", version = "0.61.0")

  crate = use_extension("@rules_rust//crate_universe:extensions.bzl", "crate")

  crate.spec(package = "futures", version = "0.3.31")
  crate.spec(package = "libc", version = "0.2.155")
  crate.spec(package = "clap", version = "4.5.4", features = ["derive"])

  crate.from_specs(name = "crate_index")
  use_repo(crate, "crate_index")

  bazel_dep(name = "rules_boost", repo_name = "com_github_nelhage_rules_boost")
  archive_override(
      module_name = "rules_boost",
      urls = ["https://github.com/nelhage/rules_boost/archive/refs/heads/master.tar.gz"],
      strip_prefix = "rules_boost-master",
  )

  bazel_dep(name = "boost.program_options", version = "1.87.0")

  bazel_dep(name = "score-baselibs", version = "0.0.0")
  git_override(
      module_name = "score-baselibs",
      remote = "https://github.com/eclipse-score/baselibs.git",
      commit = "8f041f9cc1fa596585f8a8bb71bd90c06252e017",
  )

  # TRLC dependency for requirements traceability
  bazel_dep(name = "trlc", version = "0.0.0")
  git_override(
      module_name = "trlc",
      remote = "https://github.com/bmw-software-engineering/trlc.git",
      commit = "650b51a47264a4f232b3341f473527710fc32669",  # trlc-2.0.2 release
  )

  bazel_dep(name = "communication")
  git_override(
      module_name = "communication",
      remote = "https://github.com/eclipse-score/communication.git",
      commit = "aa3fa1f42a0fd7d0a92a462cf9b0f1e93c68618d",
  )

  ```

</details>


### 1.2 Insert this into your .bazelrc:
<details>
  <summary>.bazelrc</summary>
  
  ```
  common --@score-baselibs//score/mw/log/detail/flags:KUse_Stub_Implementation_Only=False
  common --@score-baselibs//score/mw/log/flags:KRemote_Logging=False
  common --@score-baselibs//score/json:base_library=nlohmann
  common --@communication//score/mw/com/flags:tracing_library=stub

  common --registry=https://raw.githubusercontent.com/eclipse-score/bazel_registry/refs/heads/main/
  common --registry=https://bcr.bazel.build
  ```
</details>


## 1.3 Run Bazel
If you start with a plain project add a empty file called `BUILD` into your project folder.

Now you can build the project with the command `bazel build //...` (so far nothing happens, because no targets were defined).

You can choose now to continue in this guide to create a simple consumer-producer program or start on your own.

## 2. Use it :)
Now that you have setup your project so far let`s start to send and receive some messages.

### 2.1 Basic Structure
First let´s create a folder `src`in our root project directory.

Inside `src` create following folders: `consumer`, `producer` and `etc`.

### 2.2 Message
Before we start sending messages we need to define what we will send.
Therefore create the file `message_data.h` in your `src` directory.
```cpp
#ifndef SCORE_MESSAGE_DATA_H
#define SCORE_MESSAGE_DATA_H

#include "score/mw/com/types.h"

namespace com_example
{

struct Message
{
  std::string message;
};

template <typename Trait>
class IPCInterface : public Trait::Base
{
  public:
    using Trait::Base::Base;

    typename Trait::template Event<Message> message_{*this, "message"};
};

using IPCInterfaceProxy = score::mw::com::AsProxy<IPCInterface>;
using IPCInterfaceSkeleton = score::mw::com::AsSkeleton<IPCInterface>;

} // namespace com_example

#endif //SCORE_MESSAGE_DATA_H
```

Let's take a deeper look into that.
We have the struct `Message` containing our `message` as a string.

Then we have a more difficult code snippet, where we define the `IPCInterface`. This interface is necessary so producer and consumer know what to send/receive.

After defining the interface we define:
- `IPCInterfaceProxy`: client-role
- `IPCInterfaceSkeleton`: server-role

You can have a deeper look into this architecture here: [Eclipse S-Core Communication Doc](https://eclipse-score.github.io/score/main/features/communication/docs/architecture/index.html#frontend)

### 2.3 Producer
The producer will (as its name suggests) produce/send the data.

Go inside the `producer` directory and create a new file called `producer.h`.

```cpp
#ifndef SCORE_PRODUCER_H
#define SCORE_PRODUCER_H

#include "score/mw/com/impl/instance_specifier.h"
#include "src/message_data.h"

class Producer
{
  public:
    Producer(const score::mw::com::impl::InstanceSpecifier& instance_specifier);
    ~Producer() = default;

    int RunProducer(const std::chrono::milliseconds cycle_time,
                    const std::size_t num_cycles);

  private:
    score::Result<com_example::IPCInterfaceSkeleton> create_result;
};

#endif //SCORE_PRODUCER_H
```

As you can see the header is lightweight, we will only need to use the Constructor and `RunProducer` from outside.
`create_result` is our `IPCInterfaceSkeleton` specified with the `instance_specifier` out of our `score_mw_com.json`.

After that create the file `producer.cpp`.

```cpp
#include "producer.h"
#include "src/message_data.h"

Producer::Producer(const score::mw::com::impl::InstanceSpecifier& instance_specifier)
  : create_result(com_example::IPCInterfaceSkeleton::Create(instance_specifier))
{
}

int Producer::RunProducer(const std::chrono::milliseconds cycle_time,
                      const std::size_t num_cycles)
{
  if (!create_result.has_value())
  {
    std::cerr << "Skeleton was not created. Can not run producer!\n"; 
    return EXIT_FAILURE;
  }
  auto& skeleton = create_result.value();

  const auto offer_result = skeleton.OfferService();

  if (!offer_result.has_value())
  {
    std::cerr << "Unable to offer service for skeleton!\n";
    return EXIT_FAILURE;
  }

  std::cout << "Starting to send data\n";

  for (std::size_t cycle = 0U; cycle < num_cycles || num_cycles == 0U; ++cycle)
  {
    auto cycle_message = "Message " + std::to_string(cycle);
    auto message = com_example::Message{.message=cycle_message};
    std::cout << "Sending: " << cycle_message << std::endl;
    skeleton.message_.Send(std::move(message));

    std::this_thread::sleep_for(cycle_time);
  }

  skeleton.StopOfferService();

  return EXIT_SUCCESS;
}
```

Here we have a bit more code.

Let's start with the constructor, which only initializes `create_result`.

More complex is `RunProducer`, where we first check if the initialization of `create_result` was successful.
Then we offer our service and also check if it was successful.
If so, we start to send our messages in a loop.

At the end we need to stop offering the service.

### 2.4 Consumer
On the other side the consumer will consume/receive the data.

Go inside the `consumer` directory and create a new file called `consumer.h`