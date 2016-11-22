package paginator

// Option for paginator factory
type Option func(p Paginator, first bool)

/*
runs options on paginator
 */
func runOptions(p Paginator, first bool, options...Option) {
	for _, option := range options {
		option(p, first)
	}
}

/*
OptPageParam sets page param name on factory, if blank value given, it's disabled
 */
func PageParam(param string) Option {
	return Option(func(p Paginator, first bool) {
		p.PageParam(param)
	})
}

/*
OptPerPageParam sets per page param name on factory, if blank value given, it's disabled
 */
func PerPageParam(param string) Option {
	return Option(func(p Paginator, first bool) {
		p.PerPageParam(param)
	})
}

func DisablePerPage() Option {
	return PerPageParam("")
}


/*
LimitChoices sets choices for limit
 */
func PerPage(limits... int) Option {

	// if no limits provided continue
	if len(limits) == 0 {
		limits = append(limits, DEFAULT_PER_PAGE)
	}

	return Option(func(p Paginator, first bool) {

		for _, limit := range limits {
			if limit == p.GetPerPage() {
				return
			}
		}

		// set first
		p.PerPage(limits[0])
	})
}

