package rebasing

import m "github.com/jochenboesmans/go-rebase/model/market"

type neighborCacheType struct {
	Market      *m.Market
	NeighborIds map[string]map[rebaseDirection][]string
}

var neighborCache neighborCacheType

func NeighborIds(direction rebaseDirection, originalPairId string, market *m.Market) []string {
	if market == neighborCache.Market {
		return neighborCache.NeighborIds[originalPairId][direction]
	} else {
		originalPair := market.PairsById[originalPairId]
		var neighborIdsAcc []string
		if direction == BASE {
			for pairId, pair := range market.PairsById {
				if originalPair.BaseId == pair.QuoteId {
					neighborIdsAcc = append(neighborIdsAcc, pairId)
				}
			}
		} else if direction == QUOTE {
			for pairId, pair := range market.PairsById {
				if originalPair.QuoteId == pair.BaseId {
					neighborIdsAcc = append(neighborIdsAcc, pairId)
				}
			}
		}
		neighborCache.Market = market
		neighborCache.NeighborIds[originalPairId][direction] = neighborIdsAcc
		return neighborIdsAcc
	}
}
