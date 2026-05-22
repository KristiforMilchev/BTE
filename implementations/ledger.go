package implementations

import (
	"bos/constants"
	"bos/interfaces"
	"bos/utils"
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/usbwallet"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Ledger struct {
	networkProvider interfaces.INetwork
}

func (l *Ledger) Account() (*accounts.Account, error) {
	wallet, account, err := openLedger(false)
	if err != nil {
		log.Printf("Failed to connect to ledger -> %s", err)
		return nil, err
	}
	defer wallet.Close()
	return &account, nil
}

func (l *Ledger) Open() (*accounts.Wallet, *accounts.Account, error) {
	wallet, account, err := openLedger(true)
	if err != nil {
		log.Printf("Failed to connect to ledger -> %s", err)
		return nil, nil, err
	}
	return &wallet, &account, nil
}

func (l *Ledger) SendTransaction(receiver string, amount *string, overrideLimit *uint64) (*string, error) {
	client, chain, ctx, cancelationToken, err := l.networkProvider.Active()
	if err != nil {
		log.Printf("Failed to initiate a netowrk provider aborting transaction -> %s", err)
		return nil, err
	}
	defer cancelationToken()

	wallet, account, err := openLedger(true)
	if err != nil {
		log.Printf("Failed to connect to ledger -> %s", err)
		return nil, err
	}
	defer wallet.Close()

	value, err := utils.ParseEtherToWei(*amount)
	if err != nil {
		log.Printf("Failed to estimate amount in wei -> %s", err)
		return nil, err
	}

	to := common.HexToAddress(receiver)

	nonce, err := client.PendingNonceAt(ctx, account.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price: %w", err)
	}

	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
		From:  account.Address,
		To:    &to,
		Value: value,
	})
	if err != nil {
		gasLimit = 21000
	}

	tx := types.NewTransaction(
		nonce,
		to,
		value,
		gasLimit,
		gasPrice,
		nil,
	)

	signedTx, err := wallet.SignTx(account, tx, chain)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction with Ledger: %w", err)
	}

	if err := client.SendTransaction(ctx, signedTx); err != nil {
		return nil, fmt.Errorf("failed to broadcast transaction: %w", err)
	}

	txHash := signedTx.Hash().Hex()
	return &txHash, nil
}

func openLedger(pin bool) (accounts.Wallet, accounts.Account, error) {
	hub, err := usbwallet.NewLedgerHub()
	if err != nil {
		return nil, accounts.Account{}, fmt.Errorf("failed to create Ledger hub: %w", err)
	}

	wallets := hub.Wallets()
	if len(wallets) == 0 {
		return nil, accounts.Account{}, errors.New("no Ledger device found")
	}

	wallet := wallets[0]
	if err := wallet.Open(""); err != nil {
		return nil, accounts.Account{}, fmt.Errorf("failed to open Ledger: %w", err)
	}

	status, err := wallet.Status()
	if err != nil {
		_ = wallet.Close()
		return nil, accounts.Account{}, fmt.Errorf("Ledger status error: %w", err)
	}

	path, err := accounts.ParseDerivationPath(constants.DerivationPath)
	if err != nil {
		_ = wallet.Close()
		return nil, accounts.Account{}, fmt.Errorf("failed to parse derivation path: %w", err)
	}

	account, err := wallet.Derive(path, pin)
	if err != nil {
		_ = wallet.Close()
		return nil, accounts.Account{}, fmt.Errorf("Ledger status: %s | failed to derive address: %w", status, err)
	}

	return wallet, account, nil
}

func (l *Ledger) Address() (common.Address, error) {
	account, err := l.Account()
	if err != nil {
		return common.Address{}, err
	}

	return account.Address, nil
}

func (l *Ledger) SignTransaction(
	ctx context.Context,
	tx *types.Transaction,
	chainID *big.Int,
) (*types.Transaction, error) {
	wallet, account, err := l.Open()
	if err != nil {
		return nil, err
	}
	defer (*wallet).Close()

	return (*wallet).SignTx(*account, tx, chainID)
}

func NewLedger(networkProvider interfaces.INetwork) *Ledger {

	return &Ledger{
		networkProvider: networkProvider,
	}
}
