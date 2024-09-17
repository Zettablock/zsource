package handler

import (
	"github.com/Zettablock/zsource/dao/ethereum"
	"github.com/Zettablock/zsource/utils"
)

// Custom ethereum log handler should implement this interface.
type EthereumLogHandler interface {
	Handle(log ethereum.Log, deps *utils.Deps) (bool, error)
}
