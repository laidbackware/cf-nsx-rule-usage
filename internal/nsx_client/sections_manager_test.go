package nsx_client

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSgSections(t *testing.T) {
	l := log.New(os.Stderr, "", 0)
	client, err := SetupClient(mustEnv(t, "NSX_API"), mustEnv(t, "NSX_USER"), mustEnv(t, "NSX_PASS"), true, false, l)
	assert.Nil(t, err)
	sections, err := client.GetSgSections(false, l)
	assert.Nil(t, err)
	assert.NotEmpty(t, sections)
}