package types

import "bos/enums"

type ErrorPayload struct {
	Title   string
	Message string
	Return  enums.Screen
}
