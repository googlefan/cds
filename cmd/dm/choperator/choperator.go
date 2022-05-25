package choperator

import (
	"github.com/zeromicro/cds/cmd/dm/util"
	"github.com/zeromicro/cds/pkg/ckgroup"
	ckcfg "github.com/zeromicro/cds/pkg/ckgroup/config"
	"github.com/zeromicro/go-zero/core/logx"
)

type ChOperator interface {
	MysqlBatchInsert(insertData [][]interface{}, insertQuery string, arr []util.DataType, indexOfPrimKeys int) error
	ObtainClickHouseKV(targetDB, targetTable string) (map[string]string, error)
	BatchInsert(insertData [][]interface{}, insertQuery string, indexOfPrimKey int) error
}

func NewChOperator(shards [][]string) (ChOperator, error) {
	shardCfgs := make([]ckcfg.ShardGroupConfig, 0, len(shards))
	for _, i := range shards {
		cfg := ckcfg.ShardGroupConfig{ReplicaNodes: make([]string, 0, len(i)-1)}
		for index, addr := range i {
			if index == 0 {
				cfg.ShardNode = addr
			} else {
				cfg.ReplicaNodes = append(cfg.ReplicaNodes, addr)
			}
		}
		shardCfgs = append(shardCfgs, cfg)
	}
	ckConfig := ckcfg.Config{
		ShardGroups: shardCfgs,
		QueryNode:   shardCfgs[0].ShardNode,
	}
	ch, err := ckgroup.NewCKGroup(ckConfig)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	var cgo CkGroupOperator
	cgo.ckGroup = ch
	return &cgo, nil
}
