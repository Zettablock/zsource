package examples

import (
	"testing"

	"github.com/Zettablock/zsource/dao/ethereum"
	"github.com/Zettablock/zsource/testutils"
	"github.com/Zettablock/zsource/utils"
)

func FindBlockHandlerString(blockNumber string, deps *utils.Deps) (bool, error) {
	if blockNumber == "2" {
		// TODO(meng): write to destination db.
		return false, nil
	}
	return false, nil
}

func FindBlockHandlerInt64(blockNumber int64, deps *utils.Deps) (bool, error) {
	if blockNumber == 2 {
		// TODO(meng): write to destination db.
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

	runner.TestHandlerString(FindBlockHandlerString, 0)
	runner.TestHandlerInt64(FindBlockHandlerInt64, 0)
}
