package core

import "github.com/phonkee/go-paginator"

var (
	// Prepare common paginator for all api
	CommonPaginator = paginator.NewFactory(
		paginator.PerPage(20, 30),
	)
)
