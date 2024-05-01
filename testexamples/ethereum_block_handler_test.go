package examples

import (
	"fmt"
	"testing"
	"time"

	"github.com/Zettablock/zsource/configs"
	"github.com/Zettablock/zsource/dao/ethereum"
	"github.com/Zettablock/zsource/testutils"
	"github.com/Zettablock/zsource/utils"
)

// A simple handler that looks for block 2 and if found writes it to the
// destination table and custom table.
func FindBlockHandlerString(blockNumber string, deps *utils.Deps) (bool, error) {
	if deps.Config.Environment != "unittest" {
		return false, nil
	}
	if blockNumber == "2" {
		// Exercise the code to read source logs table from the handler.
		var logs []*ethereum.Log
		deps.SourceDB.Table("ethereum.logs").Where("block_number = ?", 2).Find(&logs)
		if len(logs) == 0 {
			return false, nil
		}
		deps.DestinationDB.Table("ethereum.blocks").Create(&ethereum.Block{Number: 2})
		deps.DestinationDB.Exec("INSERT INTO dest_init_example VALUES (1, 'test1')")
		deps.Logger.Info("successfully processed block 2")
		return false, nil
	}
	return false, nil
}

// A simple handler that looks for block 3 and if found writes it to the
// destination table and custom table.
func FindBlockHandlerInt64(blockNumber int64, deps *utils.Deps) (bool, error) {
	if blockNumber == 3 {
		deps.DestinationDB.Table("ethereum_mainnet.blocks").Create(&ethereum.Block{Number: 3})
		deps.DestinationDB.Exec("INSERT INTO dest_init_example VALUES (2, 'test2')")
		return false, nil
	}
	return false, nil
}

func TimeHandlerInt64(blockNumber int64, deps *utils.Deps) (bool, error) {
	if blockNumber == 3 {
		var logs []*ethereum.Log
		deps.SourceDB.Raw("SELECT * FROM ethereum.logs").Scan(&logs)
		fmt.Printf("logs: %v\n", logs[0].BlockTime)
		return false, nil
	}
	return false, nil
}

// Example of using the EthereumBlockHandlerTestRunner to test block handlers.
func TestHandlers(t *testing.T) {
	// Prepare source schema and data.
	sourceData := testutils.NewEthereumData()
	blocks := []*ethereum.Block{
		{Number: 1},
		{Number: 2},
		{Number: 3},
	}
	logs := []*ethereum.Log{
		{TransactionHash: "tx1", BlockNumber: 2, LogIndex: 1, BlockTime: time.Date(2024, time.April, 16, 10, 30, 0, 0, time.UTC)},
		{TransactionHash: "tx1", BlockNumber: 2, LogIndex: 2, BlockTime: time.Date(2024, time.April, 16, 10, 30, 0, 0, time.UTC)},
	}
	sourceData.AddSchemaData("ethereum", blocks, logs)
	sourceData.AddSchemaData("ethereum_mainnet", blocks, logs)
	sourceData.AddSchemaData("ethereum_holesky", blocks, logs)

	// Prepare dest schema.
	destData := testutils.NewEthereumData()
	destData.AddSchemaDataEmpty("ethereum")
	destData.AddSchemaDataEmpty("ethereum_mainnet")
	destData.AddSchemaDataEmpty("ethereum_holesky")

	// Prepare pipeline config.
	config := &configs.PipelineConfig{Environment: "unittest"}

	runner := testutils.NewEthereumBlockHandlerTestRunner(t, config, sourceData, "dest_init_example.sql", destData)
	defer runner.Close()

	// Returns a checker function that checks the table in destination database
	// has expected count of rows.
	checkerMaker := func(schemaName string, rowCountExp int64) testutils.DepsChecker {
		return func(deps *utils.Deps) error {
			// Verify row count in ethereum.blocks.
			var rowCountActual int64
			deps.DestinationDB.Table(schemaName + ".blocks").Count(&rowCountActual)
			if rowCountActual != rowCountExp {
				return fmt.Errorf("expected %d rows, got %d", rowCountExp, rowCountActual)
			}
			return nil
		}
	}

	customTableCheckerMaker := func(rowCountExp int64) testutils.DepsChecker {
		return func(deps *utils.Deps) error {
			// Verify row count in custom table.
			var customTableRowCountActual int64
			deps.DestinationDB.Raw("SELECT COUNT(*) FROM dest_init_example").Scan(&customTableRowCountActual)
			if customTableRowCountActual != rowCountExp {
				return fmt.Errorf("expected %d rows in dest_init_example, got %d", rowCountExp, customTableRowCountActual)
			}
			return nil
		}
	}

	runner.TestHandlerString("ethereum", "ethereum_dest", FindBlockHandlerString, checkerMaker("ethereum", 1), customTableCheckerMaker(1))
	runner.TestHandlerInt64("ethereum_mainnet", "ethereum_dest", FindBlockHandlerInt64, checkerMaker("ethereum_mainnet", 1), customTableCheckerMaker(2))
	runner.TestHandlerInt64("ethereum", "ethereum", TimeHandlerInt64)
}
