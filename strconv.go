package core

import "strconv"

/*
Atoui parses string and converts to uint, if error occures either 0 is returned or given default value (optional)
*/
func Atoui(value string, def ...uint) uint {

	// Get default value
	var defVal uint
	if len(def) > 0 {
		defVal = def[0]
	}

	var (
		err error
		u   uint64
	)
	if u, err = strconv.ParseUint(value, 10, 32); err != nil {
		return defVal
	}

	return uint(u)
}
