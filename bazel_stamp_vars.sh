#!/bin/bash
# Stable status values - these cause a rebuild when changed
cat <<EOF
STABLE_VERSION ${VERSION:-dev}
STABLE_GIT_COMMIT $(git rev-parse HEAD 2>/dev/null || echo "unknown")
STABLE_BUILD_DATE $(date -u +"%Y-%m-%dT%H:%M:%SZ")
EOF
