package app

//
//import (
//	"testing"
//
//	"github.com/eni-chain/eni-db/config"
//	"github.com/stretchr/testify/assert"
//)
//
//type TestEniDBAppOpts struct {
//}
//
//func (t TestEniDBAppOpts) Get(s string) interface{} {
//	switch s {
//	case FlagSCEnable:
//		return config.DefaultStateCommitConfig().Enable
//	case FlagSCAsyncCommitBuffer:
//		return config.DefaultStateCommitConfig().AsyncCommitBuffer
//	case FlagSCDirectory:
//		return config.DefaultStateCommitConfig().Directory
//	case FlagSCCacheSize:
//		return config.DefaultStateCommitConfig().CacheSize
//	case FlagSCSnapshotInterval:
//		return config.DefaultStateCommitConfig().SnapshotInterval
//	case FlagSCSnapshotKeepRecent:
//		return config.DefaultStateCommitConfig().SnapshotKeepRecent
//	case FlagSCSnapshotWriterLimit:
//		return config.DefaultStateCommitConfig().SnapshotWriterLimit
//	case FlagSSEnable:
//		return config.DefaultStateStoreConfig().Enable
//	case FlagSSBackend:
//		return config.DefaultStateStoreConfig().Backend
//	case FlagSSAsyncWriterBuffer:
//		return config.DefaultStateStoreConfig().AsyncWriteBuffer
//	case FlagSSDirectory:
//		return config.DefaultStateStoreConfig().DBDirectory
//	case FlagSSKeepRecent:
//		return config.DefaultStateStoreConfig().KeepRecent
//	case FlagSSPruneInterval:
//		return config.DefaultStateStoreConfig().PruneIntervalSeconds
//	case FlagSSImportNumWorkers:
//		return config.DefaultStateStoreConfig().ImportNumWorkers
//	}
//	return nil
//}
//
//func TestNewDefaultConfig(t *testing.T) {
//	// Make sure when adding a new default config, it should apply to EniDB during initialization
//	appOpts := TestEniDBAppOpts{}
//	scConfig := parseSCConfigs(appOpts)
//	ssConfig := parseSSConfigs(appOpts)
//	assert.Equal(t, scConfig, config.DefaultStateCommitConfig())
//	assert.Equal(t, ssConfig, config.DefaultStateStoreConfig())
//}
