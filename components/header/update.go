package header

import "bos/enums"

func activeHelp(focus enums.FocusArea) string {
	switch focus {
	case enums.FocusAmount:
		return "Amount • type value • h/l assets • p recipients • s simulate • S send"
	case enums.FocusTokens:
		return "Assets • j/k choose token • h amount • p recipients • s simulate"
	case enums.FocusContacts:
		return "Recipients • j/k choose • enter select • h amount • l assets"
	default:
		return "hjkl move • p recipients • s simulate • S send"
	}
}
