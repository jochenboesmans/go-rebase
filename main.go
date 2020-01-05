package main

import (
	"github.com/jochenboesmans/go-rebase/rebasing"

	"github.com/aws/aws-lambda-go/lambda"
	m "github.com/jochenboesmans/go-rebase/model/market"
)

type inputType struct {
	rebaseId  string
	pathDepth uint8
	market    m.Market
}

func rebase(i inputType) (m.Market, error) {
	rebasing.RebaseMarket(i.rebaseId, i.pathDepth, &i.market)
	return i.market, nil
}

func main() {
	lambda.Start(rebase)
}
