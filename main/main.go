package main

import (
	"fmt"
	"time"

	"github.com/hyssa-dev/edgar"
)

func main() {
	company, err := edgar.NewFilingFetcher().CompanyFolder("TSLA", edgar.FilingType10Q)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(company.Filing(edgar.FilingType10Q, time.Date(2024, 10, 24, 0, 0, 0, 0, time.UTC)))
	// fmt.Println(company)

	// company, err := edgar.NewFilingFetcher().CreateFolder(strings.NewReader("TSLA"), edgar.FilingType10Q)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(company)

}
