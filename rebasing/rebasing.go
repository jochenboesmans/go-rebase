package rebasing

import (
	"fmt"
	m "github.com/jochenboesmans/go-rebase/model/market"
	"sync"
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
	// shared data structures
	rebasedMarket := m.Market{PairsById: map[string]m.Pair{}}
	rebaseNeighbors := market.RebaseNeighbors()

	var waitGroup sync.WaitGroup
	for pairId := range market.PairsById {
		waitGroup.Add(1)
		rebasePair(pairId, rebaseId, maxPathDepth, market, &rebasedMarket, rebaseNeighbors, &waitGroup)
	}
	waitGroup.Wait()
	return &rebasedMarket
}

func rebasePair(pairId string, rebaseId string, maxPathDepth uint8, market *m.Market, rebasedMarket *m.Market, rebaseNeighbors map[string]m.Neighbors, waitGroup *sync.WaitGroup) {
	// determine all paths from the current pair to pairs based in rebaseId
	rebasePaths := rebasePathsType{
		Base:  rebasePaths(BASE, []string{pairId}, rebaseId, maxPathDepth, market, rebaseNeighbors),
		Quote: rebasePaths(QUOTE, []string{pairId}, rebaseId, maxPathDepth, market, rebaseNeighbors),
	}

	originalMarketPair := market.PairsById[pairId]

	// copy pair data to rebased pair

	// deeply rebase all rates based on the available rebasePaths
	var newExchangeMarkets []m.ExchangeMarket
	for _, emd := range originalMarketPair.ExchangeMarkets {
		newExchangeMarket := m.ExchangeMarket{
			CurrentBid: deeplyRebaseRate(emd.CurrentBid, rebaseId, rebasePaths, market),
			CurrentAsk: deeplyRebaseRate(emd.CurrentAsk, rebaseId, rebasePaths, market),
			BaseVolume: deeplyRebaseRate(emd.BaseVolume, rebaseId, rebasePaths, market),
		}
		newExchangeMarkets = append(newExchangeMarkets, newExchangeMarket)
	}

	rebasedMarket.PairsById[pairId] = m.Pair{
		BaseId:          originalMarketPair.BaseId,
		QuoteId:         originalMarketPair.QuoteId,
		ExchangeMarkets: newExchangeMarkets,
	}

	waitGroup.Done()
}

func rebasePaths(direction rebaseDirection, pathAccumulator []string, rebaseId string, maxPathDepth uint8, market *m.Market, rebaseNeighbors map[string]m.Neighbors) [][]string {
	if len(pathAccumulator) > int(maxPathDepth) {
		return [][]string{}
	} else {
		lastPairId := pathAccumulator[0]
		lastBaseId := market.PairsById[lastPairId].BaseId
		if lastBaseId == rebaseId {
			return [][]string{pathAccumulator}
		} else {
			return doRebasePaths(direction, pathAccumulator, rebaseId, maxPathDepth, market, rebaseNeighbors)
		}
	}
}

func doRebasePaths(direction rebaseDirection, pathAccumulator []string, rebaseId string, maxPathDepth uint8, market *m.Market, rebaseNeighbors map[string]m.Neighbors) [][]string {
	var nextNeighborIds []string
	if direction == BASE {
		nextNeighborIds = rebaseNeighbors[pathAccumulator[0]].Base
	} else if direction == QUOTE {
		nextNeighborIds = rebaseNeighbors[pathAccumulator[0]].Quote
	}

	var result [][]string
	for _, nextNeighborId := range nextNeighborIds {
		nextPath := append([]string{nextNeighborId}, pathAccumulator...)
		result = append(result, rebasePaths(direction, nextPath, rebaseId, maxPathDepth, market, rebaseNeighbors)...)
	}
	return result
}

func deeplyRebaseRate(rate float32, rebaseId string, rebasePaths rebasePathsType, market *m.Market) float32 {
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
		return 0.0
	} else {
		return volumeWeightedSum / combinedVolume
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
			matchingMarketPairBaseVolumeWeightedSpreadAverage := matchingMarketPair.BaseVolumeWeightedSpreadAverage()
			return matchingMarketPairBaseVolumeWeightedSpreadAverage * rate, nil
		}
	}
}
