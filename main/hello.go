package main

import (
	"fmt"
	m "github.com/jochenboesmans/go-rebase/model"
)

func main() {
	pair := m.Pair{
		BaseId: "BLA",
		QuoteId: "BLE",
	}
	fmt.Printf("pair id: %s", pair.Id())
}
