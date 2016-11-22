package classy

import (
	"fmt"
	"strings"

	"github.com/serenize/snaker"
)

/*
IsEnabledOption returns whether wildcard option is enabled
*/
func isEnabledOption(option []bool) bool {
	return len(option) > 0 && option[0]
}

/*
GetStructName returns struct name in snake case
*/
func getStructName(target interface{}, snake ...bool) (result string) {
	fullname := fmt.Sprintf("%T", target)
	splitted := strings.Split(fullname, ".")
	sl := len(splitted)

	result = splitted[sl-1]

	if isEnabledOption(snake) {
		result = snaker.CamelToSnake(result)
	}

	return
}

/*
GetViewName returns structname with stripped "view"
*/
func getViewName(target interface{}, snake ...bool) (result string) {
	result = getStructName(target)

	suffixes := []string{"apiview", "apiviewset", "view", "viewset"}

	for _, s := range suffixes {
		lowered := strings.ToLower(result)
		if strings.HasSuffix(lowered, s) {
			result = result[:len(lowered)-len(s)]
			break
		}
	}

	if isEnabledOption(snake) {
		result = snaker.CamelToSnake(result)
	}

	return
}

/*
removeDuplicates removes duplicates in string list
 */
func removeDuplicates(input []string) (output []string) {
	output = make([]string, 0, len(input))

	got := map[string]struct{}{}

	for _, item := range input {
		if _, ok := got[item]; !ok {
			output = append(output, item)
			got[item] = struct{}{}
		}
	}
	return
}

/*
StringListContains returns whether string is in list
*/
func stringListContains(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

/*
join multiple paths
*/
func joinPaths(paths ...string) (result string) {
	result = ""
	for _, path := range paths {
		if strings.HasSuffix(result, "/") && strings.HasPrefix(path, "/") {
			result = strings.TrimRight(result, "/")
		}
		result += path
	}
	return
}

/*
makeName creates name from teplate. {name} is replaced in template. If not found in template
 */
func makeName(tpl string, name string) (result string) {
	result = name
	if strings.Contains(tpl, "{name}") {
		result = strings.Replace(tpl, "{name}", name, -1)
	} else {
		result = tpl + name
	}

	return
}
