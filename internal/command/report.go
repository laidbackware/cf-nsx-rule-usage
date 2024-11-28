package command

import (
	"os"
	"path/filepath"

	"github.com/laidbackware/cf-nsx-rule-usage/internal/collect_data"
	"github.com/laidbackware/cf-nsx-rule-usage/internal/render_output"
)

func generateReport(nsxApi, nsxUsername, nsxPassword, outputType, outputFile string, skipVerify, debug bool, log Logger) {

	if outputType != "json" && outputType != "xlsx" {
		log.Fatalf("Requested output format '%s'  is invalid. Please use: [json, xlsx]", outputType)
		os.Exit(1)
	}

	if outputFile == "" {
		currentDir, err := os.Getwd()
		handleError(err)
		if outputType == "json" {
			outputFile = filepath.Join(currentDir, "report.json")
		} else {
			outputFile = filepath.Join(currentDir, "report.xlsx")
		}
	}

	log.Printf("Collecting data from NSX...\n")
	ruleUsage, err := collect_data.CollectData(nsxApi, nsxUsername, nsxPassword, skipVerify, debug, log)
	handleError(err)

	switch outputType {
	case "xlsx":
		log.Printf("Writing to spreadsheet...\n")
		handleError(render_output.WriteSheet(ruleUsage, outputFile))
	case "json":
		handleError(render_output.WriteJSON(ruleUsage, outputFile))
	}
	log.Printf("Written file: %s\n", outputFile)
}
