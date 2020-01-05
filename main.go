package main

import (
	"github.com/jochenboesmans/go-rebase/rebasing"

	"github.com/aws/aws-lambda-go/lambda"
	m "github.com/jochenboesmans/go-rebase/model/market"
)

type inputType struct {
	RebaseId  string   `json:"rebaseId"`
	PathDepth uint8    `json:"pathDepth"`
	Market    m.Market `json:"market"`
}

func rebase(i inputType) (m.Market, error) {
	rebasing.RebaseMarket(i.RebaseId, i.PathDepth, &i.Market)
	return i.Market, nil
}

func main() {
	lambda.Start(rebase)
}
