package errorview

import (
	"bos/types"
)

type Model struct {
	payload types.ErrorPayload
	width   int
	height  int
}

func New(payload types.ErrorPayload) *Model { return &Model{payload: payload} }
