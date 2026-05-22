package repositories

import (
	"bos/interfaces"
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type AccountsRepository struct {
	storage interfaces.IStorage
}

func (a *AccountsRepository) Account() (*common.Address, error) {
	sqlQuery := `
		SELECT address FROM accounts
		WHERE last_used = 1;
	`

	query := a.storage.QueryRow(context.Background(), sqlQuery)
	var account string
	err := query.Scan(&account)
	if err == sql.ErrNoRows {
		return nil, errors.New("not connected")
	}

	if err != nil {
		log.Printf("Failed to retrive last used account error communicating with database -> %s", err)
		return nil, err
	}

	address := common.HexToAddress(account)
	return &address, nil
}

func (a *AccountsRepository) AddAccount(address common.Address) error {
	ctx := context.Background()

	tx, err := a.storage.DB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		UPDATE accounts
		SET last_used = 0;
	`)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO accounts (id, address, last_used)
		VALUES (?, ?, 1)
		ON CONFLICT(address) DO UPDATE SET
			last_used = 1;
	`, uuid.NewString(), address.Hex())
	if err != nil {
		return err
	}

	return tx.Commit()
}

func NewAccountsRepository(storage interfaces.IStorage) AccountsRepository {
	return AccountsRepository{
		storage: storage,
	}
}
