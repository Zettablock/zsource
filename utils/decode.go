package utils

import (
	"fmt"

	dao "github.com/Zettablock/zsource/dao/common"
	"github.com/Zettablock/zsource/dao/ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Decoder struct {
	ABI abi.ABI
}

func NewDecoder(abi abi.ABI) *Decoder {
	return &Decoder{ABI: abi}
}

func (d *Decoder) DecodeLog(vLog types.Log, evtName string) (*ethereum.Log, dao.Event, error) {
	decodedLog := &ethereum.Log{}
	boundContract := bind.NewBoundContract(common.Address{}, d.ABI, nil, nil, nil)
	events := d.ABI.Events
	evt := dao.Event{}

	for name, event := range events {
		if event.ID.Hex() == vLog.Topics[0].Hex() && event.Name == evtName {
			decodedLog.BlockNumber = int64(vLog.BlockNumber)
			decodedLog.BlockHash = vLog.BlockHash.Hex()
			decodedLog.TransactionHash = vLog.TxHash.Hex()
			decodedLog.TransactionIndex = int32(vLog.TxIndex)
			decodedLog.LogIndex = int32(vLog.Index)
			decodedLog.Event = name
			decodedLog.EventSignature = event.Sig
			decodedLog.DecodedFromAbi = true
			decodedLog.ContractAddress = vLog.Address.Hex()

			for _, topic := range vLog.Topics {
				decodedLog.Topics = append(decodedLog.Topics, topic.Hex())
			}

			err := boundContract.UnpackLogIntoMap(evt, name, vLog)
			if err != nil {
				return nil, nil, err
			}
			inputs := event.Inputs
			for _, input := range inputs {
				decodedLog.ArgumentNames = append(decodedLog.ArgumentNames, input.Name)
				decodedLog.ArgumentTypes = append(decodedLog.ArgumentTypes, input.Type.String())
				decodedLog.ArgumentValues = append(decodedLog.ArgumentValues, fmt.Sprint(evt[input.Name]))
			}
			return decodedLog, evt, nil
		}
	}
	return nil, nil, nil
}
