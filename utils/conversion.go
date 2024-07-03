package utils

import (
	"github.com/Zettablock/zsource/dao/ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

// ConvertToRawLog converts a types.Log from RPC response to an ethereum.Log
func ConvertToRawLog(vLog types.Log) *ethereum.Log {
	rawLog := &ethereum.Log{
		BlockNumber:      int64(vLog.BlockNumber),
		BlockHash:        vLog.BlockHash.Hex(),
		TransactionHash:  vLog.TxHash.Hex(),
		TransactionIndex: int32(vLog.TxIndex),
		LogIndex:         int32(vLog.Index),
		DecodedFromAbi:   false,
		ContractAddress:  vLog.Address.Hex(),
		Removed:          vLog.Removed,
		Data:             hexutil.Encode(vLog.Data),
	}

	for _, topic := range vLog.Topics {
		rawLog.Topics = append(rawLog.Topics, topic.Hex())
	}

	return rawLog
}
