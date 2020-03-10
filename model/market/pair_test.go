package market

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var mockPair = Pair{
	ExchangeMarkets: []ExchangeMarket{
		{
			BaseVolume: 1.5,
			CurrentBid: 100.7,
			CurrentAsk: 103.5,
		},
		{
			BaseVolume: 3,
			CurrentBid: 150.1,
			CurrentAsk: 155.2,
		},
	},
}

func TestCombinedBaseVolume(t *testing.T) {
	Convey("works as expected for basic mock pair", t, func() {
		expected := mockPair.ExchangeMarkets[0].BaseVolume +
			mockPair.ExchangeMarkets[1].BaseVolume

		actual := mockPair.CombinedBaseVolume()

		So(actual, ShouldEqual, expected)
	})
}

func TestBaseVolumeWeightedCurrentBidSum(t *testing.T) {
	Convey("works as expected for basic mock pair", t, func() {
		expected := mockPair.ExchangeMarkets[0].CurrentBid*
			mockPair.ExchangeMarkets[0].BaseVolume +
			mockPair.ExchangeMarkets[1].CurrentBid*
				mockPair.ExchangeMarkets[1].BaseVolume

		actual := mockPair.BaseVolumeWeightedCurrentBidSum()

		So(actual, ShouldEqual, expected)
	})
}

func TestBaseVolumeWeightedCurrentAskSum(t *testing.T) {
	Convey("works as expected for basic mock pair", t, func() {
		expected := mockPair.ExchangeMarkets[0].CurrentAsk*
			mockPair.ExchangeMarkets[0].BaseVolume +
			mockPair.ExchangeMarkets[1].CurrentAsk*
				mockPair.ExchangeMarkets[1].BaseVolume

		actual := mockPair.BaseVolumeWeightedCurrentAskSum()

		So(actual, ShouldEqual, expected)
	})
}

func TestBaseVolumeWeightedSpreadAverage(t *testing.T) {
	Convey("works as expected for basic mock pair", t, func() {
		expected := ((mockPair.BaseVolumeWeightedCurrentBidSum() +
			mockPair.BaseVolumeWeightedCurrentAskSum()) /
			2) /
			mockPair.CombinedBaseVolume()

		actual := mockPair.BaseVolumeWeightedSpreadAverage()

		So(actual, ShouldEqual, expected)
	})
	Convey("combined base volume of pair is zero", t, func() {
		mockPair := Pair{
			ExchangeMarkets: []ExchangeMarket{
				{
					BaseVolume: float32(0),
				},
			},
		}
		expected := float32(0)
		actual := mockPair.BaseVolumeWeightedSpreadAverage()

		So(actual, ShouldEqual, expected)
	})
}
