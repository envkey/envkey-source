package shell_test

import (
	"testing"

	"github.com/envkey/envkey-fetch/fetch"
	"github.com/envkey/envkey-source/shell"
	"github.com/stretchr/testify/assert"
)

const VALID_ENVKEY = "Emzt4BE7C23QtsC7gb1z-3NvfNiG1Boy6XH2o-env-staging.envkey.com"
const INVALID_ENVKEY = "Emzt4BE7C23QtsC7gb1z-3NvfNiG1Boy6XH2oinvalid-env-staging.envkey.com"

func TestSource(t *testing.T) {
	// Test valid
	validRes := shell.Source(VALID_ENVKEY, true, fetch.FetchOptions{false, ""})
	assert.Equal(t, "export TEST='it' TEST_2='works!' TEST_INJECTION=''\"'\"'$(uname)' TEST_SINGLE_QUOTES='this'\"'\"' is ok' TEST_SPACES='it does work!'", validRes)

	// Test invalid
	invalidRes := shell.Source(INVALID_ENVKEY, true, fetch.FetchOptions{false, ""})
	assert.Equal(t, "echo 'error: ENVKEY invalid'", invalidRes)
}
