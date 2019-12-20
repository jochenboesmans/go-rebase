package rebasing

import (
	"fmt"
	m "github.com/jochenboesmans/go-rebase/model/market"
)

func RebasedRate(rate float32, rebaseId string, baseId string, market m.Market) (float32, error) {
	if rebaseId == baseId {
		return rate, nil
	} else {
		rebasePair := m.Pair{
			BaseId:  rebaseId,
			QuoteId: baseId,
		}
		if matchingMarketPair, ok := market.PairsById[rebasePair.Id()]; !ok {
			return 0, fmt.Errorf(`no pair in market to rebase baseId "%s" to rebaseId "%s"`, baseId, rebaseId)
		} else {
			if matchingMarketPairBaseVolumeWeightedSpreadAverage, err := matchingMarketPair.BaseVolumeWeightedSpreadAverage(); err != nil {
				return 0, err
			} else {
				return rate * matchingMarketPairBaseVolumeWeightedSpreadAverage, nil
			}
		}
	}
}
