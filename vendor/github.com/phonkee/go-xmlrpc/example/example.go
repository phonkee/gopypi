//go:generate ./xmlrpcgen --file $GOFILE --debug SearchService
package example

import _ "github.com/beevik/etree"

/*
SearchService xml rpc service for searching in database
*/
type SearchService struct {
}

/*
Search service method
 */
func (h *SearchService) Search(query string, page int, isit bool) ([]string, error) {
    return []string{}, nil
}