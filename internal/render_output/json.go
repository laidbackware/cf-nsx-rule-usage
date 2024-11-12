package render_output

import (
	"encoding/json"
	"os"

	"github.com/laidbackware/cf-nsx-rule-usage/internal/collect_data"
)

func WriteJSON(ruleUsage collect_data.RuleUsage, outputFile string) (err error) {

	outputJson, err := json.MarshalIndent(ruleUsage, "", "    ")
	if err != nil {
		return
	}
	err = os.WriteFile(outputFile, outputJson, 0644)
	return
}
