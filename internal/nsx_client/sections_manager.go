package nsx_client

import (
	"encoding/json"
	"fmt"
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
	LastModifiedTime int64  `json:"_last_modified_time"`
	DisplayName      string `json:"display_name"`
	Tags             []Tags `json:"tags"`
	CreateTime       int64  `json:"_create_time"`
	ID               string `json:"id"`
}

func (c *Client) GetSgSections(debug bool, log Logger) ([]Section, error) {

	var sections 				[]Section
	var returnSections 	[]Section
	
	initialEndpoint := "/api/v1/firewall/sections?page_size=500"
	endpoint := initialEndpoint

	// handle pagination
	for {
		if debug {log.Printf("Requesting: " + endpoint)}
		var response 	SectionResponse
		respBody, err := c.makeGetRequest(endpoint)
		if err != nil {return sections, err}
		
		err = json.Unmarshal(respBody, &response)
		if err != nil {return sections, err}
		
		sections = append(sections, response.Results...)

		if response.Cursor == "" {
			break
		}
		endpoint = initialEndpoint + "&cursor=" + response.Cursor
	}

	// collate only sections that contain "ncp/cf_asg_name" tag
	for _, section := range(sections) {
		_, tagerr := findTag("ncp/cf_asg_name", section.DisplayName, section.Tags)
		if tagerr == nil {
			returnSections = append(returnSections, section)
		}
	}

	log.Printf("Collected " + strconv.Itoa(len(returnSections))  + " CF tagged firewall sections")

	return returnSections, nil
}

func findTag(scope, section string, tags []Tags) (string, error) {
	for _, tag := range(tags) {
		if tag.Scope == scope {
			return tag.Tag, nil
		}
	}
	return "", fmt.Errorf("tag scope [%s] not found for section %s", scope, section)
}