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
		mockRebaseNeighbors := mockMarket.RebaseNeighbors()

		rebasePaths := rebasePathsType{
			Base:  rebasePaths(BASE, []string{mockPairB.Id()}, "1", 2, &mockMarket, mockRebaseNeighbors),
			Quote: rebasePaths(QUOTE, []string{mockPairB.Id()}, "1", 2, &mockMarket, mockRebaseNeighbors),
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
		mockRebaseNeighbors := mockMarket.RebaseNeighbors()

		rebasePaths := rebasePathsType{
			Base:  rebasePaths(BASE, []string{mockPairC.Id()}, "1", 3, &mockMarket, mockRebaseNeighbors),
			Quote: rebasePaths(QUOTE, []string{mockPairC.Id()}, "1", 3, &mockMarket, mockRebaseNeighbors),
		}

		expectedBase := [][]string{{mockPairA.Id(), mockPairB.Id(), mockPairC.Id()}}
		expectedQuote := [][]string{{mockPairA.Id(), mockPairC.Id()}}

		So(rebasePaths.Base, ShouldResemble, expectedBase)
		So(rebasePaths.Quote, ShouldResemble, expectedQuote)
	})
	Convey("multiple rebase paths in base direction", t, func() {
		mockPairA := m.Pair{
			BaseId:  "1",
			QuoteId: "2",
		}
		mockPairB := m.Pair{
			BaseId:  "1",
			QuoteId: "3",
		}
		mockPairC := m.Pair{
			BaseId:  "2",
			QuoteId: "4",
		}
		mockPairD := m.Pair{
			BaseId:  "3",
			QuoteId: "4",
		}
		mockPairE := m.Pair{
			BaseId:  "4",
			QuoteId: "6",
		}
		mockMarket := m.Market{
			PairsById: map[string]m.Pair{
				mockPairA.Id(): mockPairA,
				mockPairB.Id(): mockPairB,
				mockPairC.Id(): mockPairC,
				mockPairD.Id(): mockPairD,
				mockPairE.Id(): mockPairE,
			},
		}
		mockRebaseNeighbors := mockMarket.RebaseNeighbors()

		rebasePaths := rebasePathsType{
			Base:  rebasePaths(BASE, []string{mockPairE.Id()}, "1", 5, &mockMarket, mockRebaseNeighbors),
			Quote: rebasePaths(QUOTE, []string{mockPairE.Id()}, "1", 5, &mockMarket, mockRebaseNeighbors),
		}

		expectedBasePath1 := []string{mockPairA.Id(), mockPairC.Id(), mockPairE.Id()}
		expectedBasePath2 := []string{mockPairB.Id(), mockPairD.Id(), mockPairE.Id()}
		var expectedQuote [][]string

		So(rebasePaths.Base, ShouldContain, expectedBasePath1)
		So(rebasePaths.Base, ShouldContain, expectedBasePath2)
		So(rebasePaths.Quote, ShouldResemble, expectedQuote)
	})
}

