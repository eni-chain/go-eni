#!/bin/bash
# require success for commands
set -e

# Change to the script's parent directory (project root)
cd "$(dirname "$0")/.."

make build
make reset-eni-node
nohup ./build/enid start --home=./eni-node &