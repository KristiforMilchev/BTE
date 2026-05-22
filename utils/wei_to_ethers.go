package utils

import "math/big"

func WeiToEther(wei *big.Int) string {
	ethValue := new(big.Float).Quo(
		new(big.Float).SetInt(wei),
		big.NewFloat(1e18),
	)

	return ethValue.Text('f', 8)
}
