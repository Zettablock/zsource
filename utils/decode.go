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

// DecodeLog decodes a raw log and generate a dao.Event
func (d *Decoder) DecodeLog(vLog types.Log, rawLog *ethereum.Log) (dao.Event, error) {
	boundContract := bind.NewBoundContract(common.Address{}, d.ABI, nil, nil, nil)
	events := d.ABI.Events
	evt := dao.Event{}

	for name, event := range events {
		if event.ID.Hex() == vLog.Topics[0].Hex() {
			rawLog.Event = name
			rawLog.EventSignature = event.Sig

			err := boundContract.UnpackLogIntoMap(evt, name, vLog)
			if err != nil {
				return nil, err
			}
			inputs := event.Inputs
			for _, input := range inputs {
				rawLog.ArgumentNames = append(rawLog.ArgumentNames, input.Name)
				rawLog.ArgumentTypes = append(rawLog.ArgumentTypes, input.Type.String())
				rawLog.ArgumentValues = append(rawLog.ArgumentValues, fmt.Sprint(evt[input.Name]))
			}
			rawLog.DecodedFromAbi = true
			return evt, nil
		}
	}
	return nil, nil
}
