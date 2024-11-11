package collect_data

import (
	"fmt"
	"github.com/laidbackware/cf-nsx-rule-usage/internal/nsx_client"
)

type Rule struct {
	target string
	ports string
	protocol string
	lastUpdated string
	hitCount string
	packetCount string

}

type RuleUsage struct {
	// map[foundation_name]map[asg_name]
	UnusedRules map[string]map[string][]Rule
	AllRules		map[string]map[string][]Rule

}

func Run(nsxApi, nsxUsername, nsxPassword string, skipVerify bool) error {
	client, err := nsx_client.SetupClient(nsxApi, nsxUsername, nsxPassword)
	if err != nil {return err}
	sections, err := client.GetSgSections(nsxApi, nsxUsername, nsxPassword)
	fmt.Println(len(sections))

	return nil
}

// func processSections(sections []nsx_client.Section) {
// 	for _, section := range(sections) {
// 		// Set start time
// 		// Get rules
// 		// Get stats
// 		// Iterate stats by index
// 			// Add rule to all array
// 			// If unused add rule to unused array
// 		// Sleep for rest of 0.011 seconds if necessary
// 	}
// }