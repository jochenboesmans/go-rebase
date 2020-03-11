package main

import (
	"github.com/jochenboesmans/go-rebase/rebasing"

	"github.com/aws/aws-lambda-go/lambda"
	m "github.com/jochenboesmans/go-rebase/model/market"
)

type inputType struct {
	RebaseAssetId string   `json:"rebaseAssetId"`
	MaxPathLength uint8    `json:"maxPathLength"`
	Market        []m.Pair `json:"market"`
}

type outputType struct {
	RebaseAssetId string   `json:"rebaseAssetId"`
	Market        []m.Pair `json:"market"`
}

func (input inputType) extractMarket() m.Market {
	market := m.Market{
		PairsById: map[string]m.Pair{},
	}
	for _, pair := range input.Market {
		market.PairsById[pair.Id()] = pair
	}
	return market
}

func toOutputType(rebasedMarket m.Market, rebaseAssetId string) outputType {
	output := outputType{
		RebaseAssetId: rebaseAssetId,
		Market:        []m.Pair{},
	}
	for _, pair := range rebasedMarket.PairsById {
		output.Market = append(output.Market, pair)
	}
	return output
}

func rebase(i inputType) (outputType, error) {
	market := i.extractMarket()
	rebasedMarket := *rebasing.RebaseMarket(i.RebaseAssetId, i.MaxPathLength, &market)
	output := toOutputType(rebasedMarket, i.RebaseAssetId)
	return output, nil
}

func main() {
	lambda.Start(rebase)
}
