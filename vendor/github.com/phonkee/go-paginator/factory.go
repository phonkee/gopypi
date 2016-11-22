package paginator

// Factory function that returns paging, from can be request of gorm dg
type Factory func(from ...interface{}) Paginator

/*
NewFactory returns function that produces instance of Paginator. options control the creation of Paginators.
If you provide any compatible `froms` when calling factory, they will be passed to `From` method.

 */
func NewFactory(options ...Option) Factory {

	return Factory(func(from ...interface{}) (result Paginator) {

		result = defaultPaginator(options...)

		// apply options
		runOptions(result, true, options...)

		if len(from) > 0 {
			result.From(from...)
		}

		return
	})
}
