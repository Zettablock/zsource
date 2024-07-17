package base

import (
	"fmt"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/Zettablock/zclient-base/dao/nft"
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
	BlockTime        time.Time      `gorm:"column:block_time;not null;type:timestamp" json:"block_time"`
	DecodedFromAbi   bool           `gorm:"column:decoded_from_abi" json:"decoded_from_abi"`
	ProcessTime      time.Time      `gorm:"column:process_time;type:timestamp" json:"process_time"`
	BlockDate        time.Time      `gorm:"column:block_date;type:timestamp" json:"block_date"`
}

type LogDao struct {
	sourceDB  *gorm.DB
	replicaDB []*gorm.DB
	m         *Log
	schema    *nft.DbSchema
}

const MintAddress = "0x0000000000000000000000000000000000000000"
const TransferEvent = "Transfer"
const TransferSingleEvent = "TransferSingle"
const TransferBatchEvent = "TransferBatch"
const TransferEventTopic = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
const Erc1155TransferSingleEventTopic = "0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62"
const Erc1155TransferBatchEventTopic = "0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb"

func (d *LogDao) TableName() string {
	return fmt.Sprintf("%s.logs", d.schema.Source)
}

func NewLogDao(schema *nft.DbSchema, dbs ...*gorm.DB) *LogDao {
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
	dao.schema = schema
	return dao
}

func (d *LogDao) List(blockNumber int64) []Log {
	var o []Log
	query := fmt.Sprintf(`
        SELECT transaction_hash, transaction_index, block_number, log_index, data, topics,
		contract_address, event, argument_values, block_time, data_creation_date
        FROM %s.logs
        WHERE block_number = %d
			and event = 'Transfer'
		ORDER BY log_index asc;
    `, d.schema.Source, blockNumber)
	d.sourceDB.Raw(query).Scan(&o)
	return o
}
