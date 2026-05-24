package utils_test

import (
	"errors"
	"math/big"
	"testing"

	"bos/utils"
)

func TestParseEtherToWeiValidAmounts(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "whole ether", in: "1", want: "1000000000000000000"},
		{name: "fractional ether", in: "0.5", want: "500000000000000000"},
		{name: "trims whitespace", in: "  2.25  ", want: "2250000000000000000"},
		{name: "leading decimal point", in: ".1", want: "100000000000000000"},
		{name: "eighteen decimals", in: "0.000000000000000001", want: "1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.ParseEtherToWei(tt.in)
			if err != nil {
				t.Fatalf("ParseEtherToWei(%q) returned error: %v", tt.in, err)
			}
			if got.String() != tt.want {
				t.Fatalf("ParseEtherToWei(%q) = %s, want %s", tt.in, got.String(), tt.want)
			}
		})
	}
}

func TestParseEtherToWeiInvalidAmounts(t *testing.T) {
	for _, input := range []string{"", "   ", "0", "-1", "1.2.3", "abc", "1.0000000000000000001"} {
		t.Run(input, func(t *testing.T) {
			if got, err := utils.ParseEtherToWei(input); err == nil {
				t.Fatalf("ParseEtherToWei(%q) = %s, want error", input, got.String())
			}
		})
	}
}

func TestWeiToEther(t *testing.T) {
	tests := []struct {
		name string
		wei  string
		want string
	}{
		{name: "zero", wei: "0", want: "0.00000000"},
		{name: "one ether", wei: "1000000000000000000", want: "1.00000000"},
		{name: "fraction", wei: "1234567890000000000", want: "1.23456789"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wei, ok := new(big.Int).SetString(tt.wei, 10)
			if !ok {
				t.Fatalf("invalid test wei value %q", tt.wei)
			}
			if got := utils.WeiToEther(wei); got != tt.want {
				t.Fatalf("WeiToEther(%s) = %q, want %q", tt.wei, got, tt.want)
			}
		})
	}
}

func TestSafeWidth(t *testing.T) {
	if got := utils.SafeWidth(80); got != 100 {
		t.Fatalf("SafeWidth(80) = %d, want 100", got)
	}
	if got := utils.SafeWidth(120); got != 120 {
		t.Fatalf("SafeWidth(120) = %d, want 120", got)
	}
}

func TestErrorMessage(t *testing.T) {
	if got := utils.ErrorMessage(nil); got != "unknown error" {
		t.Fatalf("ErrorMessage(nil) = %q, want %q", got, "unknown error")
	}
	if got := utils.ErrorMessage(errors.New("boom")); got != "boom" {
		t.Fatalf("ErrorMessage(error) = %q, want %q", got, "boom")
	}
}
