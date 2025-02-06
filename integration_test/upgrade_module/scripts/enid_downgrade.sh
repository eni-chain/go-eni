#!/bin/bash

# This downgrades the binary to the currently-set UPGRADE_VERSION_LIST
# UPGRADE_VERSION_LIST is an ENV var that is the default version for upgrade tests

NODE_ID=${ID:-0}
INVARIANT_CHECK_INTERVAL=${INVARIANT_CHECK_INTERVAL:-0}
LOG_DIR="build/generated/logs"

# kill the existing service
pkill -f "enid start"

# start the service with a different UPGRADE_VERSION_LIST
UPGRADE_VERSION_LIST=$UPGRADE_VERSION_LIST enid start --chain-id eni --inv-check-period ${INVARIANT_CHECK_INTERVAL} > "$LOG_DIR/enid-$NODE_ID.log" 2>&1 &

echo "PASS"
exit 0
