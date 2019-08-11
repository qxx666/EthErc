package utils

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"math"
	"math/big"
	"strings"
)

func TransferEth(
	client *ethclient.Client,
	keystoreStr string,
	password string,
	amount *big.Int,
	toAddress string,
	gasO *big.Int,
	gasPriceO *big.Int,
) (string, error) {

	auth, err := bind.NewTransactor(strings.NewReader(keystoreStr), password)
	if err != nil {
		return "", err
	}

	nonce, err := client.PendingNonceAt(context.TODO(), auth.From)
	if err != nil {
		return "", err
	}

	key, err := keystore.DecryptKey([]byte(keystoreStr), password)

	if err != nil {
		return "", err
	}

	var gasPrice *big.Int
	//如果有设定gasprice
	if gasPriceO == nil {
		gasPrice, err = client.SuggestGasPrice(context.TODO())

		if err != nil {
			return "", err
		}
	} else {
		gasPrice = gasPriceO
	}

	var gas *big.Int
	if gasO == nil {
		gas = new(big.Int).SetUint64(params.TxGas)
	} else {
		gas = gasO
	}

	data := common.FromHex("0x")

	rawTx := types.NewTransaction(
		nonce,
		common.HexToAddress(toAddress),
		amount,
		gas,
		gasPrice,
		data)
	signer := types.NewEIP155Signer(params.MainnetChainConfig.ChainId)

	//signer := types.HomesteadSigner{}
	txs, err := types.SignTx(rawTx, signer, key.PrivateKey)

	if err != nil {
		return "", err
	}

	err = client.SendTransaction(context.Background(), txs)
	if err != nil {
		return "", err
	}

	return txs.Hash().Hex(), nil
}

func TransferToken(
	client *ethclient.Client,
	keystoreStr string,
	password string,
	amount *big.Float,
	contractAddress string,
	toAddress string,
) (string, error) {
	auth, err := bind.NewTransactor(strings.NewReader(keystoreStr), password)
	if err != nil {
		return "", err
	}

	token, err := NewToken(common.HexToAddress(contractAddress), client)
	if err != nil {
		return "", err
	}

	decimal, err := token.Decimals(nil)
	if err != nil {
		return "", err
	}

	tenDecimal := big.NewFloat(math.Pow(10, float64(decimal)))
	convertAmount, _ := new(big.Float).Mul(tenDecimal, amount).Int(&big.Int{})

	txs, err := token.Transfer(auth, common.HexToAddress(toAddress), convertAmount)

	if err != nil {
		return "", err
	}
	return txs.Hash().Hex(), nil

}

func TransferTokenOrigin(
	client *ethclient.Client,
	keystoreStr string,
	password string,
	amount *big.Int,
	contractAddress string,
	toAddress string,
	gasPrice *big.Int,
) (string, error) {
	auth, err := bind.NewTransactor(strings.NewReader(keystoreStr), password)

	if err != nil {
		return "", err
	}

	auth.GasPrice = gasPrice

	token, err := NewTokenTransactor(common.HexToAddress(contractAddress), client)
	if err != nil {
		return "", err
	}

	txs, err := token.Transfer(auth, common.HexToAddress(toAddress), amount)

	if err != nil {
		return "", err
	}
	return txs.Hash().Hex(), nil
}
