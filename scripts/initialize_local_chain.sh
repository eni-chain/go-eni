#!/bin/bash
# require success for commands
set -e

cd ../
make build
make reset-eni-node
./build/enid start --home=./eni-node