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
	BaseSymbol                     string                        `json:"baseSymbol"`
	BaseId                         string                        `json:"baseId"`
	QuoteSymbol                    string                        `json:"quoteSymbol"`
	QuoteId                        string                        `json:"quoteId"`
	ExchangeMarketDataByExchangeId map[string]ExchangeMarketData `json:"exchangeMarketDataByExchangeId"`
}

/**
Exchange-specific market data.
*/
type ExchangeMarketData struct {
	LastPrice  float32 `json:"lastPrice"`
	CurrentBid float32 `json:"currentBid"`
	CurrentAsk float32 `json:"currentAsk"`
	BaseVolume float32 `json:"baseVolume"`
}

func (p *Pair) Id() string {
	idString := fmt.Sprintf("%s/%s", p.BaseId, p.QuoteId)

	hash := sha1.New()
	hash.Write([]byte(idString))
	result := hash.Sum(nil)

	return hex.EncodeToString(result)
}

func (p *Pair) CombinedBaseVolume() float32 {
	var sum float32 = 0
	for _, emd := range p.ExchangeMarketDataByExchangeId {
		sum += emd.BaseVolume
	}
	return sum
}

func (p *Pair) BaseVolumeWeightedCurrentBidSum() float32 {
	var sum float32 = 0
	for _, emd := range p.ExchangeMarketDataByExchangeId {
		sum += emd.BaseVolume * emd.CurrentBid
	}
	return sum
}

func (p *Pair) BaseVolumeWeightedCurrentAskSum() float32 {
	var sum float32 = 0
	for _, emd := range p.ExchangeMarketDataByExchangeId {
		sum += emd.BaseVolume * emd.CurrentAsk
	}
	return sum
}

func (p *Pair) BaseVolumeWeightedSpreadAverage() (float32, error) {
	spreadAverage := (p.BaseVolumeWeightedCurrentBidSum() + p.BaseVolumeWeightedCurrentAskSum()) / 2
	if p.CombinedBaseVolume() == float32(0) {
		return 0, fmt.Errorf(`combined base volume is 0 for pair: %+v\n`, p)
	} else {
		weightedAverage := spreadAverage / p.CombinedBaseVolume()
		return weightedAverage, nil
	}
}
