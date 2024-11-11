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

func (c *Client) GetSectionsStats(api, user, password, sectionId string) ([]SectionStats, error) {
	var sectionStats []SectionStats
	url := fmt.Sprintf("https://%s/api/v1/firewall/sections/%s/rules/stats", api, sectionId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {return sectionStats, err}

	authstr := b64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, password)))

	// Add request headers
	req.Header = http.Header{
		"accept":        {"application/json"},
		"authorization": {"Basic " + authstr},
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {return sectionStats, err}
	if resp.StatusCode != 200 {return sectionStats, errors.New("Return code: " + strconv.Itoa(resp.StatusCode))}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {return sectionStats, err}
  err = json.Unmarshal(respBody, &sectionStats)
	if err != nil {return sectionStats, err}
	defer resp.Body.Close()

	return sectionStats, nil

}
