package rebasing

import (
	"fmt"
	"sync"

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

func RebaseMarket(rebaseId string, pathDepth uint8, market *m.Market) {
	var waitGroup sync.WaitGroup
	for pairId := range market.PairsById {
		waitGroup.Add(1)
		go rebasePair(pairId, rebaseId, pathDepth, market, &waitGroup)
	}
	waitGroup.Wait()
}

func rebasePair(pairId string, rebaseId string, pathDepth uint8, market *m.Market, waitGroup *sync.WaitGroup) {
	initialPath := []string{pairId}
	rebasePaths := rebasePathsType{
		Base:  rebasePaths(BASE, initialPath, rebaseId, pathDepth, market),
		Quote: rebasePaths(QUOTE, initialPath, rebaseId, pathDepth, market),
	}
	fmt.Printf("pairId: %s\n", pairId)
	fmt.Printf("rebasePaths: %+v\n", rebasePaths)

	pair := market.PairsById[pairId]

	fmt.Printf("pairBefore: %+v\n", pair)

	for exchangeId, emd := range pair.ExchangeMarketDataByExchangeId {
		if rebasedCurrentAsk, err := deeplyRebaseRate(emd.CurrentAsk, rebaseId, rebasePaths, market); err == nil {
			emd.CurrentAsk = rebasedCurrentAsk
		}
		if rebasedCurrentBid, err := deeplyRebaseRate(emd.CurrentBid, rebaseId, rebasePaths, market); err == nil {
			emd.CurrentBid = rebasedCurrentBid
		}
		if rebasedLastPrice, err := deeplyRebaseRate(emd.LastPrice, rebaseId, rebasePaths, market); err == nil {
			emd.LastPrice = rebasedLastPrice
		}
		if rebasedBaseVolume, err := deeplyRebaseRate(emd.BaseVolume, rebaseId, rebasePaths, market); err == nil {
			emd.BaseVolume = rebasedBaseVolume
		}
		pair.ExchangeMarketDataByExchangeId[exchangeId] = emd
	}

	fmt.Printf("pairAfter: %+v\n", pair)

	market.PairsById[pairId] = pair
	waitGroup.Done()
}

func rebasePaths(direction rebaseDirection, pathAccumulator []string, rebaseId string, pathDepth uint8, market *m.Market) [][]string {
	fmt.Printf("REBASE: direction: %d\n", direction)
	fmt.Printf("rebaseId: %s\n", rebaseId)
	fmt.Printf("pathAccumulator: %+v\n", pathAccumulator)

	if len(pathAccumulator) > int(pathDepth) {
		return [][]string{}
	} else {
		lastPairId := pathAccumulator[0]
		lastBaseId := market.PairsById[lastPairId].BaseId
		fmt.Printf("lastBaseId: %s\n", lastBaseId)
		if lastBaseId == rebaseId {
			return [][]string{pathAccumulator}
		} else {
			return doRebasePaths(direction, pathAccumulator, rebaseId, pathDepth, market)
		}
	}
}

func doRebasePaths(direction rebaseDirection, pathAccumulator []string, rebaseId string, pathDepth uint8, market *m.Market) [][]string {
	var nextNeighborIds []string
	if direction == BASE {
		nextNeighborIds = market.PairsById[pathAccumulator[0]].BaseNeighborIds(market)
	} else if direction == QUOTE {
		nextNeighborIds = market.PairsById[pathAccumulator[0]].QuoteNeighborIds(market)
	}
	fmt.Printf("nextNeighborIds: %+v\n", nextNeighborIds)

	var result [][]string
	for _, nextNeighborId := range nextNeighborIds {
		nextPath := append([]string{nextNeighborId}, pathAccumulator...)
		fmt.Printf("nextPath: %+v\n", nextPath)
		result = append(result, rebasePaths(direction, nextPath, rebaseId, pathDepth, market)...)
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
