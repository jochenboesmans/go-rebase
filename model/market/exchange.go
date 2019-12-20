package market

/**
Exchange-specific market data.
*/
type ExchangeMarketData struct {
	LastPrice float32
	CurrentBid float32
	CurrentAsk float32
	BaseVolume float32
}

type ExchangeId uint8
const (
	KYBER ExchangeId = iota + 1
	UNISWAP
	IDEX
	OASIS
	RADAR
)
