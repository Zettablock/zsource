// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package stellar

import (
	"time"
)

const TableNameClaimableBalance = "claimable_balances"

// ClaimableBalance mapped from table <claimable_balances>
type ClaimableBalance struct {
	ID                 string    `gorm:"column:id;primaryKey" json:"id"`
	PagingToken        string    `gorm:"column:paging_token;not null" json:"paging_token"`
	Asset              string    `gorm:"column:asset;not null" json:"asset"`
	Amount             float64   `gorm:"column:amount" json:"amount"`
	Sponsor            string    `gorm:"column:sponsor" json:"sponsor"`
	LastModifiedLedger int64     `gorm:"column:last_modified_ledger" json:"last_modified_ledger"`
	LastModifiedTime   time.Time `gorm:"column:last_modified_time" json:"last_modified_time"`
	Claimants          string    `gorm:"column:claimants" json:"claimants"`
	Flags              string    `gorm:"column:flags" json:"flags"`
	ProcessTime        time.Time `gorm:"column:process_time" json:"process_time"`
	BlockDate          time.Time `gorm:"column:block_date;not null" json:"block_date"`
}

// TableName ClaimableBalance's table name
func (*ClaimableBalance) TableName() string {
	return TableNameClaimableBalance
}
