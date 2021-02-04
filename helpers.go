package strongparams

import (
	"net/url"
	"strings"
)

func hasKey(queryValues url.Values, requireKey string) bool {
	for key := range queryValues {
		if key == requireKey || strings.HasPrefix(key, requireKey+"[") {
			return true
		}
	}
	return false
}

func cloneUrlValues(queryValues url.Values) url.Values {
	newQueryValues := make(url.Values)
	for key, value := range queryValues {
		newQueryValues[key] = make([]string, len(value))
		copy(newQueryValues[key], value)
	}
	return newQueryValues
}
