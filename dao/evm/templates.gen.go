package evm

const TableNameTemplate = "templates"

// Template mapped from table <templates>
type Template struct {
	Name            string `gorm:"column:name;primaryKey" json:"name"`
	ContractAddress string `gorm:"column:contract_address;primaryKey" json:"contract_address"`
	EventName       string `gorm:"column:event_name;primaryKey" json:"event_name"`
}

// TableName Template's table name
func (*Template) TableName() string {
	return TableNameTemplate
}
