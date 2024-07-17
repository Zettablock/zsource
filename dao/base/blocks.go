package base

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Block mapped from table <blocks>
type Block struct {
	Number            int64          `gorm:"column:number;primaryKey" json:"number"`
	Hash              string         `gorm:"column:hash;not null" json:"hash"`
	ParentHash        string         `gorm:"column:parent_hash;not null" json:"parent_hash"`
	Nonce             string         `gorm:"column:nonce;not null" json:"nonce"`
	MixHash           string         `gorm:"column:mix_hash;not null" json:"mix_hash"`
	Sha3Uncles        string         `gorm:"column:sha3_uncles;not null" json:"sha3_uncles"`
	LogsBloom         string         `gorm:"column:logs_bloom" json:"logs_bloom"`
	TransactionsRoot  string         `gorm:"column:transactions_root;not null" json:"transactions_root"`
	StateRoot         string         `gorm:"column:state_root;not null" json:"state_root"`
	ReceiptsRoot      string         `gorm:"column:receipts_root;not null" json:"receipts_root"`
	Miner             string         `gorm:"column:miner;not null" json:"miner"`
	Difficulty        float64        `gorm:"column:difficulty;not null" json:"difficulty"`
	TotalDifficulty   float64        `gorm:"column:total_difficulty;not null" json:"total_difficulty"`
	Size              int64          `gorm:"column:size;not null" json:"size"`
	GasLimit          int64          `gorm:"column:gas_limit;not null" json:"gas_limit"`
	GasUsed           int64          `gorm:"column:gas_used;not null" json:"gas_used"`
	BaseFeePerGas     int64          `gorm:"column:base_fee_per_gas" json:"base_fee_per_gas"`
	Timestamp         time.Time      `gorm:"column:timestamp;not null;type:timestamp" json:"timestamp"`
	Uncles            pq.StringArray `gorm:"column:uncles;type:text[]" json:"uncles"`
	NumOfTransactions int32          `gorm:"column:num_of_transactions;not null" json:"num_of_transactions"`
	ExtraDataRaw      string         `gorm:"column:extra_data_raw" json:"extra_data_raw"`
	ExtraData         string         `gorm:"column:extra_data" json:"extra_data"`
	ProcessTime       time.Time      `gorm:"column:process_time;type:timestamp" json:"process_time"`
	BlockDate         time.Time      `gorm:"column:block_date;type:timestamp" json:"block_date"`
}

type BlockDao struct {
	sourceDB  *gorm.DB
	replicaDB []*gorm.DB
	m         *Block
}

func NewBlockDao(ctx context.Context, dbs ...*gorm.DB) *BlockDao {
	dao := new(BlockDao)
	switch len(dbs) {
	case 0:
		panic("database connection required")
	case 1:
		dao.sourceDB = dbs[0]
		dao.replicaDB = []*gorm.DB{dbs[0]}
	default:
		dao.sourceDB = dbs[0]
		dao.replicaDB = dbs[1:]
	}
	return dao
}

func (d *BlockDao) Get(ctx context.Context, fields, where string) (*Block, error) {
	items, err := d.List(ctx, fields, where, 0, 1)
	if err != nil {
		return nil, fmt.Errorf("BlockDao: Get where=%s: %w", where, err)
	}
	if len(items) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &items[0], nil
}

func (d *BlockDao) List(ctx context.Context, fields, where string, offset, limit int) ([]Block, error) {
	var results []Block
	err := d.replicaDB[rand.Intn(len(d.replicaDB))].Model(d.m).
		Select(fields).Where(where).Offset(offset).Limit(limit).Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("BlockDao: List where=%s: %w", where, err)
	}
	return results, nil
}
