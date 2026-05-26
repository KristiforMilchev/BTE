package repositories

import "bos/interfaces"

type RepositoryRegister struct {
	Accounts     AccountsRepository
	Contacts     ContactsRepository
	Network      NetworkRepository
	Transactions TransactionsRepository
}

func NewRegister(storage interfaces.IStorage) RepositoryRegister {
	return RepositoryRegister{
		Accounts:     NewAccountsRepository(storage),
		Contacts:     NewContactsRepository(storage),
		Network:      NewNetworkRepository(storage),
		Transactions: NewTransactionsRepository(storage),
	}
}
