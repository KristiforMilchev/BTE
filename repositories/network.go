package repositories

import (
	"bos/interfaces"
	"context"
	"log"
)

type NetworkRepository struct {
	storage interfaces.IStorage
}

func (n *NetworkRepository) Name(rpc *string) (*string, error) {
	sqlQyery := `
		SELECT name FROM networks
		WHERE rpc = ?
	`

	query := n.storage.QueryRow(context.Background(), sqlQyery, rpc)
	var name string
	err := query.Scan(&name)
	if err != nil {
		log.Printf("Can't retrive network name -> %s", err)
		return nil, err
	}

	return &name, nil
}

func NewNetworkRepository(storage interfaces.IStorage) NetworkRepository {
	return NetworkRepository{
		storage: storage,
	}
}
