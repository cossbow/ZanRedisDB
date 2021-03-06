package node

import (
	"github.com/youzan/ZanRedisDB/common"
	"github.com/youzan/ZanRedisDB/engine"
)

type NamespaceConfig struct {
	// namespace full name with partition
	Name string `json:"name"`
	// namespace name without partition
	BaseName         string          `json:"base_name"`
	EngType          string          `json:"eng_type"`
	PartitionNum     int             `json:"partition_num"`
	SnapCount        int             `json:"snap_count"`
	SnapCatchup      int             `json:"snap_catchup"`
	Replicator       int             `json:"replicator"`
	OptimizedFsync   bool            `json:"optimized_fsync"`
	RaftGroupConf    RaftGroupConfig `json:"raft_group_conf"`
	ExpirationPolicy string          `json:"expiration_policy"`
	DataVersion      string          `json:"data_version"`
}

func NewNSConfig() *NamespaceConfig {
	return &NamespaceConfig{
		SnapCount:        common.DefaultSnapCount,
		SnapCatchup:      common.DefaultSnapCatchup,
		ExpirationPolicy: common.DefaultExpirationPolicy,
	}
}

type NamespaceDynamicConf struct {
	Replicator int
}

type RaftGroupConfig struct {
	GroupID   uint64        `json:"group_id"`
	SeedNodes []ReplicaInfo `json:"seed_nodes"`
}

type MachineConfig struct {
	// server node id
	NodeID                 uint64             `json:"node_id"`
	BroadcastAddr          string             `json:"broadcast_addr"`
	HttpAPIPort            int                `json:"http_api_port"`
	LocalRaftAddr          string             `json:"local_raft_addr"`
	DataRootDir            string             `json:"data_root_dir"`
	ElectionTick           int                `json:"election_tick"`
	TickMs                 int                `json:"tick_ms"`
	KeepBackup             int                `json:"keep_backup"`
	KeepWAL                int                `json:"keep_wal"`
	UseRocksWAL            bool               `json:"use_rocks_wal"`
	SharedRocksWAL         bool               `json:"shared_rocks_wal"`
	LearnerRole            string             `json:"learner_role"`
	RemoteSyncCluster      string             `json:"remote_sync_cluster"`
	StateMachineType       string             `json:"state_machine_type"`
	RocksDBOpts            engine.RockOptions `json:"rocksdb_opts"`
	RocksDBSharedConfig    engine.SharedRockConfig
	WALRocksDBOpts         engine.RockOptions `json:"wal_rocksdb_opts"`
	WALRocksDBSharedConfig engine.SharedRockConfig
}

type ReplicaInfo struct {
	NodeID    uint64 `json:"node_id"`
	ReplicaID uint64 `json:"replica_id"`
	RaftAddr  string `json:"raft_addr"`
}

type RaftConfig struct {
	GroupID uint64 `json:"group_id"`
	// group name is combined namespace-partition string
	GroupName string `json:"group_name"`
	// this is replica id
	ID uint64 `json:"id"`
	// local server transport address, it
	// can be used by several raft group
	RaftAddr       string                 `json:"raft_addr"`
	DataDir        string                 `json:"data_dir"`
	WALDir         string                 `json:"wal_dir"`
	SnapDir        string                 `json:"snap_dir"`
	RaftStorageDir string                 `json:"raft_storage_dir"`
	RaftPeers      map[uint64]ReplicaInfo `json:"raft_peers"`
	SnapCount      int                    `json:"snap_count"`
	SnapCatchup    int                    `json:"snap_catchup"`
	Replicator     int32                  `json:"replicator"`
	OptimizedFsync bool                   `json:"optimized_fsync"`
	rockEng        engine.KVEngine
	nodeConfig     *MachineConfig
}

func (rc *RaftConfig) SetEng(eng engine.KVEngine) {
	rc.rockEng = eng
}
