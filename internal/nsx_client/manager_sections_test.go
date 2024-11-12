package nsx_client

import (
	"crypto/tls"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testingClient = Client{
	HttpClient: &http.Client{
		Timeout: time.Duration(5) * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	},
	BaseUrl: "https://192.168.1.31",
	Header: http.Header{
		"accept":        {"application/json"},
		"authorization": {"Basic YWRtaW46Vk13YXJlMSFWTXdhcmUxISE"},
	},
}

func TestGetSgSections(t *testing.T) {
	sections, err := testingClient.GetSgSections()
	assert.Nil(t, err)
	assert.NotEmpty(t, sections)
}