package types

import "bos/enums"

type NavigateMsg struct {
	Screen  enums.Screen
	Payload any
}
