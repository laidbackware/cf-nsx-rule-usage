package nsx_client

import (
	
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

func (c *Client) GetSgSections(api, user, password string) ([]Section, error) {

	var response SectionResponse
	var sections []Section

	baseEndpoint := "/api/v1/search/query"
	baseQuery := "query=resource_type:(FirewallSection)%20AND%20(tags.scope:%22ncp/cf_asg_name%22)"
	baseUrl := fmt.Sprintf("https://%s%s?%s", api, baseEndpoint, baseQuery)
	
	req, err := http.NewRequest("GET", baseUrl, nil)
	if err != nil {return sections, err}

	authstr := b64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, password)))

	// Add request headers
	req.Header = http.Header{
		"accept":        {"application/json"},
		"authorization": {"Basic " + authstr},
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {return sections, err}
	if resp.StatusCode != 200 {return sections, errors.New("Return code: " + strconv.Itoa(resp.StatusCode))}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {return sections, err}
  err = json.Unmarshal(respBody, &response)
	if err != nil {return sections, err}
	defer resp.Body.Close()

	sections = append(sections, response.Results...)

	for strInt(response.Cursor) < response.ResultCount {
		req.URL.RawQuery = baseQuery + "&cursor=" + response.Cursor
		resp, err := http.DefaultClient.Do(req)
		if err != nil {return sections, err}
		if resp.StatusCode != 200 {return sections, errors.New("Return code: " + strconv.Itoa(resp.StatusCode))}
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {return sections, err}
		err = json.Unmarshal(respBody, &response)
		if err != nil {return sections, err}

		defer resp.Body.Close()

		sections = append(sections, response.Results...)
	}

	return sections, nil
}

func strInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}