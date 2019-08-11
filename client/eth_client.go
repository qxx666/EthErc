package client

import (
	"EthErc/models"
	"github.com/astaxie/beego"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

var instance *ethclient.Client

func EthClient() *ethclient.Client {
	if instance != nil {
		return instance
	} else {
		client, err := connectToRpc()
		if err != nil {
			panic("RPC连接失败")
			beego.Error(err.Error())
			return nil
		}
		instance = client
		return client
	}
}

func connectToRpc() (*ethclient.Client, error) {
	client, err := rpc.Dial(models.SysSetting().RpcHost)

	if err != nil {
		beego.Error(err.Error())
		return nil, err
	}
	conn := ethclient.NewClient(client)

	return conn, nil
}
