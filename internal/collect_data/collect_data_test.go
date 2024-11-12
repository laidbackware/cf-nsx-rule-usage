package collect_data

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/laidbackware/cf-nsx-rule-usage/internal/nsx_client"
)

func TestProcessSections(t *testing.T){
	client, err := nsx_client.SetupClient(mustEnv(t, "NSX_API"), mustEnv(t, "NSX_USER"), mustEnv(t, "NSX_PASS"))
	assert.Nil(t, err)
	sections, err := client.GetSgSections()
	assert.Nil(t, err)
	assert.NotEmpty(t, sections)
	rulesUsage, err := processSections(client, sections)
	assert.Nil(t, err)
	assert.NotEmpty(t, rulesUsage)
}

func mustEnv(t *testing.T, k string) string {
	t.Helper()

	if v, ok := os.LookupEnv(k); ok {
		return v
	}

	t.Fatalf("expected environment variable %q", k)
	return ""
}
