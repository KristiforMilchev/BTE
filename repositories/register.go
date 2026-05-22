package repositories

import "bos/interfaces"

type RepositoryRegister struct {
	Accounts AccountsRepository
}

func NewRegister(storage interfaces.IStorage) RepositoryRegister {
	return RepositoryRegister{
		Accounts: NewAccountsRepository(storage),
	}
}
