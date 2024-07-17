package handler

import (
	"github.com/Zettablock/zsource/dao/ethereum"
	"github.com/Zettablock/zsource/utils"
)

// Example for custom Ethereum log handler.
type ExampleEthereumLogHandler struct {
}

// Verify interface compliance.
var _ EthereumLogHandler = (*ExampleEthereumLogHandler)(nil)

func (h *ExampleEthereumLogHandler) Handle(log ethereum.Log, deps *utils.Deps) (bool, error) {
	return false, nil
}

// Example for custom Ethereum block handler.
type ExampleEthereumBlockHandler struct {
}

// Verify interface compliance.
var _ EthereumBlockHandler = (*ExampleEthereumBlockHandler)(nil)

func (h *ExampleEthereumBlockHandler) Handle(blockNumber string, deps *utils.Deps) (bool, error) {
	return false, nil
}
