//go:generate ./xmlrpcgen --file $GOFILE SearchService
package core

import (
	"strings"

	"github.com/jinzhu/gorm"
)

/*
SearchResult item
*/
type SearchResult struct {
	_pypi_ordering int
	version        string
	name           string
	summary        string
}

/*
SearchService xml rpc service for searching packages
*/
type SearchService struct {
	Config Config
}

/*
search xml rpc method
*/
func (s *SearchService) search(query string) (result []SearchResult, err error) {
	query = NormalizePackageName(strings.TrimSpace(query))

	packages := []Package{}

	if err = s.Config.DB().Preload("Versions", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	}).Find(&packages, "name LIKE ?", "%"+query).Error; err != nil {
		return
	}

	result = []SearchResult{}

	// prepare result
	for _, pack := range packages {
		for _, version := range pack.Versions {
			_ = version
			result = append(result, SearchResult{
				name:           pack.Name,
				version:        version.Version,
				summary:        version.Description,
				_pypi_ordering: 999,
			})
		}
	}

	return
}
