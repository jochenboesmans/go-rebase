package main

import (
	"github.com/jochenboesmans/go-rebase/rebasing"

	"github.com/aws/aws-lambda-go/lambda"
	m "github.com/jochenboesmans/go-rebase/model/market"
)

func rebase(rebaseId string, pathDepth uint8, market m.Market) (m.Market, error) {
	rebasing.RebaseMarket(rebaseId, pathDepth, &market)
	return market, nil
}

func main() {
	lambda.Start(rebase)
}
