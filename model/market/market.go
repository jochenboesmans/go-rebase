package market

/**
Collection of pairs for which all market data is based in each pair's base token.
*/
type Market struct {
	PairsById map[string]Pair `json:"pairsById"`
}

type Neighbors struct {
	Base  []string
	Quote []string
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

func (m *Market) RebaseNeighbors() map[string]Neighbors {
	rebaseNeighbors := map[string]Neighbors{}

	for pairId := range m.PairsById {
		rebaseNeighbors[pairId] = Neighbors{
			Base:  []string{},
			Quote: []string{},
		}
	}
	for pairAId, pairA := range m.PairsById {
		for pairBId, pairB := range m.PairsById {
			if pairA.BaseId == pairB.QuoteId {
				rebaseNeighbors[pairAId] = Neighbors{
					Base:  append(rebaseNeighbors[pairAId].Base, pairBId),
					Quote: rebaseNeighbors[pairAId].Quote,
				}
				rebaseNeighbors[pairBId] = Neighbors{
					Base:  rebaseNeighbors[pairBId].Base,
					Quote: append(rebaseNeighbors[pairBId].Quote, pairAId),
				}
			}
		}
	}

	return rebaseNeighbors
}
