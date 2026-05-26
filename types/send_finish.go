package types

type SendFinishedMsg struct {
	TxHash string
	Draft  TxDraft
	Err    error
}
