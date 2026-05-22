package sent

import "bos/views"

type Model struct {
	payload views.SentPayload
	width   int
	height  int
}

func New(payload views.SentPayload) *Model { return &Model{payload: payload} }
