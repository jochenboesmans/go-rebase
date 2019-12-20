package market

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

/**
Market pair containing data from one or many exchanges.
*/
type Pair struct {
	BaseSymbol string
	BaseId string
	QuoteSymbol string
	QuoteId string
	ExchangeMarketDataByExchangeId map[ExchangeId]ExchangeMarketData
}

func (p *Pair) Id() string {
	idString := fmt.Sprintf("%s/%s", p.BaseId, p.QuoteId)

	hash := sha1.New()
	hash.Write([]byte(idString))
	result := hash.Sum(nil)

	return hex.EncodeToString(result)
}

func (p *Pair) combinedBaseVolume() float32 {
	var sum float32 = 0
	for _, emd := range p.ExchangeMarketDataByExchangeId {
		sum += emd.BaseVolume
	}
	return sum
}

func (p *Pair) baseVolumeWeightedCurrentBidSum() float32 {
	var sum float32 = 0
	for _, emd := range p.ExchangeMarketDataByExchangeId {
		sum += emd.BaseVolume * emd.CurrentBid
	}
	return sum
}

func (p *Pair) baseVolumeWeightedCurrentAskSum() float32 {
	var sum float32 = 0
	for _, emd := range p.ExchangeMarketDataByExchangeId {
		sum += emd.BaseVolume * emd.CurrentAsk
	}
	return sum
}

func (p *Pair) baseVolumeWeightedSpreadAverage() float32 {
	spreadAverage := (p.baseVolumeWeightedCurrentBidSum() + p.baseVolumeWeightedCurrentAskSum()) / 2
	weightedAverage := spreadAverage / p.combinedBaseVolume()
	return weightedAverage
}






