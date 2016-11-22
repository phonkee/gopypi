package metadata

import "strings"

/*
CleanMethod cleans http method name
*/
func cleanMethod(method string) string {
	return strings.ToUpper(strings.TrimSpace(method))
}

// taken from encoding/json
// tagOptions is the string following a comma in a struct field's "json"
// tag, or the empty string. It does not include the leading comma.
type tagOptions string

// parseTag splits a struct field's json tag into its name and
// comma-separated options.
func ParseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		t := strings.TrimSpace(tag[idx+1:])
		return tag[:idx], tagOptions(t)
	}
	return tag, tagOptions("")
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}

/*
stringListContains returns whether string is in stringlist
 */
func stringListContains(where []string, what string) bool {
	for _, value := range(where) {
		if value == what {
			return true
		}
	}
	return false
}