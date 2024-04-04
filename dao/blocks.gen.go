package dao

import (
	"context"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"math/rand"
	"time"
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
	Timestamp         time.Time      `gorm:"column:timestamp;not null" json:"timestamp"`
	Uncles            pq.StringArray `gorm:"column:uncles;type:text[]" json:"uncles"`
	NumOfTransactions int32          `gorm:"column:num_of_transactions;not null" json:"num_of_transactions"`
	ExtraDataRaw      string         `gorm:"column:extra_data_raw" json:"extra_data_raw"`
	ExtraData         string         `gorm:"column:extra_data" json:"extra_data"`
	ProcessTime       time.Time      `gorm:"column:process_time" json:"process_time"`
	BlockDate         time.Time      `gorm:"column:block_date" json:"block_date"`
}

// TableName Block's table name
func (*Block) TableName() string {
	return "ethereum_holesky.blocks"
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

func (d *BlockDao) Create(ctx context.Context, obj *Block) error {
	err := d.sourceDB.Model(d.m).Create(&obj).Error
	if err != nil {
		return fmt.Errorf("BlockDao: %w", err)
	}
	return nil
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

func (d *BlockDao) Update(ctx context.Context, where string, update map[string]interface{}, args ...interface{}) error {
	err := d.sourceDB.Model(d.m).Where(where, args...).
		Updates(update).Error
	if err != nil {
		return fmt.Errorf("BlockDao:Update where=%s: %w", where, err)
	}
	return nil
}

func (d *BlockDao) Delete(ctx context.Context, where string, args ...interface{}) error {
	if len(where) == 0 {
		return fmt.Errorf("BlockDao: Delete where=%s", where)
	}
	if err := d.sourceDB.Where(where, args...).Delete(d.m).Error; err != nil {
		return fmt.Errorf("BlockDao: Delete where=%s: %w", where, err)
	}
	return nil
}
