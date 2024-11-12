package nsx_client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSgSections(t *testing.T) {
	client, err := SetupClient(mustEnv(t, "NSX_API"), mustEnv(t, "NSX_USER"), mustEnv(t, "NSX_PASS"))
	assert.Nil(t, err)
	sections, err := client.GetSgSections()
	assert.Nil(t, err)
	assert.NotEmpty(t, sections)
}