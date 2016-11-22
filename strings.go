/*
Various string helpers
 */
package core

import "strings"

/*
StringListContains returns whether string is in stringlist
 */
func StringListContains(where []string, what string) bool {
	for _, value := range(where) {
		if value == what {
			return true
		}
	}
	return false
}

/*
StringParseBool returns boolean value by given string value
 */
func StringParseBool(value string) (result *bool) {
	if value = strings.ToLower(strings.TrimSpace(value)); value == "" {
		return
	}
	switch value {
	case "true", "1", "on":
		tmp := true
		return &tmp
	case "false", "0", "off":
		tmp := false
		return &tmp
	}
	return
}