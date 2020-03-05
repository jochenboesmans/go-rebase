package rebasing

import (
	m "github.com/jochenboesmans/go-rebase/model/market"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRebasePaths(t *testing.T) {
	Convey("rebase path in base direction but not in quote direction", t, func() {
		mockPairA := m.Pair{
			BaseId:  "1",
			QuoteId: "2",
		}
		mockPairB := m.Pair{
			BaseId:  "2",
			QuoteId: "3",
		}
		mockMarket := m.Market{
			PairsById: map[string]m.Pair{
				mockPairA.Id(): mockPairA,
				mockPairB.Id(): mockPairB,
			},
		}

		rebasePaths := rebasePathsType{
			Base:  rebasePaths(BASE, []string{mockPairB.Id()}, "1", 2, &mockMarket),
			Quote: rebasePaths(QUOTE, []string{mockPairB.Id()}, "1", 2, &mockMarket),
		}

		expectedBase := [][]string{{mockPairA.Id(), mockPairB.Id()}}
		var expectedQuote [][]string

		So(rebasePaths.Base, ShouldResemble, expectedBase)
		So(rebasePaths.Quote, ShouldResemble, expectedQuote)
	})
	Convey("rebase path in both the quote and the base direction", t, func() {
		mockPairA := m.Pair{
			BaseId:  "1",
			QuoteId: "2",
		}
		mockPairB := m.Pair{
			BaseId:  "2",
			QuoteId: "3",
		}
		mockPairC := m.Pair{
			BaseId:  "3",
			QuoteId: "1",
		}

		mockMarket := m.Market{
			PairsById: map[string]m.Pair{
				mockPairA.Id(): mockPairA,
				mockPairB.Id(): mockPairB,
				mockPairC.Id(): mockPairC,
			},
		}

		rebasePaths := rebasePathsType{
			Base:  rebasePaths(BASE, []string{mockPairC.Id()}, "1", 3, &mockMarket),
			Quote: rebasePaths(QUOTE, []string{mockPairC.Id()}, "1", 3, &mockMarket),
		}

		expectedBase := [][]string{{mockPairA.Id(), mockPairB.Id(), mockPairC.Id()}}
		expectedQuote := [][]string{{mockPairA.Id(), mockPairC.Id()}}

		So(rebasePaths.Base, ShouldResemble, expectedBase)
		So(rebasePaths.Quote, ShouldResemble, expectedQuote)
	})
}

func TestShallowlyRebaseRate(t *testing.T) {
	Convey("rebase ids have matching pair in mockMarket", t, func() {
		mockPairA := m.Pair{
			BaseId:  "1",
			QuoteId: "2",
			ExchangeMarketDataByExchangeId: map[string]m.ExchangeMarketData{
				"ex1": {
					LastPrice:  2,
					CurrentBid: 2,
					CurrentAsk: 2,
					BaseVolume: 1,
				},
			},
		}
		mockPairB := m.Pair{
			BaseId:  "2",
			QuoteId: "3",
			ExchangeMarketDataByExchangeId: map[string]m.ExchangeMarketData{
				"ex1": {
					LastPrice:  3,
					CurrentBid: 3,
					CurrentAsk: 3,
					BaseVolume: 1,
				},
			},
		}

		rate := float32(1.1)

		mockMarket := m.Market{
			PairsById: map[string]m.Pair{
				mockPairA.Id(): mockPairA,
				mockPairB.Id(): mockPairB,
			},
		}
		rebasePairBaseVolumeWeightedSpreadAverage, err1 := mockPairA.BaseVolumeWeightedSpreadAverage()
		// only test in case BaseVolumeWeightedSpreadAverage returns valid response
		if err1 == nil {
			expected := rate * rebasePairBaseVolumeWeightedSpreadAverage
			actual, err := shallowlyRebaseRate(rate, "1", mockPairB.BaseId, &mockMarket)

			So(err, ShouldBeNil)
			So(actual, ShouldResemble, expected)
		} else {
			expected := 0
			actual, err := shallowlyRebaseRate(rate, "1", mockPairB.BaseId, &mockMarket)

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
		mockMarket := m.Market{
			PairsById: map[string]m.Pair{
				rebasePair.Id():   rebasePair,
				originalPair.Id(): originalPair,
			},
		}
		expected := rate
		actual, err := shallowlyRebaseRate(rate, rebaseId, baseId, &mockMarket)

		So(err, ShouldBeNil)
		So(actual, ShouldResemble, expected)
	})
}

func TestRebaseMarket(t *testing.T) {
	Convey("returns expected values for simple mockMarket", t, func() {
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
	Convey("more complex mockMarket with longer path to rebase pair", t, func() {
		mockPairA := m.Pair{
			BaseId:  "1",
			QuoteId: "2",
			ExchangeMarketDataByExchangeId: map[string]m.ExchangeMarketData{
				"ex1": {
					LastPrice:  2,
					CurrentBid: 2,
					CurrentAsk: 2,
					BaseVolume: 1,
				},
			},
		}
		mockPairB := m.Pair{
			BaseId:  "2",
			QuoteId: "3",
			ExchangeMarketDataByExchangeId: map[string]m.ExchangeMarketData{
				"ex1": {
					LastPrice:  3,
					CurrentBid: 3,
					CurrentAsk: 3,
					BaseVolume: 1,
				},
			},
		}
		mockPairC := m.Pair{
			BaseId:  "3",
			QuoteId: "4",
			ExchangeMarketDataByExchangeId: map[string]m.ExchangeMarketData{
				"ex1": {
					LastPrice:  4,
					CurrentBid: 4,
					CurrentAsk: 4,
					BaseVolume: 1,
				},
			},
		}
		mockMarket := m.Market{
			PairsById: map[string]m.Pair{
				mockPairA.Id(): mockPairA,
				mockPairB.Id(): mockPairB,
				mockPairC.Id(): mockPairC,
			},
		}

		RebaseMarket("1", 3, &mockMarket)

		// expect pair a's rates not to have changed since it's based in "1" already
		expectedPairA := m.Pair{
			BaseId:  "1",
			QuoteId: "2",
			ExchangeMarketDataByExchangeId: map[string]m.ExchangeMarketData{
				"ex1": {
					LastPrice:  2,
					CurrentBid: 2,
					CurrentAsk: 2,
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
					BaseVolume: 2,
				},
			},
		}
		// expect pair c's rates to be the product of its own rates and pair b's rates since it's possible to rebase via id "3"
		expectedPairC := m.Pair{
			BaseId:  "3",
			QuoteId: "4",
			ExchangeMarketDataByExchangeId: map[string]m.ExchangeMarketData{
				"ex1": {
					LastPrice:  24,
					CurrentBid: 24,
					CurrentAsk: 24,
					BaseVolume: 6,
				},
			},
		}

		expectedMarket := m.Market{
			PairsById: map[string]m.Pair{
				expectedPairA.Id(): expectedPairA,
				expectedPairB.Id(): expectedPairB,
				expectedPairC.Id(): expectedPairC,
			},
		}

		So(mockMarket, ShouldResemble, expectedMarket)
	})
	Convey("doesn't change rates when there is no path to rebase id", t, func() {
		mockPairA := m.Pair{
			BaseId:  "1",
			QuoteId: "2",
			ExchangeMarketDataByExchangeId: map[string]m.ExchangeMarketData{
				"ex1": {
					LastPrice:  2,
					CurrentBid: 2,
					CurrentAsk: 2,
					BaseVolume: 1,
				},
			},
		}
		mockPairB := m.Pair{
			BaseId:  "3",
			QuoteId: "4",
			ExchangeMarketDataByExchangeId: map[string]m.ExchangeMarketData{
				"ex1": {
					LastPrice:  3,
					CurrentBid: 3,
					CurrentAsk: 3,
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
					LastPrice:  2,
					CurrentBid: 2,
					CurrentAsk: 2,
					BaseVolume: 1,
				},
			},
		}

		// expect pair b's rates not to have changed since there's not path to the rebase id
		expectedPairB := m.Pair{
			BaseId:  "3",
			QuoteId: "4",
			ExchangeMarketDataByExchangeId: map[string]m.ExchangeMarketData{
				"ex1": {
					LastPrice:  3,
					CurrentBid: 3,
					CurrentAsk: 3,
					BaseVolume: 1,
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
