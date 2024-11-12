package nsx_client

import (
	
	"encoding/json"
	"strconv"
)

type SectionResponse struct {
	Results     []Section `json:"results"`
	ResultCount int       `json:"result_count"`
	Cursor      string    `json:"cursor"`
}
type Tags struct {
	Scope string `json:"scope"`
	Tag   string `json:"tag"`
}
type Section struct {
	Description      string `json:"description"`
	LastModifiedTime int64  `json:"_last_modified_time"`
	DisplayName      string `json:"display_name"`
	Tags             []Tags `json:"tags"`
	CreateTime       int64  `json:"_create_time"`
	ID               string `json:"id"`
	Locked           bool   `json:"locked"`
	RuleCount        int    `json:"rule_count"`
}

func (c *Client) GetSgSections() ([]Section, error) {

	var response SectionResponse
	var sections []Section

	endpoint := "/api/v1/search/query?query=resource_type:(FirewallSection)%20AND%20(tags.scope:%22ncp/cf_asg_name%22)"
	respBody, err := c.makeGetRequest(endpoint)
	if err != nil {return sections, err}

  err = json.Unmarshal(respBody, &response)
	if err != nil {return sections, err}

	sections = append(sections, response.Results...)

	for strInt(response.Cursor) < response.ResultCount {
		respBody, err := c.makeGetRequest(endpoint + "&cursor=" + response.Cursor)
		if err != nil {return sections, err}

		err = json.Unmarshal(respBody, &response)
		if err != nil {return sections, err}

		sections = append(sections, response.Results...)
	}
	return sections, nil
}

func strInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}