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
	Base [][]string
	Quote [][]string
}

func RebaseMarket(rebaseId string, pathDepth uint8, market *m.Market) {
	for pairId := range market.PairsById {
		go rebasePair(pairId, rebaseId, pathDepth, market)
	}
}

func rebasePair(pairId string, rebaseId string, pathDepth uint8, market *m.Market) {
	initialPath := []string{pairId}
	rebasePaths := rebasePathsType{
		Base:  rebasePaths(BASE, initialPath, rebaseId, pathDepth, market),
		Quote: rebasePaths(QUOTE, initialPath, rebaseId, pathDepth, market),
	}

	pair := market.PairsById[pairId]

	for exchangeId, emd := range pair.ExchangeMarketDataByExchangeId {
		emd.CurrentAsk = deeplyRebaseRate(emd.CurrentAsk, rebasePaths)
		emd.CurrentBid = deeplyRebaseRate(emd.CurrentBid, rebasePaths)
		emd.LastPrice = deeplyRebaseRate(emd.LastPrice, rebasePaths)
		emd.BaseVolume = deeplyRebaseRate(emd.BaseVolume, rebasePaths)
		pair.ExchangeMarketDataByExchangeId[exchangeId] = emd
	}

	market.PairsById[pairId] = pair
}

func rebasePaths(direction rebaseDirection, pathAccumulator []string, rebaseId string, pathDepth uint8, market *m.Market) [][]string {
	if len(pathAccumulator) >= int(pathDepth) {
		return [][]string{}
	} else {
		lastPairId := pathAccumulator[0]
		lastBaseId := market.PairsById[lastPairId].BaseId
		if lastBaseId == rebaseId {
			return [][]string{pathAccumulator}
		} else {
			return doRebasePaths(direction, pathAccumulator, rebaseId, pathDepth, market)
		}
	}
}

func doRebasePaths(direction rebaseDirection, pathAccumulator []string, rebaseId string, pathDepth uint8, market *m.Market) [][]string {
	nextNeighbors := n.Neighbors(direction, pathAccumulator[0])

	var result [][]string
	for _, nextNeighbor := range nextNeighbors {
		nextPath := append([]string{nextNeighbor})
		result = append(result, rebasePaths(direction, nextPath, rebaseId, pathDepth, market)...)
	}
	return result
}

func deeplyRebaseRate(rate float32, rebaseId string, rebasePaths rebasePathsType, market *m.Market) (float32, error) {
	baseCombinedVolume := 0
	baseVolumeWeightedSum := 0
	for _, baseRebasePath := range rebasePaths.Base {

	}

	quoteCombinedVolume := 0
	quoteVolumeWeightedSum := 0
	for _, quoteRebasePath := range rebasePaths.Quote {

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
