package utils

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"strings"
)

type InputData struct {
	FromAddress common.Address
	ToAddress   common.Address
	Value       *big.Int
}

func DealInputData(inputData string) (*InputData, error) {

	transfer := "0xa9059cbb"
	transferFrom := "0x23b872dd"

	if strings.Contains(inputData, transfer) {

		//transfer
		inputData = strings.Replace(inputData, transfer, "", 1)
		splitM := len(inputData) / 2

		address := common.HexToAddress(inputData[0:splitM])

		value := inputData[splitM:]
		endZero := 0
		for i := 0; i < len(value); i++ {
			if value[i] != '0' {
				endZero = i
				break
			}
		}

		value = value[endZero:]

		valueBig, err := hexutil.DecodeBig("0x" + value)
		if err != nil {
			return nil, err
		}

		return &InputData{ToAddress: address, Value: valueBig}, nil
	} else if strings.Contains(inputData, transferFrom) {
		inputData = strings.Replace(inputData, transferFrom, "", 1)
		splitM := len(inputData) / 3
		fromAddress := common.HexToAddress(inputData[0:splitM])

		toAddress := common.HexToAddress(inputData[splitM : 2*splitM])

		value := inputData[2*splitM:]
		endZero := 0
		for i := 0; i < len(value); i++ {
			if value[i] != '0' {
				endZero = i
				break
			}
		}
		value = value[endZero:]
		valueBig, err := hexutil.DecodeBig("0x" + value)
		if err != nil {
			return nil, err
		}

		return &InputData{ToAddress: toAddress, FromAddress: fromAddress, Value: valueBig}, nil
	}
	return nil, errors.New("xxxx")
}
