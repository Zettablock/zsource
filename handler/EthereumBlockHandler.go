package handler

import "github.com/Zettablock/zsource/utils"

// Custom ethereum block handler should implement this interface.
type EthereumBlockHandler interface {
	Handle(blockNumber string, deps *utils.Deps) (bool, error)
}
