package shell

import (
	"encoding/json"
	"os"
	"sort"
	"strings"

	"github.com/envkey/envkey-fetch/fetch"
)

func Source(envkey string, force bool, options fetch.FetchOptions, pamCompatible bool) string {
	if envkey == "" {
		return "echo 'error: ENVKEY missing.'; false"
	}

	fetchRes, err := fetch.Fetch(envkey, options)

	if err != nil {
		return "echo 'error: " + err.Error() + "'; false"
	}

	if fetchRes == "" {
		return "echo 'error: ENVKEY invalid.'; false"
	}

	var resMap map[string]string
	err = json.Unmarshal([]byte(fetchRes), &resMap)

	if err != nil {
		return "echo 'error: There was a problem parsing EnvKey's response.'; false"
	}

	if len(resMap) == 0 {
		return "echo 'No vars set'"
	}

	var res string
	if pamCompatible {
		res = ""
	} else {
		res = "export"
	}

	var keys []string
	for k := range resMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		v := resMap[k]
		var key, val string

		if pamCompatible {
			// Remove newlines. Leave quotes alone.
			key = strings.Replace(k, "\n", "", -1)
		} else {
			// Quote quotes.
			key = strings.Replace(k, "'", `'"'"'`, -1)
		}

		if !force && os.Getenv(k) != "" {
			val = os.Getenv(k)
		} else {
			if pamCompatible {
				// Remove newlines. Leave quotes alone.
				val = strings.Replace(v, "\n", "", -1)
			} else {
				// Quote quotes.
				val = strings.Replace(v, "'", `'"'"'`, -1)
			}
		}

		if pamCompatible {
			if i > 0 {
				res = res + "\n"
			}
			// Do not quote keys, but quote values.
			res = res + "export " + key + "='" + val + "'"
		} else {
			// Quote both keys and values.
			res = res + " '" + key + "'='" + val + "'"
		}
	}

	return res
}
