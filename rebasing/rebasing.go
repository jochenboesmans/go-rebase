package rebasing

import (
	"fmt"
	m "github.com/jochenboesmans/go-rebase/model/market"
)

type rebaseDirection uint8

const (
	BASE = iota + 1
	QUOTE
)

type rebasePathsType struct {
	Base  [][]string
	Quote [][]string
}

func RebaseMarket(rebaseId string, maxPathDepth uint8, market *m.Market) *m.Market {
	rebasedMarket := m.Market{PairsById: map[string]m.Pair{}}
	for pairId := range market.PairsById {
		rebasedMarket.PairsById[pairId] = *rebasePair(pairId, rebaseId, maxPathDepth, market)
	}
	return &rebasedMarket
}

func rebasePair(pairId string, rebaseId string, maxPathDepth uint8, market *m.Market) *m.Pair {
	// determine all paths from the current pair to pairs based in rebaseId
	rebasePaths := rebasePathsType{
		Base:  rebasePaths(BASE, []string{pairId}, rebaseId, maxPathDepth, market),
		Quote: rebasePaths(QUOTE, []string{pairId}, rebaseId, maxPathDepth, market),
	}

	originalPair := market.PairsById[pairId]

	rebasedPair := m.Pair{
		BaseSymbol:                     originalPair.BaseSymbol,
		QuoteSymbol:                    originalPair.QuoteSymbol,
		BaseId:                         originalPair.BaseId,
		QuoteId:                        originalPair.QuoteId,
		ExchangeMarketDataByExchangeId: map[string]m.ExchangeMarketData{},
	}

	// deeply rebase all rates based on the available rebasePaths
	for exchangeId, emd := range originalPair.ExchangeMarketDataByExchangeId {
		rebasedEmd := m.ExchangeMarketData{}
		if rebasedCurrentAsk, err := deeplyRebaseRate(emd.CurrentAsk, rebaseId, rebasePaths, market); err == nil {
			rebasedEmd.CurrentAsk = rebasedCurrentAsk
		}
		if rebasedCurrentBid, err := deeplyRebaseRate(emd.CurrentBid, rebaseId, rebasePaths, market); err == nil {
			rebasedEmd.CurrentBid = rebasedCurrentBid
		}
		if rebasedLastPrice, err := deeplyRebaseRate(emd.LastPrice, rebaseId, rebasePaths, market); err == nil {
			rebasedEmd.LastPrice = rebasedLastPrice
		}
		if rebasedBaseVolume, err := deeplyRebaseRate(emd.BaseVolume, rebaseId, rebasePaths, market); err == nil {
			rebasedEmd.BaseVolume = rebasedBaseVolume
		}
		rebasedPair.ExchangeMarketDataByExchangeId[exchangeId] = rebasedEmd
	}

	return &rebasedPair
}

func rebasePaths(direction rebaseDirection, pathAccumulator []string, rebaseId string, maxPathDepth uint8, market *m.Market) [][]string {
	if len(pathAccumulator) > int(maxPathDepth) {
		return [][]string{}
	} else {
		lastPairId := pathAccumulator[0]
		lastBaseId := market.PairsById[lastPairId].BaseId
		if lastBaseId == rebaseId {
			return [][]string{pathAccumulator}
		} else {
			return doRebasePaths(direction, pathAccumulator, rebaseId, maxPathDepth, market)
		}
	}
}

func doRebasePaths(direction rebaseDirection, pathAccumulator []string, rebaseId string, maxPathDepth uint8, market *m.Market) [][]string {
	var nextNeighborIds []string
	if direction == BASE {
		nextNeighborIds = market.PairsById[pathAccumulator[0]].BaseNeighborIds(market)
	} else if direction == QUOTE {
		nextNeighborIds = market.PairsById[pathAccumulator[0]].QuoteNeighborIds(market)
	}

	var result [][]string
	for _, nextNeighborId := range nextNeighborIds {
		nextPath := append([]string{nextNeighborId}, pathAccumulator...)
		result = append(result, rebasePaths(direction, nextPath, rebaseId, maxPathDepth, market)...)
	}
	return result
}

func deeplyRebaseRate(rate float32, rebaseId string, rebasePaths rebasePathsType, market *m.Market) (float32, error) {
	combinedVolume := float32(0)
	volumeWeightedSum := float32(0)
	for _, baseRebasePath := range rebasePaths.Base {
		rebasedRateAcc := rate
		weightedSumAcc := float32(0)
		for i := len(baseRebasePath) - 2; i >= 0; i-- {
			pair := market.PairsById[baseRebasePath[i]]
			baseId := pair.BaseId
			quoteId := pair.QuoteId
			if rebasedRate, err := shallowlyRebaseRate(rebasedRateAcc, baseId, quoteId, market); err == nil {
				rebasedRateAcc = rebasedRate
			}
			combinedVolume := pair.CombinedBaseVolume()
			if rebasedCombinedVolume, err := shallowlyRebaseRate(combinedVolume, rebaseId, baseId, market); err == nil {
				weightedSumAcc += rebasedCombinedVolume
			}
		}
		weight := weightedSumAcc / float32(len(baseRebasePath))
		combinedVolume += weight
		volumeWeightedSum += weight * rebasedRateAcc
	}

	for _, quoteRebasePath := range rebasePaths.Quote {
		rebasedRateAcc := rate
		weightedSumAcc := float32(0)
		for i := len(quoteRebasePath) - 1; i >= 0; i-- {
			pair := market.PairsById[quoteRebasePath[i]]
			baseId := pair.BaseId
			quoteId := pair.QuoteId
			if i == 0 {
				combinedVolume := pair.CombinedBaseVolume()
				if rebasedCombinedVolume, err := shallowlyRebaseRate(combinedVolume, rebaseId, baseId, market); err == nil {
					weightedSumAcc += rebasedCombinedVolume
				}
			} else if i == len(quoteRebasePath)-1 {
				if rebasedRate, err := shallowlyRebaseRate(rebasedRateAcc, quoteId, baseId, market); err == nil {
					rebasedRateAcc = rebasedRate
				}
			} else {
				pair := market.PairsById[quoteRebasePath[i]]
				baseId := pair.BaseId
				quoteId := pair.QuoteId
				if rebasedRate, err := shallowlyRebaseRate(rebasedRateAcc, quoteId, baseId, market); err == nil {
					rebasedRateAcc = rebasedRate
				}
				combinedVolume := pair.CombinedBaseVolume()
				if rebasedCombinedVolume, err := shallowlyRebaseRate(combinedVolume, rebaseId, baseId, market); err == nil {
					weightedSumAcc += rebasedCombinedVolume
				}
			}
		}
		weight := weightedSumAcc / float32(len(quoteRebasePath))
		combinedVolume += weight
		volumeWeightedSum += weight * rebasedRateAcc
	}

	if combinedVolume == 0 {
		return 0, fmt.Errorf("division by 0")
	} else {
		return volumeWeightedSum / combinedVolume, nil
	}
}

func shallowlyRebaseRate(rate float32, rebaseId string, baseId string, market *m.Market) (float32, error) {
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
