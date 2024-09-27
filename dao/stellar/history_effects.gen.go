// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package stellar

import (
	"time"
)

const TableNameHistoryEffect = "history_effects"

// HistoryEffect mapped from table <history_effects>
type HistoryEffect struct {
	ID           string    `gorm:"column:id;primaryKey" json:"id"`
	PagingToken  string    `gorm:"column:paging_token;not null" json:"paging_token"`
	OperationID  int64     `gorm:"column:operation_id;not null" json:"operation_id"`
	Sequence     int32     `gorm:"column:sequence;not null" json:"sequence"`
	Account      string    `gorm:"column:account;not null" json:"account"`
	Type         string    `gorm:"column:type;not null" json:"type"`
	TypeI        int32     `gorm:"column:type_i;not null" json:"type_i"`
	CreatedAt    time.Time `gorm:"column:created_at;not null" json:"created_at"`
	Details      string    `gorm:"column:details" json:"details"`
	LedgerNumber int64     `gorm:"column:ledger_number;not null" json:"ledger_number"`
	ProcessTime  time.Time `gorm:"column:process_time" json:"process_time"`
	BlockDate    time.Time `gorm:"column:block_date;not null" json:"block_date"`
}

// TableName HistoryEffect's table name
func (*HistoryEffect) TableName() string {
	return TableNameHistoryEffect
}