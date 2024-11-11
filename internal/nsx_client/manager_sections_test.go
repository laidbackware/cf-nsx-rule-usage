package nsx_client

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testingClient = Client{
	HttpClient: &http.Client{Timeout: time.Duration(1) * time.Second},
}

func TestBuildTableArray(t *testing.T) {
	_, err := testingClient.GetSgSections("192.168.1.31", "admin", "VMware1!VMware1!!")
	assert.Nil(t, err)
}