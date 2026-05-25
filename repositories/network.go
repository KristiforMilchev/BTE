package repositories

import (
	"bos/interfaces"
	"bos/types"
	"context"
	"database/sql"
	"errors"
	"log"
	"math/big"

	"github.com/google/uuid"
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

func (n *NetworkRepository) Networks() (*[]types.Network, error) {
	sqlQuery := `
		SELECT id, name, rpc, symbol, chain_id, block_explorer FROM networks
	`

	query, err := n.storage.Query(context.Background(), sqlQuery)
	if err == sql.ErrNoRows {
		log.Println("No networks saved in storage")
		return &[]types.Network{}, nil
	}

	if err != nil {
		log.Printf("Failed to fetch networks from storage aborting -> %s", err)
		return nil, err
	}

	defer query.Close()
	var networks []types.Network
	var chain *int64
	for query.Next() {
		var network types.Network
		err := query.Scan(&network.Id, &network.Name, &network.Rpc, &network.Symbol, &chain, &network.Explorer)
		if err != nil {
			log.Printf("Failed to map network data to entity while fetching all netowrks -> %s", err)
			return nil, err
		}
		network.Chain = big.NewInt(*chain)
		networks = append(networks, network)
	}

	return &networks, nil
}

func (n *NetworkRepository) Network(id *uuid.UUID) (*types.Network, error) {
	sqlQuery := `
		SELECT id, name, rpc, symbol, chain_id, block_explorer FROM networks
		WHERE id = ?
	`

	query := n.storage.QueryRow(context.Background(), sqlQuery, id)
	var network types.Network
	var chain *int64

	err := query.Scan(&network.Id, &network.Name, &network.Rpc, &network.Symbol, &chain, &network.Explorer)
	if err == sql.ErrNoRows {
		log.Println("No networks saved in storage")
		return &types.Network{}, nil
	}

	if err != nil {
		log.Printf("Failed to map network data to entity while fetching all netowrks -> %s", err)
		return nil, err
	}

	network.Chain = big.NewInt(*chain)
	return &network, nil
}

func (n *NetworkRepository) NetworkByRpc(rpc *string) (*types.Network, error) {
	sqlQuery := `
		SELECT id, name, rpc, symbol, chain_id, block_explorer FROM networks
		WHERE rpc = ?
	`

	query := n.storage.QueryRow(context.Background(), sqlQuery, rpc)
	var network types.Network
	var chain *int64

	err := query.Scan(&network.Id, &network.Name, &network.Rpc, &network.Symbol, &chain, &network.Explorer)
	if err == sql.ErrNoRows {
		log.Println("No networks saved in storage")
		return &types.Network{}, nil
	}

	if err != nil {
		log.Printf("Failed to map network data to entity while fetching all netowrks -> %s", err)
		return nil, err
	}

	network.Chain = big.NewInt(*chain)
	return &network, nil
}

func (n *NetworkRepository) Create(name *string, rpc *string, symbol *string, chainId *int64, blockExplorer *string) error {
	sqlQeury := `
		INSERT INTO networks
		(id, name, rpc, symbol, chain_id,block_explorer)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	if blockExplorer == nil {
		*blockExplorer = "-"
	}

	id := uuid.New()
	query, err := n.storage.Exec(context.Background(), sqlQeury, id, &name, &rpc, &symbol, &chainId, &blockExplorer)
	if err != nil {
		log.Printf("Failed to save network aborting -> %s", err)
		return err
	}

	count, err := query.RowsAffected()
	if err != nil {
		log.Printf("Failed to verify inserted rows after saving network assuming failed -> %s", err)
		return err
	}

	if count < 1 {
		log.Println("Failed to save network")
		return errors.New("Rows not saved when creating a network")
	}

	return nil
}

func NewNetworkRepository(storage interfaces.IStorage) NetworkRepository {
	return NetworkRepository{
		storage: storage,
	}
}
