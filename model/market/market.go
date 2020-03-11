package market

/**
Collection of pairs for which all market data is based in each pair's base token.
*/
type Market struct {
	PairsById map[string]Pair
}

type Neighbors struct {
	Base  []string
	Quote []string
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
			if pairA.BaseAssetId == pairB.QuoteAssetId {
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
