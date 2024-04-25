package examples

import (
	"fmt"
	"testing"

	"github.com/Zettablock/zsource/dao/ethereum"
	"github.com/Zettablock/zsource/testutils"
	"github.com/Zettablock/zsource/utils"
)

func FindBlockHandlerString(blockNumber string, deps *utils.Deps) (bool, error) {
	if blockNumber == "2" {
		deps.DestinationDB.Create(&ethereum.Block{Number: 2})
		return false, nil
	}
	return false, nil
}

func FindBlockHandlerInt64(blockNumber int64, deps *utils.Deps) (bool, error) {
	if blockNumber == 3 {
		deps.DestinationDB.Create(&ethereum.Block{Number: 3})
		return false, nil
	}
	return false, nil
}

func TestHandlers(t *testing.T) {
	sourceData := []*ethereum.Block{
		{Number: 1},
		{Number: 2},
		{Number: 3},
	}
	runner := testutils.NewEthereumBlockHandlerTestRunner(t, sourceData)
	defer runner.Close()

	// Returns a checker function that checks the table in destination database
	// has expected count of rows.
	checkerMaker := func(rowCountExp int64) testutils.DepsChecker {
		return func(deps *utils.Deps) error {
			var rowCountActual int64
			deps.DestinationDB.Table("ethereum.blocks").Count(&rowCountActual)
			if rowCountActual != rowCountExp {
				return fmt.Errorf("expected %d rows, got %d", rowCountExp, rowCountActual)
			}
			return nil
		}
	}

	runner.TestHandlerString(FindBlockHandlerString, checkerMaker(1))
	runner.TestHandlerInt64(FindBlockHandlerInt64, checkerMaker(2))
}
