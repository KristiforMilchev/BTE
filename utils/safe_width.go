package utils

func SafeWidth(width int) int {
	if width < 100 {
		return 100
	}
	return width
}
