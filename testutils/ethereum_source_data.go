package testutils

import (
	"context"

	"github.com/Zettablock/zsource/dao/ethereum"
	"gorm.io/gorm"
)

// This struct is used to hold the data provided by the user for testing the
// Ethereum handlers and initializing the database according to the data (by
// creating tables and inserting data).
type EthereumData struct {
	blocks []*ethereum.Block
	logs   []*ethereum.Log
}

func NewEthereumData(blocks []*ethereum.Block, logs []*ethereum.Log) *EthereumData {
	if blocks == nil {
		blocks = []*ethereum.Block{}
	}
	if logs == nil {
		logs = []*ethereum.Log{}
	}
	return &EthereumData{
		blocks: blocks,
		logs:   logs,
	}
}

func NewEthereumDataEmpty() *EthereumData {
	return NewEthereumData(
		[]*ethereum.Block{},
		[]*ethereum.Log{},
	)
}

// Populates the database with the blocks and logs provided.
func (d *EthereumData) PopulateDb(db *gorm.DB) error {
	err := d.initDb(db)
	if err != nil {
		return err
	}

	blockDao := ethereum.NewBlockDao(context.Background(), db)
	for _, block := range d.blocks {
		if err := blockDao.Create(context.Background(), block); err != nil {
			return err
		}
	}

	logDao := ethereum.NewLogDao(context.Background(), db)
	for _, log := range d.logs {
		if err := logDao.Create(context.Background(), log); err != nil {
			return err
		}
	}

	return nil
}

// This will create the corresponding tables if they do not already exist and
// migrate the schema according to the model if they already exist (highly
// unlikely during test).
func (d *EthereumData) initDb(db *gorm.DB) error {
	if err := db.AutoMigrate(&ethereum.Block{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&ethereum.Log{}); err != nil {
		return err
	}
	return nil
}
