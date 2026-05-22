package types

type SendFinishedMsg struct {
	TxHash string
	Err    error
}
