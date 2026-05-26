package header

import "bos/enums"

func activeHelp(focus enums.FocusArea) string {
	switch focus {
	case enums.FocusAmount:
		return "Amount • type value • h/l assets • p recipients • a add contact • i import token • s simulate • c contract • S send"
	case enums.FocusTokens:
		return "Assets • j/k choose token • enter select • l transactions"
	case enums.FocusTransactions:
		return "Transactions • j/k choose • enter/space select • h assets"
	case enums.FocusContacts:
		return "Recipients • j/k choose • a add contact • enter select • h amount • l assets"
	default:
		return "hjkl move • t assets • x transactions • p recipients • a add contact • i import token • s simulate • c contract • S send"
	}
}
