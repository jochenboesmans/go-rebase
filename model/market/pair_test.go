package market

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var mockPair = Pair{
	ExchangeMarketDataByExchangeId: map[string]ExchangeMarketData{
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

func TestPair_BaseNeighborIds(t *testing.T) {
	mockPairA := Pair{
		BaseId: "1",
	}
	Convey("a pair B with a matching quote id to a pair A's base id should be a base neighbor of pair A", t, func() {
		mockPairBMatching := Pair{
			QuoteId: "1",
		}
		mockMarket := Market{
			PairsById: map[string]Pair{
				mockPairA.Id(): mockPairA,
				mockPairBMatching.Id(): mockPairBMatching,
			},
		}

		expected := []string{mockPairBMatching.Id()}

		actual := mockPairA.BaseNeighborIds(&mockMarket)

		So(actual, ShouldResemble, expected)
	})
	Convey("a pair B with a non-matching quote id to a pair A's base id should not be a base neighbor of pair A", t, func() {
		mockPairBNonMatching := Pair{
			QuoteId: "2",
		}
		mockMarket := Market{
			PairsById: map[string]Pair{
				mockPairA.Id(): mockPairA,
				mockPairBNonMatching.Id(): mockPairBNonMatching,
			},
		}

		var expected []string

		actual := mockPairA.BaseNeighborIds(&mockMarket)

		So(actual, ShouldResemble, expected)
	})
}

func TestPair_QuoteNeighborIds(t *testing.T) {
	mockPairA := Pair{
		QuoteId: "1",
	}
	Convey("a pair B with a matching base id to a pair A's quote id should be a quote neighbor of pair A", t, func() {
		mockPairBMatching := Pair{
			BaseId: "1",
		}
		mockMarket := Market{
			PairsById: map[string]Pair{
				mockPairA.Id(): mockPairA,
				mockPairBMatching.Id(): mockPairBMatching,
			},
		}

		expected := []string{mockPairBMatching.Id()}

		actual := mockPairA.BaseNeighborIds(&mockMarket)

		So(actual, ShouldResemble, expected)
	})
	Convey("a pair B with a non-matching base id to a pair A's quote id should not be a quote neighbor of pair A", t, func() {
		mockPairBNonMatching := Pair{
			BaseId: "2",
		}
		mockMarket := Market{
			PairsById: map[string]Pair{
				mockPairA.Id(): mockPairA,
				mockPairBNonMatching.Id(): mockPairBNonMatching,
			},
		}

		var expected []string

		actual := mockPairA.QuoteNeighborIds(&mockMarket)

		So(actual, ShouldResemble, expected)
	})
}

func TestCombinedBaseVolume(t *testing.T) {
	Convey("works as expected for basic mock pair", t, func() {
		expected := mockPair.ExchangeMarketDataByExchangeId["KYBER"].BaseVolume +
			mockPair.ExchangeMarketDataByExchangeId["UNISWAP"].BaseVolume

		result := mockPair.CombinedBaseVolume()

		So(result, ShouldEqual, expected)
	})
}

func TestBaseVolumeWeightedCurrentBidSum(t *testing.T) {
	Convey("works as expected for basic mock pair", t, func() {
		expected := mockPair.ExchangeMarketDataByExchangeId["KYBER"].CurrentBid*
			mockPair.ExchangeMarketDataByExchangeId["KYBER"].BaseVolume +
			mockPair.ExchangeMarketDataByExchangeId["UNISWAP"].CurrentBid*
				mockPair.ExchangeMarketDataByExchangeId["UNISWAP"].BaseVolume

		result := mockPair.BaseVolumeWeightedCurrentBidSum()

		So(result, ShouldEqual, expected)
	})
}

func TestBaseVolumeWeightedCurrentAskSum(t *testing.T) {
	Convey("works as expected for basic mock pair", t, func() {
		expected := mockPair.ExchangeMarketDataByExchangeId["KYBER"].CurrentAsk*
			mockPair.ExchangeMarketDataByExchangeId["KYBER"].BaseVolume +
			mockPair.ExchangeMarketDataByExchangeId["UNISWAP"].CurrentAsk*
				mockPair.ExchangeMarketDataByExchangeId["UNISWAP"].BaseVolume

		result := mockPair.BaseVolumeWeightedCurrentAskSum()

		So(result, ShouldEqual, expected)
	})
}

func TestBaseVolumeWeightedSpreadAverage(t *testing.T) {
	Convey("works as expected for basic mock pair", t, func() {
		expected := ((mockPair.BaseVolumeWeightedCurrentBidSum() +
			mockPair.BaseVolumeWeightedCurrentAskSum()) /
			2) /
			mockPair.CombinedBaseVolume()

		result, err := mockPair.BaseVolumeWeightedSpreadAverage()

		So(err, ShouldBeNil)
		So(result, ShouldEqual, expected)
	})
}
