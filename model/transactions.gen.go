package model

import (
	"context"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

// Transaction mapped from table <transactions>
type Transaction struct {
	Hash                   string         `gorm:"column:hash;primaryKey" json:"hash"`
	Nonce                  int64          `gorm:"column:nonce;not null" json:"nonce"`
	BlockHash              string         `gorm:"column:block_hash;not null" json:"block_hash"`
	BlockNumber            int64          `gorm:"column:block_number;not null" json:"block_number"`
	TransactionIndex       int32          `gorm:"column:transaction_index;not null" json:"transaction_index"`
	FromAddress            string         `gorm:"column:from_address;not null" json:"from_address"`
	ToAddress              string         `gorm:"column:to_address;not null" json:"to_address"`
	Value                  float64        `gorm:"column:value;not null" json:"value"`
	Type                   string         `gorm:"column:type" json:"type"`
	GasPrice               float64        `gorm:"column:gas_price;not null" json:"gas_price"`
	Input                  string         `gorm:"column:input" json:"input"`
	V                      string         `gorm:"column:v;not null" json:"v"`
	S                      string         `gorm:"column:s;not null" json:"s"`
	R                      string         `gorm:"column:r;not null" json:"r"`
	MaxFeePerGas           float64        `gorm:"column:max_fee_per_gas" json:"max_fee_per_gas"`
	MaxPriorityFeePerGas   float64        `gorm:"column:max_priority_fee_per_gas" json:"max_priority_fee_per_gas"`
	ChainID                string         `gorm:"column:chain_id" json:"chain_id"`
	AccessList             string         `gorm:"column:access_list;type:text[]" json:"access_list"`
	GasLimit               int64          `gorm:"column:gas_limit;not null" json:"gas_limit"`
	FuncName               string         `gorm:"column:func_name" json:"func_name"`
	FuncSignature          string         `gorm:"column:func_signature" json:"func_signature"`
	ArgumentNames          pq.StringArray `gorm:"column:argument_names;type:text[]" json:"argument_names"`
	ArgumentTypes          pq.StringArray `gorm:"column:argument_types;type:text[]" json:"argument_types"`
	ArgumentValues         pq.StringArray `gorm:"column:argument_values;type:text[]" json:"argument_values"`
	BlockTime              time.Time      `gorm:"column:block_time;not null" json:"block_time"`
	Status                 int32          `gorm:"column:status" json:"status"`
	GasUsed                int64          `gorm:"column:gas_used;not null" json:"gas_used"`
	CumulativeGasUsed      int64          `gorm:"column:cumulative_gas_used;not null" json:"cumulative_gas_used"`
	EffectiveGasPrice      float64        `gorm:"column:effective_gas_price" json:"effective_gas_price"`
	ReceiptContractAddress string         `gorm:"column:receipt_contract_address" json:"receipt_contract_address"`
	DecodedFromAbi         bool           `gorm:"column:decoded_from_abi" json:"decoded_from_abi"`
	ProcessTime            time.Time      `gorm:"column:process_time" json:"process_time"`
	BlockDate              time.Time      `gorm:"column:block_date" json:"block_date"`
}

// TableName Transaction's table name
func (*Transaction) TableName() string {
	return "ethereum_holesky.transactions"
}

type TransactionDao struct {
	sourceDB  *gorm.DB
	replicaDB []*gorm.DB
	m         *Transaction
}

func NewTransactionDao(ctx context.Context, dbs ...*gorm.DB) *TransactionDao {
	dao := new(TransactionDao)
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

func (d *TransactionDao) Create(ctx context.Context, obj *Transaction) error {
	err := d.sourceDB.Model(d.m).Create(&obj).Error
	if err != nil {
		return fmt.Errorf("TransactionDao: %w", err)
	}
	return nil
}

func (d *TransactionDao) Get(ctx context.Context, fields, where string) (*Transaction, error) {
	items, err := d.List(ctx, fields, where, 0, 1)
	if err != nil {
		return nil, fmt.Errorf("TransactionDao: Get where=%s: %w", where, err)
	}
	if len(items) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &items[0], nil
}

func (d *TransactionDao) List(ctx context.Context, fields, where string, offset, limit int) ([]Transaction, error) {
	var results []Transaction
	err := d.replicaDB[rand.Intn(len(d.replicaDB))].Model(d.m).
		Select(fields).Where(where).Offset(offset).Limit(limit).Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("TransactionDao: List where=%s: %w", where, err)
	}
	return results, nil
}

func (d *TransactionDao) Update(ctx context.Context, where string, update map[string]interface{}, args ...interface{}) error {
	err := d.sourceDB.Model(d.m).Where(where, args...).
		Updates(update).Error
	if err != nil {
		return fmt.Errorf("TransactionDao:Update where=%s: %w", where, err)
	}
	return nil
}

func (d *TransactionDao) Delete(ctx context.Context, where string, args ...interface{}) error {
	if len(where) == 0 {
		return fmt.Errorf("TransactionDao: Delete where=%s", where)
	}
	if err := d.sourceDB.Where(where, args...).Delete(d.m).Error; err != nil {
		return fmt.Errorf("TransactionDao: Delete where=%s: %w", where, err)
	}
	return nil
}
