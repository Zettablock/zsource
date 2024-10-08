// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package stellar

import (
	"time"
)

const TableNameHistoryLedger = "history_ledgers"

// HistoryLedger mapped from table <history_ledgers>
type HistoryLedger struct {
	Number                     int64     `gorm:"column:number;primaryKey" json:"number"`
	ID                         string    `gorm:"column:id;not null" json:"id"`
	PagingToken                string    `gorm:"column:paging_token;not null" json:"paging_token"`
	Hash                       string    `gorm:"column:hash;not null" json:"hash"`
	PrevHash                   string    `gorm:"column:prev_hash;not null" json:"prev_hash"`
	TransactionCount           int32     `gorm:"column:transaction_count" json:"transaction_count"`
	SuccessfulTransactionCount int32     `gorm:"column:successful_transaction_count" json:"successful_transaction_count"`
	FailedTransactionCount     int32     `gorm:"column:failed_transaction_count" json:"failed_transaction_count"`
	OperationCount             int32     `gorm:"column:operation_count" json:"operation_count"`
	TxSetOperationCount        int32     `gorm:"column:tx_set_operation_count" json:"tx_set_operation_count"`
	Timestamp                  time.Time `gorm:"column:timestamp" json:"timestamp"`
	TotalCoins                 float64   `gorm:"column:total_coins" json:"total_coins"`
	FeePool                    float64   `gorm:"column:fee_pool" json:"fee_pool"`
	BaseFeeInStroops           int32     `gorm:"column:base_fee_in_stroops" json:"base_fee_in_stroops"`
	BaseReserveInStroops       int32     `gorm:"column:base_reserve_in_stroops" json:"base_reserve_in_stroops"`
	MaxTxSetSize               int32     `gorm:"column:max_tx_set_size" json:"max_tx_set_size"`
	ProtocolVersion            int32     `gorm:"column:protocol_version" json:"protocol_version"`
	HeaderXdr                  string    `gorm:"column:header_xdr" json:"header_xdr"`
	ProcessTime                time.Time `gorm:"column:process_time" json:"process_time"`
	BlockDate                  time.Time `gorm:"column:block_date;not null" json:"block_date"`
}

// TableName HistoryLedger's table name
func (*HistoryLedger) TableName() string {
	return TableNameHistoryLedger
}
