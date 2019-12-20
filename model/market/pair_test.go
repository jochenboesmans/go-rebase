package market

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var mockPair = Pair{
	ExchangeMarketDataByExchangeId: map[ExchangeId]ExchangeMarketData{
		KYBER: {
			BaseVolume: 1.5,
			CurrentBid: 100.7,
			CurrentAsk: 103.5,
		},
		UNISWAP: {
			BaseVolume: 3,
			CurrentBid: 150.1,
			CurrentAsk: 155.2,
		},
	},
}

func TestCombinedBaseVolume(t *testing.T) {
	Convey("works as expected for basic mock pair", t, func() {
		expected := mockPair.ExchangeMarketDataByExchangeId[KYBER].BaseVolume +
			mockPair.ExchangeMarketDataByExchangeId[UNISWAP].BaseVolume

		result := mockPair.CombinedBaseVolume()

		So(result, ShouldEqual, expected)
	})
}

func TestBaseVolumeWeightedCurrentBidSum(t *testing.T) {
	Convey("works as expected for basic mock pair", t, func() {
		expected := mockPair.ExchangeMarketDataByExchangeId[KYBER].CurrentBid*
			mockPair.ExchangeMarketDataByExchangeId[KYBER].BaseVolume +
			mockPair.ExchangeMarketDataByExchangeId[UNISWAP].CurrentBid*
				mockPair.ExchangeMarketDataByExchangeId[UNISWAP].BaseVolume

		result := mockPair.BaseVolumeWeightedCurrentBidSum()

		So(result, ShouldEqual, expected)
	})
}

func TestBaseVolumeWeightedCurrentAskSum(t *testing.T) {
	Convey("works as expected for basic mock pair", t, func() {
		expected := mockPair.ExchangeMarketDataByExchangeId[KYBER].CurrentAsk*
			mockPair.ExchangeMarketDataByExchangeId[KYBER].BaseVolume +
			mockPair.ExchangeMarketDataByExchangeId[UNISWAP].CurrentAsk*
				mockPair.ExchangeMarketDataByExchangeId[UNISWAP].BaseVolume

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
