package shell

import (
	"encoding/json"
	"os"
	"sort"
	"strings"

	"github.com/envkey/envkey-fetch/fetch"
)

func Source(envkey string, force bool, options fetch.FetchOptions) string {
	fetchRes := fetch.Fetch(envkey, options)

	if strings.HasPrefix(fetchRes, "error: ") {
		return "echo '" + fetchRes + "'"
	}

	if fetchRes == "" {
		return "echo 'error: ENVKEY invalid.'"
	}

	var resMap map[string]string
	err := json.Unmarshal([]byte(fetchRes), &resMap)

	if err != nil {
		return "echo 'error: There was a problem parsing EnvKey's response.'"
	}

	if len(resMap) == 0 {
		return "echo 'No vars set'"
	}

	res := "export"

	var keys []string
	for k := range resMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := resMap[k]
		var val string
		if !force && os.Getenv(k) != "" {
			val = os.Getenv(k)
		} else {
			val = strings.Replace(v, "'", `'"'"'`, -1)
		}

		res = res + " " + k + "='" + val + "'"
	}

	return res
}
