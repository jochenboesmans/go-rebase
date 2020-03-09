package main

import (
	"github.com/jochenboesmans/go-rebase/rebasing"

	"github.com/aws/aws-lambda-go/lambda"
	m "github.com/jochenboesmans/go-rebase/model/market"
)

type inputType struct {
	RebaseId      string   `json:"rebaseId"`
	MaxPathLength uint8    `json:"maxPathLength"`
	Market        []m.Pair `json:"market"`
}

type outputType struct {
	RebaseId string   `json:"rebaseId"`
	Market   []m.Pair `json:"market"`
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

func toOutputType(rebasedMarket m.Market, rebaseId string) outputType {
	output := outputType{
		RebaseId: rebaseId,
		Market:   []m.Pair{},
	}
	for _, pair := range rebasedMarket.PairsById {
		output.Market = append(output.Market, pair)
	}
	return output
}

func rebase(i inputType) (outputType, error) {
	market := i.extractMarket()
	rebasedMarket := *rebasing.RebaseMarket(i.RebaseId, i.MaxPathLength, &market)
	output := toOutputType(rebasedMarket, i.RebaseId)
	return output, nil
}

func main() {
	lambda.Start(rebase)
}
