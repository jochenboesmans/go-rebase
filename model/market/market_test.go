package market

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMarket_RebaseNeighbors(t *testing.T) {
	Convey("adjacent pairs see each other as neighbors", t, func() {
		mockPairA := Pair{
			BaseAssetId:  "1",
			QuoteAssetId: "2",
		}
		mockPairB := Pair{
			BaseAssetId:  "2",
			QuoteAssetId: "3",
		}
		mockMarket := Market{
			PairsById: map[string]Pair{
				mockPairA.Id(): mockPairA,
				mockPairB.Id(): mockPairB,
			},
		}

		expectedRebaseNeighbors := map[string]Neighbors{
			mockPairA.Id(): {
				Base:  []string{},
				Quote: []string{mockPairB.Id()},
			},
			mockPairB.Id(): {
				Base:  []string{mockPairA.Id()},
				Quote: []string{},
			},
		}

		actualRebaseNeighbors := mockMarket.RebaseNeighbors()

		So(actualRebaseNeighbors, ShouldResemble, expectedRebaseNeighbors)
	})
	Convey("non-adjacent pairs don't see each other as neighbors", t, func() {
		mockPairA := Pair{
			BaseAssetId:  "1",
			QuoteAssetId: "2",
		}
		mockPairB := Pair{
			BaseAssetId:  "3",
			QuoteAssetId: "4",
		}
		mockMarket := Market{
			PairsById: map[string]Pair{
				mockPairA.Id(): mockPairA,
				mockPairB.Id(): mockPairB,
			},
		}

		expectedRebaseNeighbors := map[string]Neighbors{
			mockPairA.Id(): {
				Base:  []string{},
				Quote: []string{},
			},
			mockPairB.Id(): {
				Base:  []string{},
				Quote: []string{},
			},
		}

		actualRebaseNeighbors := mockMarket.RebaseNeighbors()

		So(actualRebaseNeighbors, ShouldResemble, expectedRebaseNeighbors)
	})
}
