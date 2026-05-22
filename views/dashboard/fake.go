package dashboard

import "bos/views"

func fakeTokens() []views.Token {
	return []views.Token{
		{Symbol: "ETH", Name: "Native Network Token", Balance: "0", Address: "native", Native: true},
		{Symbol: "USDT", Name: "Tether USD", Balance: "1240.22", Address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"},
		{Symbol: "LINK", Name: "Chainlink", Balance: "4.12", Address: "0x514910771af9ca656af840dff83e8264ecf986ca"},
		{Symbol: "UNI", Name: "Uniswap", Balance: "42.11", Address: "0x1f9840a85d5af5bf1d1762f925bdaddc4201f984"},
	}
}

func fakeContacts() []views.Contact {
	return []views.Contact{
		{Name: "Treasury Wallet", Address: "0x1111111111111111111111111111111111111111"},
		{Name: "Personal Wallet", Address: "0x2222222222222222222222222222222222222222"},
		{Name: "Binance Deposit", Address: "0x3333333333333333333333333333333333333333"},
	}
}
