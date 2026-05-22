package types

type WalletLoadedMsg struct {
	Address string
	Balance string
	ChainID string
	Err     error
}
