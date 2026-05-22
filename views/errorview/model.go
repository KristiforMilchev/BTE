package errorview

import "bos/views"

type Model struct {
	payload views.ErrorPayload
	width   int
	height  int
}

func New(payload views.ErrorPayload) *Model { return &Model{payload: payload} }
