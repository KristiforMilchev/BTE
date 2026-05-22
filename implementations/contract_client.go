// implementations/contract_client.go
package implementations

import (
	"context"
	"math/big"

	"bos/interfaces"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ContractClient struct {
	networkProvider interfaces.INetwork
	signer          interfaces.ITransactionSigner
}

func (c *ContractClient) From() (common.Address, error) {
	return c.signer.Address()
}

func (c *ContractClient) Call(
	ctx context.Context,
	contract common.Address,
	parsed abi.ABI,
	method string,
	args ...any,
) ([]any, error) {
	client, _, _, cancel, err := c.networkProvider.Active()
	if err != nil {
		return nil, err
	}
	defer cancel()

	data, err := parsed.Pack(method, args...)
	if err != nil {
		return nil, err
	}

	raw, err := client.CallContract(ctx, ethereum.CallMsg{
		To:   &contract,
		Data: data,
	}, nil)
	if err != nil {
		return nil, err
	}

	return parsed.Unpack(method, raw)
}

func (c *ContractClient) Transact(
	ctx context.Context,
	contract common.Address,
	parsed abi.ABI,
	method string,
	value *big.Int,
	args ...any,
) (common.Hash, error) {
	client, chainID, networkCtx, cancel, err := c.networkProvider.Active()
	if err != nil {
		return common.Hash{}, err
	}
	defer cancel()

	from, err := c.signer.Address()
	if err != nil {
		return common.Hash{}, err
	}

	data, err := parsed.Pack(method, args...)
	if err != nil {
		return common.Hash{}, err
	}

	if value == nil {
		value = big.NewInt(0)
	}

	nonce, err := client.PendingNonceAt(networkCtx, from)
	if err != nil {
		return common.Hash{}, err
	}

	gasPrice, err := client.SuggestGasPrice(networkCtx)
	if err != nil {
		return common.Hash{}, err
	}

	gasLimit, err := client.EstimateGas(networkCtx, ethereum.CallMsg{
		From:  from,
		To:    &contract,
		Value: value,
		Data:  data,
	})
	if err != nil {
		return common.Hash{}, err
	}

	tx := types.NewTransaction(
		nonce,
		contract,
		value,
		gasLimit,
		gasPrice,
		data,
	)

	signedTx, err := c.signer.SignTransaction(ctx, tx, chainID)
	if err != nil {
		return common.Hash{}, err
	}

	if err := client.SendTransaction(networkCtx, signedTx); err != nil {
		return common.Hash{}, err
	}

	return signedTx.Hash(), nil
}

func NewContractClient(
	networkProvider interfaces.INetwork,
	signer interfaces.ITransactionSigner,
) *ContractClient {
	return &ContractClient{
		networkProvider: networkProvider,
		signer:          signer,
	}
}
