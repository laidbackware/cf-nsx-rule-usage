package nsx_client

import (
	"encoding/json"
	"fmt"
)

type RulesResponse struct {
	Cursor      string    `json:"cursor"`
	SortBy      string    `json:"sort_by"`
	ResultCount int       `json:"result_count"`
	Results     []Rule 		`json:"results"`
}
type Destinations struct {
	TargetDisplayName string `json:"target_display_name"`
	IsValid           bool   `json:"is_valid"`
	TargetType        string `json:"target_type"`
	TargetID          string `json:"target_id"`
}
type Service struct {
	ResourceType 			string 		`json:"resource_type"`
	IcmpType     			int    		`json:"icmp_type"`
	Protocol     			string 		`json:"protocol"`
	L4Protocol     		string 		`json:"l4_protocol"`
	IcmpCode     			int    		`json:"icmp_code"`
	Destination_ports []string 	`json:"destination_ports"`
}
type Services struct {
	Service Service `json:"service"`
}
type Rule struct {
	ID                   string         `json:"id"`
	DisplayName          string         `json:"display_name"`
	Notes                string         `json:"notes"`
	DestinationsExcluded bool           `json:"destinations_excluded"`
	Destinations         []Destinations `json:"destinations,omitempty"`
	Services             []Services     `json:"services,omitempty"`
	IPProtocol           string         `json:"ip_protocol"`
	RuleTag              string         `json:"rule_tag"`
	Logged               bool           `json:"logged"`
	Action               string         `json:"action"`
	SourcesExcluded      bool           `json:"sources_excluded"`
	Disabled             bool           `json:"disabled"`
	Direction            string         `json:"direction"`
	Revision             int            `json:"_revision"`
}

func (c *Client) GetSectionRules(sectionId string) ([]Rule, error) {
	var response RulesResponse
	var rules []Rule

	endpoint := fmt.Sprintf("/api/v1/firewall/sections/%s/rules", sectionId)
	respBody, err := c.makeGetRequest(endpoint)
	if err != nil {return rules, err}

  err = json.Unmarshal(respBody, &response)
	if err != nil {return rules, err}

	rules = append(rules, response.Results...)

	for response.Cursor != "" {
		respBody, err := c.makeGetRequest(endpoint + "&cursor=" + response.Cursor)
		if err != nil {return rules, err}
		
		err = json.Unmarshal(respBody, &response)
		if err != nil {return rules, err}

		rules = append(rules, response.Results...)
	}

	return rules, nil
}
