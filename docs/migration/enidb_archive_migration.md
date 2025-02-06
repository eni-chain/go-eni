# EniDB Archive Migration Guide

## Overview
EniDB is the next generation of chain storage in EniV2.
One issue for running EniDB on archive nodes is that we need to keep the full state of the chain, so we can't
state sync it and clear out previous iavl data.

In order to run an archive node with EniDB, we need to run a migration from iavl to eni-db.

The overall process will work as follows:

1. Update config to enable EniDB (state committment + state store)
2. Stop the node and Run SC Migration
3. Note down MIGRATION_HEIGHT
4. Re start enid with `--migrate-iavl` enabled (migrating state store in background)
5. Verify migration at various sampled heights once state store is complete
6. Restart enid normally and verify node runs properly
7. Clear out iavl and restart enid normally, now only using EniDB fully


## Requirements

### Additional Disk Space

Ensure you have at least 10TB of free disk space available for the migration process. This extra space is needed because during migration, both IAVL and EniDB state stores will exist simultaneously. Note
that this is ONLY during the migration process, and will not be necessary after completion


## Migration Steps

### Step 1: Add EniDB Configurations
We can enable EniDB by adding the following configs to app.toml file.
Usually you can find this file under ~/.eni/config/app.toml.
```bash
#############################################################################
###                             EniDB Configuration                       ###
#############################################################################

[state-commit]
# Enable defines if the EniDB should be enabled to override existing IAVL db backend.
sc-enable = true

# AsyncCommitBuffer defines the size of asynchronous commit queue, this greatly improve block catching-up
# performance, <=0 means synchronous commit.
sc-async-commit-buffer = 100

# SnapshotKeepRecent defines how many memiavl snapshots (beyond the latest one) to keep
# Recommend to set to 1 to make sure IBC relayers work.
sc-keep-recent = 1

# SnapshotInterval defines the number of blocks interval the memiavl snapshot is taken, default to 10000 blocks.
# Adjust based on your needs:
# Setting this too low could lead to lot of extra heavy disk IO
# Setting this too high could lead to slow restart
sc-snapshot-interval = 10000

# SnapshotWriterLimit defines the max concurrency for taking memiavl snapshot
sc-snapshot-writer-limit = 1

# CacheSize defines the size of the LRU cache for each store on top of the tree, default to 100000.
sc-cache-size = 100000

[state-store]
# Enable defines if the state-store should be enabled for historical queries.
# In order to use state-store, you need to make sure to enable state-commit at the same time.
# Validator nodes should turn this off.
# State sync node or full nodes should turn this on.
ss-enable = true

# DBBackend defines the backend database used for state-store.
# Supported backends: pebbledb, rocksdb
# defaults to pebbledb (recommended)
ss-backend = "pebbledb"

# AsyncWriteBuffer defines the async queue length for commits to be applied to State Store
# Set <= 0 for synchronous writes, which means commits also need to wait for data to be persisted in State Store.
# defaults to 100
ss-async-write-buffer = 100

# KeepRecent defines the number of versions to keep in state store
# Setting it to 0 means keep everything, default to 100000
ss-keep-recent = 0

# PruneIntervalSeconds defines the minimum interval in seconds + some random delay to trigger pruning.
# It is more efficient to trigger pruning less frequently with large interval.
# default to 600 seconds
ss-prune-interval = 600

# ImportNumWorkers defines the concurrency for state sync import
# defaults to 1
ss-import-num-workers = 1
```


### Step 2: Stop the node and Run SC Migration

```bash
systemctl stop enid
enid tools migrate-iavl --home-dir /root/.eni
```

You can also run this sc migration in the background:
```bash
enid tools migrate-iavl --home-dir /root/.eni > migrate-sc.log 2>&1 &
```

This may take a couple hours to run. You will see logs of form
`Start restoring SC store for height`


### Step 3: Note down MIGRATION_HEIGHT
Note down the latest height as outputted from the sc migration log. 

As an example the sc migration log would show:
```
latest version: 111417590
D[2024-10-29|15:26:03.657] Finished loading IAVL tree
D[2024-10-29|15:26:03.657] Finished loading IAVL tree
D[2024-10-29|15:26:03.657] Finished loading IAVL tree
```

Save that `latest version` (in example 111417590) as an env var $MIGRATION_HEIGHT.

```bash
MIGRATION_HEIGHT=<>
```


### Step 4: Restart enid with background SS migration

If you are using systemd, make sure to update your service configuration to use this command.
Always be sure to run with these flags until migration is complete.
```bash
enid start --migrate-iavl --migrate-height $MIGRATION_HEIGHT --chain-id pacific-1
```

Enid will run normally and the migration will run in the background. Data from iavl
will be written to SS and new writes will be directed at SS not iavl.

You will see logs of form 
`EniDB Archive Migration: Iterating through %s module...` and 
`EniDB Archive Migration: Last 1,000,000 iterations took:...`

NOTE: While this is running, any historical queries will be routed to iavl if
they are for a height BEFORE the migrate-height. Any queries on heights
AFTER the migrate-height will be routed to state store (pebbledb).


### Step 5: Verify State Store Migration after completion
Once State Store Migration is complete, you will see logs of form
`EniDB Archive Migration: DB scanning completed. Total time taken:...`

You DO NOT immediately need to do anything. Your node will continue to run
and will operate normally. However we added a verification tool that will iterate through
all keys in iavl at a specific height and verify they exist in State Store.

You should run the following command for a selection of different heights
```bash
enid tools verify-migration --version $VERIFICATION_HEIGHT
```

This will output `Verification Succeeded` if the verification was successful.


### Step 6: Restart enid normally and verify node runs properly
Once the verification has completed, we can restart enid normally and verify
that the node operates.

If you are using systemd, make sure to update your service configuration to use this command:
```bash
enid start --chain-id pacific-1
```

Note how we are not using the `--migrate-iavl` and `--migration-height` flags.
We can let this run for a couple hours and verify node oeprates normally.


### Step 7: Clear out Iavl and restart enid
Once it has been confirmed that the node has been running normally,
we can proceed to clear out the iavl and restart enid normally.

```bash
systemctl stop enid
rm -rf ~/.eni/data/application.db
enid start --chain-id pacific-1
```


## Metrics

During the State Store Migration, there are exported metrics that are helpful to keep track of
the progress.

`eni_migration_leaf_nodes_exported` keeps track of how many nodes have been exported from iavl.

`eni_migration_nodes_imported` keeps track of how many nodes have been imported into EniDB (pebbledb).

Both of these metrics have a `module` label which indicates what module is currently being exported / imported.


## FAQ

### Can the state store migration be stopped and restarted?

The state store migration can be stopped and restarted at any time. The migration
process saves the latest `module` and `key` written to State Store (pebbledb) and will
automatically resume the migration from that latest key once restarted.

All one needs to do is restart enid with the migration command as in step 4
```bash
enid start --migrate-iavl --migrate-height $MIGRATION_HEIGHT --chain-id pacific-1
```
