package nsx_client

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupClient(t *testing.T) {
	l := log.New(os.Stderr, "", 0)
	_, err := SetupClient(mustEnv(t, "NSX_API"), mustEnv(t, "NSX_USER"), mustEnv(t, "NSX_PASS"), true, false, l)
	assert.Nil(t, err)
}

func TestCheckConnectivity(t *testing.T){
	err := CheckConnectivity(mustEnv(t, "NSX_API"))
	assert.Nil(t, err)

}

func TestRequest(t *testing.T){
	l := log.New(os.Stderr, "", 0)
	c, err := SetupClient(mustEnv(t, "NSX_API"), mustEnv(t, "NSX_USER"), mustEnv(t, "NSX_PASS"), true, false, l) 
	assert.Nil(t, err)
	respBody, err := c.makeGetRequest("/api/v1/firewall/sections")
	assert.Nil(t, err)
	assert.NotEmpty(t, respBody)
}

func mustEnv(t *testing.T, k string) string {
	t.Helper()

	if v, ok := os.LookupEnv(k); ok {
		return v
	}

	t.Fatalf("expected environment variable %q", k)
	return ""
}