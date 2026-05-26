package repositories

import "bos/interfaces"

type RepositoryRegister struct {
	Accounts     AccountsRepository
	Network      NetworkRepository
	Transactions TransactionsRepository
}

func NewRegister(storage interfaces.IStorage) RepositoryRegister {
	return RepositoryRegister{
		Accounts:     NewAccountsRepository(storage),
		Network:      NewNetworkRepository(storage),
		Transactions: NewTransactionsRepository(storage),
	}
}
