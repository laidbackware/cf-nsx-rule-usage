package render_output

import (
	"sort"
	"strconv"

	"github.com/laidbackware/cf-nsx-rule-usage/internal/collect_data"
	"github.com/xuri/excelize/v2"
	"golang.org/x/exp/maps"
)

type sheetContent struct {
	sheetName    string
	sheetHeaders []string
	columnWidths []float64
	tableData    [][]string
}

func WriteSheet(ruleUsage collect_data.RuleUsage, outputFile string) (err error) {
	headers := []string{
		"Foundation Name",
		"ASG Name",
		"Rule index",
		"Target",
		"Ports",
		"Protocol",
		"Created",
		"Last Updated",
		"Hit Count",
		"Packet Count",
	}
	widths := []float64{25, 30, 10, 20, 20, 10, 20, 20, 15, 15}
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

func buildTableArray(sourceData map[string]map[string][]collect_data.Rule) (tableArray [][]string) {
	// Get keys to enable sorting so that each sheet has a predictable order
	foundations := maps.Keys(sourceData)
	sort.Strings(foundations)
	idx := 0
	for _, foundationName := range foundations {
		asgs := maps.Keys(sourceData[foundationName])
		sort.Strings(asgs)
		for _, asgName := range asgs {
			for ruleIdx, rule := range sourceData[foundationName][asgName] {
				tableArray = append(tableArray, []string{
					foundationName,
					asgName,
					strconv.Itoa(ruleIdx),
					rule.Target,
					rule.Ports,
					rule.Protocol,
					rule.Created,
					rule.LastUpdated,
					strconv.Itoa(rule.HitCount),
					strconv.Itoa(rule.PacketCount),
				})
				idx++
			}
		}
	}
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
func writeLine(f *excelize.File, sheetName string, content []string, rowIdx int) {
	for columnIdx, cellContent := range content {
		cellName, _ := excelize.CoordinatesToCellName(columnIdx+1, rowIdx+1)
		_ = f.SetCellValue(sheetName, cellName, cellContent)
	}
}
