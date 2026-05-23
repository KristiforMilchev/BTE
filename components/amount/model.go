package amount

import (
	"bos/types"

	"github.com/charmbracelet/bubbles/textinput"
)

type Model struct {
	token       types.Token
	amountInput textinput.Model
	active      bool
}

func (m *Model) Blur() {
	m.amountInput.Blur()
}

func (m *Model) Focus() {
	m.amountInput.Focus()
}

func (m *Model) SetSymbol(symbol types.Token) {
	m.token = symbol
}

func (m *Model) Value() string {
	return m.amountInput.Value()
}

func New() *Model {

	amount := textinput.New()
	amount.Placeholder = "0.01"
	amount.CharLimit = 32
	amount.Width = 0
	amount.Focus()
	amount.SetValue("12")

	return &Model{
		active:      false,
		amountInput: amount,
	}
}
