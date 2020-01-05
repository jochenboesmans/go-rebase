package rebasing

import (
	m "github.com/jochenboesmans/go-rebase/model/market"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestShallowlyRebaseRate(t *testing.T) {
	Convey("rebase ids have matching pair in market", t, func() {
		rate := float32(1.1)
		rebaseId := "0xfoo"
		baseId := "0xbar"
		quoteId := "0xheh"
		originalPair := m.Pair{
			BaseId:  baseId,
			QuoteId: quoteId,
			ExchangeMarketDataByExchangeId: map[string]m.ExchangeMarketData{
				"KYBER": {
					BaseVolume: 1.5,
					CurrentBid: 100.7,
					CurrentAsk: 103.5,
				},
				"UNISWAP": {
					BaseVolume: 3,
					CurrentBid: 150.1,
					CurrentAsk: 155.2,
				},
			},
		}
		rebasePair := m.Pair{
			BaseId:  rebaseId,
			QuoteId: baseId,
			ExchangeMarketDataByExchangeId: map[string]m.ExchangeMarketData{
				"KYBER": {
					BaseVolume: 1.5,
					CurrentBid: 100.7,
					CurrentAsk: 103.5,
				},
				"UNISWAP": {
					BaseVolume: 3,
					CurrentBid: 150.1,
					CurrentAsk: 155.2,
				},
			},
		}
		market := m.Market{
			PairsById: map[string]m.Pair{
				rebasePair.Id():   rebasePair,
				originalPair.Id(): originalPair,
			},
		}
		rebasePairBaseVolumeWeightedSpreadAverage, err1 := rebasePair.BaseVolumeWeightedSpreadAverage()
		// only test in case BaseVolumeWeightedSpreadAverage returns valid response
		if err1 == nil {
			expected := rate * rebasePairBaseVolumeWeightedSpreadAverage
			actual, err := shallowlyRebaseRate(rate, rebaseId, baseId, &market)

			So(err, ShouldBeNil)
			So(actual, ShouldResemble, expected)
		} else {
			expected := 0
			actual, err := shallowlyRebaseRate(rate, rebaseId, baseId, &market)

			So(err, ShouldNotBeNil)
			So(actual, ShouldResemble, expected)
		}
	})
	Convey("rebase id is base id", t, func() {
		rate := float32(1.1)
		rebaseId := "0xfoo"
		baseId := "0xfoo"
		quoteId := "0xheh"
		originalPair := m.Pair{
			BaseId:  baseId,
			QuoteId: quoteId,
		}
		rebasePair := m.Pair{
			BaseId:  rebaseId,
			QuoteId: baseId,
		}
		market := m.Market{
			PairsById: map[string]m.Pair{
				rebasePair.Id():   rebasePair,
				originalPair.Id(): originalPair,
			},
		}
		expected := rate
		actual, err := shallowlyRebaseRate(rate, rebaseId, baseId, &market)

		So(err, ShouldBeNil)
		So(actual, ShouldResemble, expected)
	})
}
