package command

import (
	"os"
	"path/filepath"

	"github.com/laidbackware/cf-nsx-rule-usage/internal/collect_data"
	"github.com/laidbackware/cf-nsx-rule-usage/internal/render_output"
)

func generateReport(nsxApi, nsxUsername, nsxPassword, outputType, outputFile string, skipVerify bool, log Logger) {

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

	// var healthState collect_data.HealthState
	ruleUsage, err := collect_data.CollectData(nsxApi, nsxUsername, nsxPassword, skipVerify)
	handleError(err)

	switch outputType {
	case "xlsx":
		// handleError(render_output.WriteSheet(healthState, outputFile))
	case "json":
		handleError(render_output.WriteJSON(ruleUsage, outputFile))
	}
	log.Printf("Written file: %s\n", outputFile)
}
