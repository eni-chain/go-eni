#!/bin/bash

mkdir $HOME/eni-snapshot
mkdir $HOME/key_backup
# Move priv_validator_state out so it isn't used by anyone else
mv $HOME/.eni/data/priv_validator_state.json $HOME/key_backup
# Create backups
cd $HOME/eni-snapshot
tar -czf data.tar.gz -C $HOME/.eni data/
tar -czf wasm.tar.gz -C $HOME/.eni wasm/
echo "Data and Wasm snapshots created in $HOME/eni-snapshot"