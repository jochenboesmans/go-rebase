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
				baseAlreadyIn := false
				quoteAlreadyIn := false
				for _, pairId := range rebaseNeighbors[pairAId].Base {
					if pairId == pairBId {
						baseAlreadyIn = true
					}
				}
				for _, pairId := range rebaseNeighbors[pairBId].Quote {
					if pairId == pairAId {
						quoteAlreadyIn = true
					}
				}
				aBaseUpdated := rebaseNeighbors[pairAId].Base
				if !baseAlreadyIn {
					aBaseUpdated = append(rebaseNeighbors[pairAId].Base, pairBId)
				}
				bQuoteUpdated := rebaseNeighbors[pairBId].Quote
				if !quoteAlreadyIn {
					bQuoteUpdated = append(rebaseNeighbors[pairBId].Quote, pairAId)
				}
				rebaseNeighbors[pairAId] = Neighbors{
					Base:  aBaseUpdated,
					Quote: rebaseNeighbors[pairAId].Quote,
				}
				rebaseNeighbors[pairBId] = Neighbors{
					Base:  rebaseNeighbors[pairBId].Base,
					Quote: bQuoteUpdated,
				}
			}
			if pairB.BaseId == pairA.QuoteId {
				baseAlreadyIn := false
				quoteAlreadyIn := false
				for _, pairId := range rebaseNeighbors[pairBId].Base {
					if pairId == pairAId {
						baseAlreadyIn = true
					}
				}
				for _, pairId := range rebaseNeighbors[pairAId].Quote {
					if pairId == pairBId {
						quoteAlreadyIn = true
					}
				}
				bBaseUpdated := rebaseNeighbors[pairBId].Base
				if !baseAlreadyIn {
					bBaseUpdated = append(rebaseNeighbors[pairBId].Base, pairAId)
				}
				aQuoteUpdated := rebaseNeighbors[pairAId].Quote
				if !quoteAlreadyIn {
					aQuoteUpdated = append(rebaseNeighbors[pairAId].Quote, pairBId)
				}
				rebaseNeighbors[pairBId] = Neighbors{
					Base:  bBaseUpdated,
					Quote: rebaseNeighbors[pairBId].Quote,
				}
				rebaseNeighbors[pairAId] = Neighbors{
					Base:  rebaseNeighbors[pairAId].Base,
					Quote: aQuoteUpdated,
				}
			}
		}
	}

	return rebaseNeighbors
}
