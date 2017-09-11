package shell

import (
	"encoding/json"
	"os"

	"github.com/envkey/envkey-fetch/fetch"
)

func Source(envkey string, force bool, options fetch.FetchOptions) string {
	fetchedJson := fetch.Fetch(envkey, options)

	if fetchedJson == "" {
		return "echo 'error: ENVKEY invalid.'"
	}

	var resMap map[string]string
	err := json.Unmarshal([]byte(fetchedJson), &resMap)

	if err != nil {
		return "echo 'error: There was a problem parsing EnvKey's response.'"
	}

	if len(resMap) == 0 {
		return "echo 'No vars set'"
	}

	res := "export"

	for k, v := range resMap {
		var val string
		if !force && os.Getenv(k) != "" {
			val = os.Getenv(k)
		} else {
			val = v
		}

		res = res + " " + k + "=" + val
	}

	return res
}
