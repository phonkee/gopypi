/*
Various database (gorm) helpers
*/
package core

import (
	"github.com/jinzhu/gorm"
	"github.com/phonkee/go-paginator"
)

/*
Initial setup of database connection
*/
func setupDB(DB *gorm.DB) {
	// remove original functions to update timestamps, we will maintain that by ourself in models in
	// BeforeCreate, BeforeSave methods
	DB.Callback().Create().Remove("gorm:update_time_stamp")
}

/*
LimitQueryset sets limits to queryset and returns it
*/
func LimitQueryset(db *gorm.DB, p paginator.Paginator) *gorm.DB {
	limit, offset := p.GetLimitOffset()
	return db.Limit(limit).Offset(offset)
}

/*
CountQueryset performs count query and sets paginator count
*/
func CountQueryset(db *gorm.DB, p paginator.Paginator) *gorm.DB {
	count := 0
	db.Limit(0).Offset(0).Count(&(count))
	p.Count(count)
	return db
}
