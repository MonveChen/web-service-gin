/*
 * @Author: Monve
 * @Date: 2023-07-24 18:38:25
 * @LastEditors: Monve
 * @LastEditTime: 2023-07-25 12:27:40
 * @FilePath: /web-service-gin/utils/chains/chains.go
 */
package chains

import (
	"fmt"
	token "web-service-gin/third_party/erc20"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TokenInfo struct {
	Symbol   string `json:"symbol" example:"PLTC"`
	Decimals uint8  `json:"decimals" example:"18"`
}

func EthContractInfo(address string) (*TokenInfo, error) {
	client, err := ethclient.Dial("https://eth.llamarpc.com")
	if err != nil {
		return nil, err
	}
	contract_address := common.HexToAddress(address)

	// use erc20.go
	contract, err := token.NewTokenCaller(contract_address, client)
	if err != nil {
		return nil, err
	}
	symbol, err := contract.Symbol(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	fmt.Println("symbol:", symbol)

	decimals, err := contract.Decimals(nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("decimals:", decimals)

	return &TokenInfo{
		Symbol:   symbol,
		Decimals: decimals,
	}, nil
}
