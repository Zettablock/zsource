package testutils

import (
	"github.com/Zettablock/zsource/dao/ethereum"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/gorm"
)

// These tables will be automatically created in each schema.
var tableModels = []interface{}{
	&ethereum.Block{},
	&ethereum.Log{},
	&ethereum.Trace{},
	&ethereum.Transaction{},
}

// EthereumData is used to hold the data provided by the user for testing the
// Ethereum handlers and initializing the database according to the data (by
// creating tables and inserting data).
type EthereumData struct {
	// Map from schema name to the data for that schema.
	data map[string]*EthereumSchemaData
}

func NewEthereumData() *EthereumData {
	return &EthereumData{
		data: make(map[string]*EthereumSchemaData),
	}
}

func (d *EthereumData) AddSchemaData(schemaName string, blocks []*ethereum.Block, logs []*ethereum.Log) {
	d.data[schemaName] = NewEthereumSchemaData(schemaName, blocks, logs)
}

func (d *EthereumData) AddSchemaDataEmpty(schemaName string) {
	d.data[schemaName] = NewEthereumSchemaData(schemaName, []*ethereum.Block{}, []*ethereum.Log{})
}

func (d *EthereumData) PopulateDb(container *postgres.PostgresContainer, db *gorm.DB) error {
	if err := d.initDb(container); err != nil {
		return err
	}
	for schemaName, schemaData := range d.data {
		if err := schemaData.PopulateDb(db, schemaName); err != nil {
			return err
		}
	}
	return nil
}

// This will create the corresponding tables if they do not already exist and
// migrate the schema according to the model if they already exist (highly
// unlikely during test).
func (d *EthereumData) initDb(container *postgres.PostgresContainer) error {
	for schemaName := range d.data {
		// We open a new gorm.DB for every schema because the AutoMigrate call
		// doesn't support changing the default schema while we still want to use
		// AutoMigrate for its ability to automatically create the corresponding
		// table with only the model struct (no need to provide explicit CREATE
		// TABLE statement).
		db, err := GetDbFromContianer(container, schemaName)
		if err != nil {
			return err
		}
		if err := db.Exec("CREATE SCHEMA IF NOT EXISTS " + schemaName).Error; err != nil {
			return err
		}
		if err := db.AutoMigrate(tableModels...); err != nil {
			return err
		}
	}
	return nil
}

// EthereumSchemaData is for a single schema.
type EthereumSchemaData struct {
	schemaName string
	blocks     []*ethereum.Block
	logs       []*ethereum.Log
}

func NewEthereumSchemaData(schemaName string, blocks []*ethereum.Block, logs []*ethereum.Log) *EthereumSchemaData {
	if blocks == nil {
		blocks = []*ethereum.Block{}
	}
	if logs == nil {
		logs = []*ethereum.Log{}
	}
	return &EthereumSchemaData{
		schemaName: schemaName,
		blocks:     blocks,
		logs:       logs,
	}
}

// PopulateDb populates the database with the blocks and logs provided.
func (d *EthereumSchemaData) PopulateDb(db *gorm.DB, schemaName string) error {
	if len(d.blocks) != 0 {
		if err := db.Table(schemaName + ".blocks").Create(d.blocks).Error; err != nil {
			return err
		}
	}
	if len(d.logs) != 0 {
		if err := db.Table(schemaName + ".logs").Create(d.logs).Error; err != nil {
			return err
		}
	}
	return nil
}
