package main

import (
	"fmt"

	m "github.com/jochenboesmans/go-rebase/model/market"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest() string {
	pair := m.Pair{
		BaseId: "BLA",
		QuoteId: "BLE",
	}
	return fmt.Sprintf("pair id: %s", pair.Id())
}

func main() {
	lambda.Start(handleRequest)
}
