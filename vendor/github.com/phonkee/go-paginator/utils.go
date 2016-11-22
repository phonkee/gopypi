package paginator

/*
IsEnabled is function that returns value of optional switches in functions

Example:

	func (normalize...bool) {
		isset := IsEnabled(normalize)
	}
 */
func IsEnabled(options []bool) bool {
	if len(options) == 0 {
		return false
	}

	return options[0]
}