func TestShallowlyRebaseRate(t *testing.T) {
	Convey("rebase pair not in market", t, func() {
		mockRebaseId := "1"
		mockBaseId := "2"
		// pair not in market
		mockMarket := m.Market{
			PairsById: map[string]m.Pair{},
		}
		mockRate := float32(1.0)

		actualRebaseRate, actualError := shallowlyRebaseRate(mockRate, mockRebaseId, mockBaseId, &mockMarket)

		So(actualRebaseRate, ShouldEqual, float32(0.0))
		So(actualError, ShouldNotBeNil)
	})
	Convey("rebase ids have matching pair in mockMarket", t, func() {
		mockPairA := m.Pair{
			BaseId:  "1",
			QuoteId: "2",
			ExchangeMarkets: []m.ExchangeMarket{
				{
					CurrentBid: 2,
					CurrentAsk: 2,
					BaseVolume: 1,
				},
			},
		}
		mockPairB := m.Pair{
			BaseId:  "2",
			QuoteId: "3",
			ExchangeMarkets: []m.ExchangeMarket{
				{
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
		rebasePairBaseVolumeWeightedSpreadAverage := mockPairA.BaseVolumeWeightedSpreadAverage()
		expected := rate * rebasePairBaseVolumeWeightedSpreadAverage
		actual, err := shallowlyRebaseRate(rate, "1", mockPairB.BaseId, &mockMarket)

		So(err, ShouldBeNil)
		So(actual, ShouldResemble, expected)
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
			ExchangeMarkets: []m.ExchangeMarket{
				{
					CurrentBid: 3,
					CurrentAsk: 3,
					BaseVolume: 1,
				},
			},
		}
		mockPairB := m.Pair{
			BaseId:  "2",
			QuoteId: "3",
			ExchangeMarkets: []m.ExchangeMarket{
				{
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

		actualMarket := RebaseMarket("1", 2, &mockMarket)

		// expect pair a's rates not to have changed since it's based in "1" already
		expectedPairA := m.Pair{
			BaseId:  "1",
			QuoteId: "2",
			ExchangeMarkets: []m.ExchangeMarket{
				{
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
			ExchangeMarkets: []m.ExchangeMarket{
				{
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

		So(actualMarket, ShouldResemble, &expectedMarket)
	})
	Convey("more complex mockMarket with longer path to rebase pair", t, func() {
		mockPairA := m.Pair{
			BaseId:  "1",
			QuoteId: "2",
			ExchangeMarkets: []m.ExchangeMarket{
				{
					CurrentBid: 2,
					CurrentAsk: 2,
					BaseVolume: 1,
				},
			},
		}
		mockPairB := m.Pair{
			BaseId:  "2",
			QuoteId: "3",
			ExchangeMarkets: []m.ExchangeMarket{
				{
					CurrentBid: 3,
					CurrentAsk: 3,
					BaseVolume: 1,
				},
			},
		}
		mockPairC := m.Pair{
			BaseId:  "3",
			QuoteId: "4",
			ExchangeMarkets: []m.ExchangeMarket{
				{
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

		actualMarket := RebaseMarket("1", 3, &mockMarket)

		// expect pair a's rates not to have changed since it's based in "1" already
		expectedPairA := m.Pair{
			BaseId:  "1",
			QuoteId: "2",
			ExchangeMarkets: []m.ExchangeMarket{
				{
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
			ExchangeMarkets: []m.ExchangeMarket{
				{
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
			ExchangeMarkets: []m.ExchangeMarket{
				{
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

		So(actualMarket, ShouldResemble, &expectedMarket)
	})
	Convey("doesn't change rates when there is no path to rebase id", t, func() {
		mockPairA := m.Pair{
			BaseId:  "1",
			QuoteId: "2",
			ExchangeMarkets: []m.ExchangeMarket{
				{
					CurrentBid: 2,
					CurrentAsk: 2,
					BaseVolume: 1,
				},
			},
		}
		mockPairB := m.Pair{
			BaseId:  "3",
			QuoteId: "4",
			ExchangeMarkets: []m.ExchangeMarket{
				{
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

		actualMarket := RebaseMarket("1", 2, &mockMarket)

		// expect pair a's rates not to have changed since it's based in "1" already
		expectedPairA := m.Pair{
			BaseId:  "1",
			QuoteId: "2",
			ExchangeMarkets: []m.ExchangeMarket{
				{
					CurrentBid: 2,
					CurrentAsk: 2,
					BaseVolume: 1,
				},
			},
		}

		// expect pair b's rates to be 0 because they can't be rebased
		expectedPairB := m.Pair{
			BaseId:  "3",
			QuoteId: "4",
			ExchangeMarkets: []m.ExchangeMarket{
				{
					CurrentBid: 0,
					CurrentAsk: 0,
					BaseVolume: 0,
				},
			},
		}

		expectedMarket := m.Market{
			PairsById: map[string]m.Pair{
				expectedPairA.Id(): expectedPairA,
				expectedPairB.Id(): expectedPairB,
			},
		}

		So(actualMarket, ShouldResemble, &expectedMarket)
	})
}
