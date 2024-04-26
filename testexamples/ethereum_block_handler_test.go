package examples

import (
	"fmt"
	"testing"

	"github.com/Zettablock/zsource/dao/ethereum"
	"github.com/Zettablock/zsource/testutils"
	"github.com/Zettablock/zsource/utils"
)

// A simple handler that looks for block 2 and if found writes it to the
// destination table and custom table.
func FindBlockHandlerString(blockNumber string, deps *utils.Deps) (bool, error) {
	if blockNumber == "2" {
		deps.DestinationDB.Create(&ethereum.Block{Number: 2})
		deps.DestinationDB.Exec("INSERT INTO dest_init_example VALUES (1, 'test1')")
		return false, nil
	}
	return false, nil
}

// A simple handler that looks for block 3 and if found writes it to the
// destination table and custom table.
func FindBlockHandlerInt64(blockNumber int64, deps *utils.Deps) (bool, error) {
	if blockNumber == 3 {
		deps.DestinationDB.Create(&ethereum.Block{Number: 3})
		deps.DestinationDB.Exec("INSERT INTO dest_init_example VALUES (2, 'test2')")
		return false, nil
	}
	return false, nil
}

// Example of using the EthereumBlockHandlerTestRunner to test block handlers.
func TestHandlers(t *testing.T) {
	sourceData := []*ethereum.Block{
		{Number: 1},
		{Number: 2},
		{Number: 3},
	}
	runner := testutils.NewEthereumBlockHandlerTestRunner(t, sourceData, "dest_init_example.sql")
	defer runner.Close()

	// Returns a checker function that checks the table in destination database
	// has expected count of rows.
	checkerMaker := func(rowCountExp int64) testutils.DepsChecker {
		return func(deps *utils.Deps) error {
			// Verify row count in ethereum.blocks.
			var rowCountActual int64
			deps.DestinationDB.Table("ethereum.blocks").Count(&rowCountActual)
			if rowCountActual != rowCountExp {
				return fmt.Errorf("expected %d rows, got %d", rowCountExp, rowCountActual)
			}

			// Verify row count in custom table.
			var customTableRowCountActual int64
			deps.DestinationDB.Raw("SELECT COUNT(*) FROM dest_init_example").Scan(&customTableRowCountActual)
			if customTableRowCountActual != rowCountExp {
				return fmt.Errorf("expected %d rows in dest_init_example, got %d", rowCountExp, customTableRowCountActual)
			}
			return nil
		}
	}

	runner.TestHandlerString(FindBlockHandlerString, checkerMaker(1))
	runner.TestHandlerInt64(FindBlockHandlerInt64, checkerMaker(2))
}
