package utils

import (
	"errors"
	"math/big"
	"strings"
)

func ParseEtherToWei(input string) (*big.Int, error) {
	input = strings.TrimSpace(input)

	if input == "" {
		return nil, errors.New("amount is empty")
	}

	if strings.HasPrefix(input, "-") {
		return nil, errors.New("amount cannot be negative")
	}

	parts := strings.Split(input, ".")
	if len(parts) > 2 {
		return nil, errors.New("invalid amount")
	}

	whole := parts[0]
	fraction := ""

	if len(parts) == 2 {
		fraction = parts[1]
	}

	if whole == "" {
		whole = "0"
	}

	if len(fraction) > 18 {
		return nil, errors.New("amount has more than 18 decimal places")
	}

	for _, ch := range whole + fraction {
		if ch < '0' || ch > '9' {
			return nil, errors.New("amount must contain only digits and optional decimal point")
		}
	}

	fraction = fraction + strings.Repeat("0", 18-len(fraction))

	weiString := whole + fraction
	weiString = strings.TrimLeft(weiString, "0")

	if weiString == "" {
		weiString = "0"
	}

	wei, ok := new(big.Int).SetString(weiString, 10)
	if !ok {
		return nil, errors.New("failed to parse amount")
	}

	if wei.Sign() == 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	return wei, nil
}
