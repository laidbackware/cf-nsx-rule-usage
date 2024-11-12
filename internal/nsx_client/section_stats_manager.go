package nsx_client

import (
	"encoding/json"
	"fmt"
)

type SectionStats struct {
	Results []RuleStats `json:"results"`
}
type RuleStats struct {
	RuleID                    string `json:"rule_id"`
	PacketCount               int    `json:"packet_count"`
	ByteCount                 int    `json:"byte_count"`
	SessionCount              int    `json:"session_count"`
	HitCount                  int    `json:"hit_count"`
	L7AcceptCount             int    `json:"l7_accept_count"`
	L7RejectCount             int    `json:"l7_reject_count"`
	L7RejectWithResponseCount int    `json:"l7_reject_with_response_count"`
	PopularityIndex           int    `json:"popularity_index"`
	MaxPopularityIndex        int    `json:"max_popularity_index"`
	MaxSessionCount           int    `json:"max_session_count"`
	TotalSessionCount         int    `json:"total_session_count"`
	Schema                    string `json:"_schema"`
}

func (c *Client) GetSectionStats(sectionId string) (SectionStats, error) {
	var sectionStats SectionStats

	endpoint := fmt.Sprintf("/api/v1/firewall/sections/%s/rules/stats", sectionId)
	respBody, err := c.makeGetRequest(endpoint)
	if err != nil {return sectionStats, err}

  err = json.Unmarshal(respBody, &sectionStats)
	if err != nil {return sectionStats, err}

	return sectionStats, nil
}
