package shell_test

import (
	"testing"

	"github.com/envkey/envkey-fetch/fetch"
	"github.com/envkey/envkey-source/shell"
	"github.com/envkey/envkey-source/version"
	"github.com/stretchr/testify/assert"
)

const VALID_ENVKEY = "Emzt4BE7C23QtsC7gb1z-3NvfNiG1Boy6XH2o-env-staging.envkey.com"
const INVALID_ENVKEY = "Emzt4BE7C23QtsC7gb1z-3NvfNiG1Boy6XH2oinvalid-env-staging.envkey.com"

func TestSource(t *testing.T) {
	// Test valid
	validRes := shell.Source(VALID_ENVKEY, true, fetch.FetchOptions{false, "", "envkey-source", version.Version}, false)
	assert.Equal(t, "export 'TEST'='it' 'TEST_2'='works!' 'TEST_INJECTION'=''\"'\"'$(uname)' 'TEST_SINGLE_QUOTES'='this'\"'\"' is ok' 'TEST_SPACES'='it does work!'", validRes)

	// Test --pam-compatible
	validRes2 := shell.Source(VALID_ENVKEY, true, fetch.FetchOptions{false, "", "envkey-source", version.Version}, true)
	assert.Equal(t, "export TEST='it'\nexport TEST_2='works!'\nexport TEST_INJECTION=''$(uname)'\nexport TEST_SINGLE_QUOTES='this' is ok'\nexport TEST_SPACES='it does work!'", validRes2)


	// Test invalid
	invalidRes := shell.Source(INVALID_ENVKEY, true, fetch.FetchOptions{false, "", "envkey-source", version.Version}, false)
	assert.Equal(t, "echo 'error: ENVKEY invalid'; false", invalidRes)
}
