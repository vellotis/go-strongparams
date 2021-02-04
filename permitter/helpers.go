package permitter

import "strings"

func noQuotes(val string) string {
	valLen := len(val)
	if valLen >= 2 && val[0] == '\'' && val[valLen-1] == '\'' {
		return val[1 : valLen-1]
	}
	return val
}

func removeProcessedRule(rules, toRemove string) string {
	return strings.Replace(rules, toRemove, "", 1)
}
