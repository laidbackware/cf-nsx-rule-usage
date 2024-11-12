package nsx_client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckConnectivity(t *testing.T){
	err := CheckConnectivity("192.168.1.31")
	assert.Nil(t, err)

}

func TestRequest(t *testing.T){
	c, err := SetupClient("192.168.1.31", "admin", "VMware1!VMware1!!") 
	assert.Nil(t, err)
	respBody, err := c.makeGetRequest("/api/v1/firewall/sections")
	assert.Nil(t, err)
	assert.NotEmpty(t, respBody)
}