package render_output

import (
	"testing"

	"github.com/laidbackware/cf-nsx-rule-usage/internal/collect_data"
	"github.com/stretchr/testify/assert"
)

func TestBuildTableArray(t *testing.T) {
	sheetContents := map[string]map[string][]collect_data.Rule{
		"f1": {
			"asg1":{
				collect_data.Rule{
					Target: 			"1.1.1.1",
					Ports: 				"443",
					Protocol: 		"TCP",
					Created: 			"1",
					LastUpdated: 	"2",
					HitCount: 		0,
					PacketCount: 	0,
				},
				collect_data.Rule{
					Target: 			"1.1.1.2",
					Ports: 				"80",
					Protocol: 		"TCP",
					Created: 			"1",
					LastUpdated: 	"2",
					HitCount: 		0,
					PacketCount: 	0,
				},
			},
			"asg2":{
				collect_data.Rule{
					Target: 			"1.1.1.1",
					Ports: 				"443",
					Protocol: 		"TCP",
					Created: 			"1",
					LastUpdated: 	"2",
					HitCount: 		0,
					PacketCount: 	0,
				},
			},
		},
		"f2": {
			"asg1":{
				collect_data.Rule{
					Target: 			"1.1.1.1",
					Ports: 				"443",
					Protocol: 		"TCP",
					Created: 			"1",
					LastUpdated: 	"2",
					HitCount: 		0,
					PacketCount: 	0,
				},
			},
		},
	}
	tableArray := buildTableArray(sheetContents)
	assert.Equal(t, len(tableArray), 4)
}
