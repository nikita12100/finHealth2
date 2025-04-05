package main

import (
	"log"
	"test2/internal/common"
	"test2/internal/parser"
	"test2/internal/inserter"

	"github.com/xuri/excelize/v2"
)

const (
	brokerReportFile = "broker_report.xlsx"
	listName         = "broker_rep"
)

func main() {
	f, err := excelize.OpenFile(brokerReportFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rows, err := f.GetRows(listName)
	if err != nil {
		log.Fatal(err)
	}

	// Stats>
	operations := parser.FetchOperations(rows)
	count := parser.CalcCount(operations)
	countSorted := common.Sort(parser.CalcCount(operations))
	avgPrice := parser.CalcAvgPrice(operations)
	// <Stats

	// for k, v := range parser.CalcCount(operations) {
	// 	fmt.Printf("%v: %v, avgPrice=%.3f\n", k, v, avgPrice[k])
	// }

	inserter.InsertIntoSheet(count, avgPrice, countSorted)
}

// GMKN 100 -> 400 todo
