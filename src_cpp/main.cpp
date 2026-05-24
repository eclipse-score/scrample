/********************************************************************************
 * Copyright (c) 2026 Contributors to the Eclipse Foundation
 *
 * See the NOTICE file(s) distributed with this work for additional
 * information regarding copyright ownership.
 *
 * This program and the accompanying materials are made available under the
 * terms of the Apache License Version 2.0 which is available at
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * SPDX-License-Identifier: Apache-2.0
 ********************************************************************************/

#include "score/mw/log/logging.h"

#include <climits>
#include <cstdlib>
#include <string>
#include <unistd.h>

int main()
{
    const char* runfiles_dir = std::getenv("RUNFILES_DIR");
    if (runfiles_dir != nullptr)
    {
        std::string config_path = std::string(runfiles_dir) + "/_main/config/logging.json";
        setenv("MW_LOG_CONFIG_FILE", config_path.c_str(), 1);
    }
    else
    {
        char exe_path[PATH_MAX];
        ssize_t len = readlink("/proc/self/exe", exe_path, sizeof(exe_path) - 1);
        if (len != -1)
        {
            exe_path[len] = '\0';
            std::string config_path = std::string(exe_path) + ".runfiles/_main/config/logging.json";
            setenv("MW_LOG_CONFIG_FILE", config_path.c_str(), 1);
        }
    }

    score::mw::log::LogInfo("scrample") << "Hello from SCRAMPLE!";
    return 0;
}
