package client

import (
	"EthErc/models"
	"github.com/onrik/ethrpc"
)

var (
	clientOther *ethrpc.EthRPC
)

func EthClientOther() *ethrpc.EthRPC {

	if clientOther == nil {
		conn := ethrpc.NewEthRPC(models.SysSetting().RpcHost)
		clientOther = conn
	}

	return clientOther
}

func GetCurrentBlockNumber() (int, error) {
	blockNumber, err := EthClientOther().EthBlockNumber()
	if err != nil {
		return 0, err
	}

	return blockNumber, nil
}
