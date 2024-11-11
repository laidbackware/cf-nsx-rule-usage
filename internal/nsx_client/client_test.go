package nsx_client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckConnectivity(t *testing.T){
	err := CheckConnectivity("192.168.1.31")
	assert.Nil(t, err)

}