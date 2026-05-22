package enums

type Screen int

const (
	ScreenLoading Screen = iota
	ScreenDashboard
	ScreenConfirm
	ScreenSending
	ScreenSent
	ScreenError
)
