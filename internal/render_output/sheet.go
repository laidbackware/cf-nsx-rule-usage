package render_output

import (
	"fmt"
	"sort"

	"github.com/laidbackware/cf-nsx-rule-usage/internal/collect_data"
	"github.com/xuri/excelize/v2"
	"golang.org/x/exp/maps"
)

type sheetContent struct {
	sheetName    string
	sheetHeaders []excelize.Cell
	columnWidths []float64
	tableData    [][]excelize.Cell
}

func WriteSheet(ruleUsage collect_data.RuleUsage, unusedMonths int, outputFile string) (err error) {
	headers := []excelize.Cell{
		{Value: "Foundation Name"},
		{Value: "ASG Name"},
		{Value: "Target"},
		{Value: "Ports"},
		{Value: "Protocol"},
		{Value: "ASG Created"},
		{Value: "ASG Last Updated"},
		{Value: "Hit Count"},
		{Value: "Packet Count"},
	}
	widths := []float64{25, 30, 20, 20, 10, 20, 20, 15, 15}
	allSheetsContents := []sheetContent{
		{
			sheetName:    "unused_rules",
			sheetHeaders: headers,
			columnWidths: widths,
			tableData:    buildTableArray(ruleUsage.UnusedRules),
		},
		{
			sheetName: 		"all_rules",
			sheetHeaders: headers,
			columnWidths: widths,
			tableData:    buildTableArray(ruleUsage.AllRules),
		},
		{
			sheetName: 		fmt.Sprintf("unused_in_%d_months", unusedMonths),
			sheetHeaders: headers,
			columnWidths: widths,
			tableData:    buildTableArray(ruleUsage.UnusedRulesMonths),
		},
		{
			sheetName: 		"ununsed_by_asg",
			sheetHeaders: []excelize.Cell{
				{Value: "Foundation Name"},
				{Value: "ASG Name"},
				{Value: "Unused rules"},
				{Value: "Total rules"},
				{Value: "Used percentage"},
			},
			columnWidths: []float64{25, 30, 15, 15, 15},
			tableData:    buildTableArrayUnusedCount(ruleUsage.AllRules),
		},
	}

	err = renderSheet(allSheetsContents, outputFile)
	return
}

func renderSheet(allSheetsContents []sheetContent, outputFile string) (error) {
	// Initialize file
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			return
		}
	}()

	// Walk sheets and render
	for _, sheetContents := range allSheetsContents {
		_, err := f.NewSheet(sheetContents.sheetName)
		if err != nil {return err}
		setColumnWidths(f, sheetContents.sheetName, sheetContents.columnWidths)

		// Write headers
		writeLine(f, sheetContents.sheetName, sheetContents.sheetHeaders, 0)

		for row, line := range sheetContents.tableData {
			writeLine(f, sheetContents.sheetName, line, row + 1)
		}

	}
	_ = f.DeleteSheet("Sheet1")
	err := f.SaveAs(outputFile)
	return err

}

func buildTableArray(sourceData map[string]map[string][]collect_data.Rule) (tableArray [][]excelize.Cell) {
	// Get keys to enable sorting so that each sheet has a predictable order
	foundations := maps.Keys(sourceData)
	sort.Strings(foundations)
	for _, foundationName := range foundations {
		asgs := maps.Keys(sourceData[foundationName])
		sort.Strings(asgs)
		for _, asgName := range asgs {
			for _, rule := range sourceData[foundationName][asgName] {
				tableArray = append(tableArray, []excelize.Cell{
					{Value: foundationName},
					{Value: asgName},
					{Value: rule.Target},
					{Value: rule.Ports},
					{Value: rule.Protocol},
					{Value: rule.Created},
					{Value: rule.LastUpdated},
					{Value: rule.HitCount},
					{Value: rule.PacketCount},
				})
			}
		}
	}
	return
}

func buildTableArrayUnusedCount(sourceData map[string]map[string][]collect_data.Rule) (tableArray [][]excelize.Cell) {
	// Get keys to enable sorting so that each sheet has a predictable order
	foundations := maps.Keys(sourceData)
	sort.Strings(foundations)
	rowIdx := 1
	for _, foundationName := range foundations {
		asgs := maps.Keys(sourceData[foundationName])
		sort.Strings(asgs)
		for _, asgName := range asgs {
			tableArray = append(tableArray, []excelize.Cell{
				{Value: foundationName},
				{Value: asgName},
				{Formula: fmt.Sprintf("COUNTIFS(unused_rules!A:A,%s,unused_rules!B:B,%s)", getCellName(1, rowIdx+1), getCellName(2, rowIdx+1))},
				{Formula: fmt.Sprintf("COUNTIFS(all_rules!A:A,%s,all_rules!B:B,%s)", getCellName(1, rowIdx+1), getCellName(2, rowIdx+1))},
				{Formula: fmt.Sprintf("ROUND(%s/%s,2)*100", getCellName(3, rowIdx+1), getCellName(4, rowIdx+1))}})
			rowIdx++
		}
	}
	return
}

func getCellName(columnIdx, rolIdx int) (cellName string) {
	cellName, _ = excelize.CoordinatesToCellName(columnIdx, rolIdx)
	return
}

// Set widths of colums
func setColumnWidths(f *excelize.File, sheetName string, columnWidths []float64) {
	for columnIdx, columnWidth := range columnWidths {
		columnName, _ := excelize.ColumnNumberToName(columnIdx + 1)
		_ = f.SetColWidth(sheetName, columnName, columnName, columnWidth)
	}
}

// Write line to a worksheet based on an array of strings
func writeLine(f *excelize.File, sheetName string, content []excelize.Cell, rowIdx int) {
	for columnIdx, cellContent := range content {
		cellName, _ := excelize.CoordinatesToCellName(columnIdx+1, rowIdx+1)
		if cellContent.Value != nil {
			_ = f.SetCellValue(sheetName, cellName, cellContent.Value)
		} else {
			_ = f.SetCellFormula(sheetName, cellName, cellContent.Formula)
		}
	}
}
