package dao

import (
	"context"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

// Trace mapped from table <traces>
type Trace struct {
	TransactionHash   string         `gorm:"column:transaction_hash" json:"transaction_hash"`
	TransactionIndex  int32          `gorm:"column:transaction_index" json:"transaction_index"`
	BlockNumber       int32          `gorm:"column:block_number;not null" json:"block_number"`
	BlockHash         string         `gorm:"column:block_hash;not null" json:"block_hash"`
	BlockTime         time.Time      `gorm:"column:block_time;not null" json:"block_time"`
	FromAddress       string         `gorm:"column:from_address" json:"from_address"`
	ToAddress         string         `gorm:"column:to_address" json:"to_address"`
	Value             float64        `gorm:"column:value;not null" json:"value"`
	Input             string         `gorm:"column:input" json:"input"`
	Output            string         `gorm:"column:output" json:"output"`
	TraceType         string         `gorm:"column:trace_type;not null" json:"trace_type"`
	CallType          string         `gorm:"column:call_type" json:"call_type"`
	RewardType        string         `gorm:"column:reward_type" json:"reward_type"`
	Gas               float64        `gorm:"column:gas" json:"gas"`
	GasUsed           int64          `gorm:"column:gas_used" json:"gas_used"`
	Subtraces         int64          `gorm:"column:subtraces;not null" json:"subtraces"`
	TraceAddress      pq.StringArray `gorm:"column:trace_address;type:text[]" json:"trace_address"`
	Error             string         `gorm:"column:error" json:"error"`
	Status            int32          `gorm:"column:status;not null" json:"status"`
	TransactionStatus int32          `gorm:"column:transaction_status" json:"transaction_status"`
	FuncName          string         `gorm:"column:func_name" json:"func_name"`
	FuncSignature     string         `gorm:"column:func_signature" json:"func_signature"`
	ArgumentNames     pq.StringArray `gorm:"column:argument_names;type:text[]" json:"argument_names"`
	ArgumentTypes     pq.StringArray `gorm:"column:argument_types;type:text[]" json:"argument_types"`
	ArgumentValues    pq.StringArray `gorm:"column:argument_values;type:text[]" json:"argument_values"`
	OutputParameters  string         `gorm:"column:output_parameters" json:"output_parameters"`
	OutputNames       pq.StringArray `gorm:"column:output_names;type:text[]" json:"output_names"`
	OutputTypes       pq.StringArray `gorm:"column:output_types;type:text[]" json:"output_types"`
	OutputValues      pq.StringArray `gorm:"column:output_values;type:text[]" json:"output_values"`
	TraceID           string         `gorm:"column:trace_id;primaryKey" json:"trace_id"`
	TraceIndex        int32          `gorm:"column:trace_index" json:"trace_index"`
	DecodedFromAbi    bool           `gorm:"column:decoded_from_abi" json:"decoded_from_abi"`
	ProcessTime       time.Time      `gorm:"column:process_time" json:"process_time"`
	BlockDate         time.Time      `gorm:"column:block_date;primaryKey" json:"block_date"`
}

// TableName Trace's table name
func (*Trace) TableName() string {
	return "ethereum_holesky.traces"
}

type TraceDao struct {
	sourceDB  *gorm.DB
	replicaDB []*gorm.DB
	m         *Trace
}

func NewTraceDao(ctx context.Context, dbs ...*gorm.DB) *TraceDao {
	dao := new(TraceDao)
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

func (d *TraceDao) Create(ctx context.Context, obj *Trace) error {
	err := d.sourceDB.Model(d.m).Create(&obj).Error
	if err != nil {
		return fmt.Errorf("TraceDao: %w", err)
	}
	return nil
}

func (d *TraceDao) Get(ctx context.Context, fields, where string) (*Trace, error) {
	items, err := d.List(ctx, fields, where, 0, 1)
	if err != nil {
		return nil, fmt.Errorf("TraceDao: Get where=%s: %w", where, err)
	}
	if len(items) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &items[0], nil
}

func (d *TraceDao) List(ctx context.Context, fields, where string, offset, limit int) ([]Trace, error) {
	var results []Trace
	err := d.replicaDB[rand.Intn(len(d.replicaDB))].Model(d.m).
		Select(fields).Where(where).Offset(offset).Limit(limit).Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("TraceDao: List where=%s: %w", where, err)
	}
	return results, nil
}

func (d *TraceDao) Update(ctx context.Context, where string, update map[string]interface{}, args ...interface{}) error {
	err := d.sourceDB.Model(d.m).Where(where, args...).
		Updates(update).Error
	if err != nil {
		return fmt.Errorf("TraceDao:Update where=%s: %w", where, err)
	}
	return nil
}

func (d *TraceDao) Delete(ctx context.Context, where string, args ...interface{}) error {
	if len(where) == 0 {
		return fmt.Errorf("TraceDao: Delete where=%s", where)
	}
	if err := d.sourceDB.Where(where, args...).Delete(d.m).Error; err != nil {
		return fmt.Errorf("TraceDao: Delete where=%s: %w", where, err)
	}
	return nil
}
