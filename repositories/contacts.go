package repositories

import (
	"bos/interfaces"
	"bos/types"
	"context"
	"log"

	"github.com/google/uuid"
)

type ContactsRepository struct {
	storage interfaces.IStorage
}

func (c ContactsRepository) Create(contact types.Contact) error {
	_, err := c.storage.Exec(context.Background(), `
		INSERT INTO contacts (id, name, address)
		VALUES (?, ?, ?)
		ON CONFLICT(address) DO UPDATE SET
			name = excluded.name;
	`, uuid.NewString(), contact.Name, contact.Address)
	if err != nil {
		log.Printf("Failed to save contact -> %s", err)
		return err
	}
	return nil
}

func (c ContactsRepository) Contacts() (*[]types.Contact, error) {
	rows, err := c.storage.Query(context.Background(), `
		SELECT name, address FROM contacts
		ORDER BY name COLLATE NOCASE;
	`)
	if err != nil {
		log.Printf("Failed to fetch contacts -> %s", err)
		return nil, err
	}
	defer rows.Close()

	contacts := []types.Contact{}
	for rows.Next() {
		var contact types.Contact
		if err := rows.Scan(&contact.Name, &contact.Address); err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &contacts, nil
}

func NewContactsRepository(storage interfaces.IStorage) ContactsRepository {
	return ContactsRepository{
		storage: storage,
	}
}
