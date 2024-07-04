package utils

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/Zettablock/zsource/dao/ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	HexHeader     = "0x"
	DelimiterChar = "\x1f"
)

var StopChars = []string{"\n", "\r", DelimiterChar}

// ToRawLog converts a types.Log from RPC response to an ethereum.Log
func ToRawLog(vLog types.Log, blockTime time.Time, blockDate time.Time) *ethereum.Log {
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
		BlockTime:        blockTime,
		ProcessTime:      time.Now(),
		BlockDate:        blockDate,
	}

	topics := make([]string, len(vLog.Topics))

	for i, topic := range vLog.Topics {
		topics[i] = topic.Hex()
	}
	rawLog.Topics = topics

	return rawLog
}

func ToBlock(b *types.Block) *ethereum.Block {
	difficulty, _ := b.Difficulty().Float64()
	uncles := b.Uncles()
	ul := make([]string, len(uncles))
	for i, u := range uncles {
		ul[i] = u.Hash().Hex()
	}

	extraDataRaw := hexutil.Encode(b.Extra())
	extraData := ToHumanReadableField(extraDataRaw)

	t := time.Unix(int64(b.Header().Time), 0)
	date := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	block := &ethereum.Block{
		Number:            int64(b.Number().Uint64()),
		Hash:              b.Hash().Hex(),
		ParentHash:        b.ParentHash().Hex(),
		Nonce:             fmt.Sprintf("%x", b.Nonce()),
		MixHash:           b.MixDigest().Hex(),
		Sha3Uncles:        b.UncleHash().Hex(),
		LogsBloom:         hexutil.Encode(b.Bloom().Bytes()),
		TransactionsRoot:  b.Root().Hex(),
		StateRoot:         b.Root().Hex(),
		ReceiptsRoot:      b.ReceiptHash().Hex(),
		Miner:             b.Coinbase().Hex(),
		Difficulty:        difficulty,
		TotalDifficulty:   0,
		Size:              int64(b.Size()),
		GasLimit:          int64(b.GasLimit()),
		GasUsed:           int64(b.GasUsed()),
		BaseFeePerGas:     b.BaseFee().Int64(),
		Timestamp:         t,
		Uncles:            ul,
		NumOfTransactions: int32(len(b.Transactions())),
		ExtraDataRaw:      extraDataRaw,
		ExtraData:         extraData,
		ProcessTime:       time.Now(),
		BlockDate:         date,
	}

	return block
}

func ToHumanReadableField(text string) string {
	if strings.HasPrefix(text, HexHeader) {
		text = strings.TrimPrefix(text, HexHeader)
	}
	decoded, err := hex.DecodeString(text)
	if err != nil {
		return ""
	}
	ret := string(decoded)
	for _, c := range StopChars {
		ret = strings.ReplaceAll(ret, c, "")
	}
	return ret
}
