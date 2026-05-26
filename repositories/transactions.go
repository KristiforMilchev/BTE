package repositories

import (
	"bos/interfaces"
	"bos/types"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type TransactionsRepository struct {
	storage interfaces.IStorage
}

func (t *TransactionsRepository) SaveNativeTransfer(draft types.TxDraft, txHash string, network types.Network) error {
	if txHash == "" {
		return errors.New("transaction hash is required")
	}
	if draft.RecipientAddress == "" {
		return errors.New("recipient address is required")
	}
	if network.Id == uuid.Nil {
		return errors.New("network id is required")
	}

	ctx := context.Background()
	tx, err := t.storage.DB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	accountID, err := accountID(ctx, tx, draft.FromAddress)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO contact_transactions (id, tx_hash, recipient, token, amount, account_id, network_id)
		VALUES (?, ?, ?, NULL, ?, ?, ?)
		ON CONFLICT(tx_hash) DO UPDATE SET
			recipient = excluded.recipient,
			amount = excluded.amount,
			account_id = excluded.account_id,
			network_id = excluded.network_id,
			token = excluded.token;
	`, uuid.NewString(), txHash, draft.RecipientAddress, draft.Amount, accountID, network.Id.String())
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (t *TransactionsRepository) GetTransactions(network *uuid.UUID, account *string) (*[]types.Transaction, error) {
	if network == nil || *network == uuid.Nil {
		return nil, errors.New("network id is required")
	}
	if account == nil || *account == "" {
		return nil, errors.New("account address is required")
	}

	ctx := context.Background()
	accountID, err := t.accountIDByAddress(ctx, *account)
	if err != nil {
		return nil, err
	}

	sqlQuery := `
		SELECT recipient, tx_hash, amount
		FROM (
			SELECT recipient, tx_hash, amount, 0 AS source_order
			FROM token_transactions
			WHERE network_id = ? AND account_id = ?

			UNION ALL

			SELECT recipient, tx_hash, amount, 1 AS source_order
			FROM contact_transactions
			WHERE network_id = ? AND account_id = ?
		)
		ORDER BY source_order
	`

	query, err := t.storage.Query(ctx, sqlQuery, network.String(), accountID, network.String(), accountID)
	if err == sql.ErrNoRows {
		log.Printf("Account has no transactions on the network -> %s | %s", *network, *account)
		return &[]types.Transaction{}, nil
	}

	if err != nil {
		log.Printf("Can't reach the database failing to retrive transactions -> %s", err)
		return nil, err
	}

	defer query.Close()

	var transactions []types.Transaction
	for query.Next() {
		var transaction types.Transaction
		err := query.Scan(&transaction.To, &transaction.TxHash, &transaction.Amount)
		if err != nil {
			log.Printf("Failed to map transaction data to transaction -> %s", err)
			return nil, err
		}

		transactions = append(transactions, transaction)
	}
	if err := query.Err(); err != nil {
		return nil, err
	}

	return &transactions, nil
}

func (t *TransactionsRepository) accountIDByAddress(ctx context.Context, address string) (string, error) {
	var id string
	err := t.storage.QueryRow(ctx, `
		SELECT id FROM accounts
		WHERE address = ?;
	`, address).Scan(&id)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("account %s is not saved", address)
	}
	if err != nil {
		return "", err
	}
	return id, nil
}

func accountID(ctx context.Context, tx *sql.Tx, address string) (string, error) {
	var id string
	err := tx.QueryRowContext(ctx, `
		SELECT id FROM accounts
		WHERE address = ?;
	`, address).Scan(&id)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("account %s is not saved", address)
	}
	if err != nil {
		return "", err
	}
	return id, nil
}

func NewTransactionsRepository(storage interfaces.IStorage) TransactionsRepository {
	return TransactionsRepository{
		storage: storage,
	}
}
