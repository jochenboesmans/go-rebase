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

func TestRebaseMarket(t *testing.T) {
	Convey("returns expected values for simple market", t, func() {
		mockPairA := m.Pair{
			BaseId:  "1",
			QuoteId: "2",
			ExchangeMarketDataByExchangeId: map[string]m.ExchangeMarketData{
				"ex1": {
					LastPrice:  3,
					CurrentBid: 3,
					CurrentAsk: 3,
					BaseVolume: 1,
				},
			},
		}
		mockPairB := m.Pair{
			BaseId:  "2",
			QuoteId: "3",
			ExchangeMarketDataByExchangeId: map[string]m.ExchangeMarketData{
				"ex1": {
					LastPrice:  2,
					CurrentBid: 2,
					CurrentAsk: 2,
					BaseVolume: 1,
				},
			},
		}
		mockMarket := m.Market{
			PairsById: map[string]m.Pair{
				mockPairA.Id(): mockPairA,
				mockPairB.Id(): mockPairB,
			},
		}

		RebaseMarket("1", 2, &mockMarket)

		// expect pair a's rates not to have changed since it's based in "1" already
		expectedPairA := m.Pair{
			BaseId:  "1",
			QuoteId: "2",
			ExchangeMarketDataByExchangeId: map[string]m.ExchangeMarketData{
				"ex1": {
					LastPrice:  3,
					CurrentBid: 3,
					CurrentAsk: 3,
					BaseVolume: 1,
				},
			},
		}

		// expect pair b's rates to be the product of its own rates and pair a's rates since it's possible to rebase via id "2"
		expectedPairB := m.Pair{
			BaseId:  "2",
			QuoteId: "3",
			ExchangeMarketDataByExchangeId: map[string]m.ExchangeMarketData{
				"ex1": {
					LastPrice:  6,
					CurrentBid: 6,
					CurrentAsk: 6,
					BaseVolume: 3,
				},
			},
		}

		expectedMarket := m.Market{
			PairsById: map[string]m.Pair{
				expectedPairA.Id(): expectedPairA,
				expectedPairB.Id(): expectedPairB,
			},
		}

		So(mockMarket, ShouldResemble, expectedMarket)
	})
}
