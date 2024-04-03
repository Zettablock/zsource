package model

import (
	"context"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

// Log mapped from table <logs>
type Log struct {
	TransactionHash  string         `gorm:"column:transaction_hash;primaryKey" json:"transaction_hash"`
	TransactionIndex int32          `gorm:"column:transaction_index;not null" json:"transaction_index"`
	BlockNumber      int64          `gorm:"column:block_number;primaryKey" json:"block_number"`
	BlockHash        string         `gorm:"column:block_hash;not null" json:"block_hash"`
	Removed          bool           `gorm:"column:removed" json:"removed"`
	LogIndex         int32          `gorm:"column:log_index;primaryKey" json:"log_index"`
	Data             string         `gorm:"column:data" json:"data"`
	Topics           pq.StringArray `gorm:"column:topics;type:text[]" json:"topics"`
	ContractAddress  string         `gorm:"column:contract_address;not null" json:"contract_address"`
	Anonymous        bool           `gorm:"column:anonymous" json:"anonymous"`
	Event            string         `gorm:"column:event" json:"event"`
	EventSignature   string         `gorm:"column:event_signature" json:"event_signature"`
	ArgumentNames    pq.StringArray `gorm:"column:argument_names;type:text[]" json:"argument_names"`
	ArgumentTypes    pq.StringArray `gorm:"column:argument_types;type:text[]" json:"argument_types"`
	ArgumentValues   pq.StringArray `gorm:"column:argument_values;type:text[]" json:"argument_values"`
	BlockTime        time.Time      `gorm:"column:block_time;not null" json:"block_time"`
	DecodedFromAbi   bool           `gorm:"column:decoded_from_abi" json:"decoded_from_abi"`
	ProcessTime      time.Time      `gorm:"column:process_time" json:"process_time"`
	BlockDate        time.Time      `gorm:"column:block_date" json:"block_date"`
}

// TableName Log's table name
func (*Log) TableName() string {
	return "ethereum_holesky.logs"
}

type LogDao struct {
	sourceDB  *gorm.DB
	replicaDB []*gorm.DB
	m         *Log
}

func NewLogDao(ctx context.Context, dbs ...*gorm.DB) *LogDao {
	dao := new(LogDao)
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

func (d *LogDao) Create(ctx context.Context, obj *Log) error {
	err := d.sourceDB.Model(d.m).Create(&obj).Error
	if err != nil {
		return fmt.Errorf("LogDao: %w", err)
	}
	return nil
}

func (d *LogDao) Get(ctx context.Context, fields, where string) (*Log, error) {
	items, err := d.List(ctx, fields, where, 0, 1)
	if err != nil {
		return nil, fmt.Errorf("LogDao: Get where=%s: %w", where, err)
	}
	if len(items) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &items[0], nil
}

func (d *LogDao) List(ctx context.Context, fields, where string, offset, limit int) ([]Log, error) {
	var results []Log
	err := d.replicaDB[rand.Intn(len(d.replicaDB))].Model(d.m).
		Select(fields).Where(where).Offset(offset).Limit(limit).Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("LogDao: List where=%s: %w", where, err)
	}
	return results, nil
}

func (d *LogDao) Update(ctx context.Context, where string, update map[string]interface{}, args ...interface{}) error {
	err := d.sourceDB.Model(d.m).Where(where, args...).
		Updates(update).Error
	if err != nil {
		return fmt.Errorf("LogDao:Update where=%s: %w", where, err)
	}
	return nil
}

func (d *LogDao) Delete(ctx context.Context, where string, args ...interface{}) error {
	if len(where) == 0 {
		return fmt.Errorf("LogDao: Delete where=%s", where)
	}
	if err := d.sourceDB.Where(where, args...).Delete(d.m).Error; err != nil {
		return fmt.Errorf("LogDao: Delete where=%s: %w", where, err)
	}
	return nil
}
