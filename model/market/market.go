package market

/**
Collection of pairs for which all market data is based in each pair's base token.
*/
type Market struct {
	PairsById map[string]Pair `json:"pairsById"`
}

/**
Returns a collection of exchanges currently included in the market.
*/
func (m *Market) ExchangeIds() []string {
	var result []string
	for _, p := range m.PairsById {
		for eId := range p.ExchangeMarketDataByExchangeId {
			result = append(result, eId)
		}
	}
	return result
}
