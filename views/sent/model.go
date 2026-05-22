package sent

import "bos/types"

type Model struct {
	payload types.SentPayload
	width   int
	height  int
}

func New(payload types.SentPayload) *Model { return &Model{payload: payload} }
