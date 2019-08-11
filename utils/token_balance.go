package utils

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math"
	"math/big"
)

var (
	tokens map[string]*TokenCaller = make(map[string]*TokenCaller)
)

func GetTokenBalance(client *ethclient.Client, contractAddress string, address string) (balance *big.Float, err error) {

	token := tokens[contractAddress]

	if token == nil {
		token, err = NewTokenCaller(common.HexToAddress(contractAddress), client)

		if err != nil {
			return nil, err
		}
	}

	decimal, err := token.Decimals(nil)

	if err != nil {
		return nil, err
	}

	balancex, err := token.BalanceOf(nil, common.HexToAddress(address))

	if err != nil {
		return nil, err
	}

	balanceFloat := new(big.Float).SetInt(balancex)
	balanceFloat = new(big.Float).Mul(balanceFloat, big.NewFloat(math.Pow(10, -float64(decimal))))
	if err != nil {
		return nil, err
	}
	return balanceFloat, nil
}

func GetMutilTokenBalance(client *ethclient.Client, contractAddress string, address string) (balance *big.Int, balanceF *big.Float, err error) {

	token := tokens[contractAddress]

	if token == nil {
		token, err = NewTokenCaller(common.HexToAddress(contractAddress), client)

		if err != nil {
			return nil, nil, err
		}
	}

	decimal, err := token.Decimals(nil)

	if err != nil {
		return nil, nil, err
	}

	balancex, err := token.BalanceOf(nil, common.HexToAddress(address))

	if err != nil {
		return nil, nil, err
	}

	balanceFloat := new(big.Float).SetInt(balancex)
	balanceFloat = new(big.Float).Mul(balanceFloat, big.NewFloat(math.Pow(10, -float64(decimal))))
	if err != nil {
		return nil, nil, err
	}
	return balancex, balanceFloat, nil
}

func GetEthBalance(client *ethclient.Client, address string) (balance *big.Float, err error) {
	balancex, err := client.BalanceAt(
		context.TODO(),
		common.HexToAddress(address), nil)

	if err != nil {
		return nil, err
	}
	balanceFloat := new(big.Float).SetInt(balancex)
	balanceFloat = new(big.Float).Mul(balanceFloat, big.NewFloat(math.Pow(10, -float64(18))))

	return balanceFloat, nil
}

func GetMutilEthBalance(client *ethclient.Client, address string) (balance *big.Int, balanceF *big.Float, err error) {
	balance, err = client.BalanceAt(
		context.TODO(),
		common.HexToAddress(address), nil)

	if err != nil {
		return nil, nil, err
	}
	balanceFloat := new(big.Float).SetInt(balance)
	balanceFloat = new(big.Float).Mul(balanceFloat, big.NewFloat(math.Pow(10, -float64(18))))

	return balance, balanceFloat, nil
}
