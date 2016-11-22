/*
Filter maintains various filters from http forms (get/post)
*/
package core

import (
	"net/http"

	"github.com/gorilla/schema"
	"github.com/jinzhu/gorm"
)

/*
BindFilter binds filter to given url values
*/
func BindFilter(r *http.Request, filter Filter) (err error) {
	decoder := schema.NewDecoder()
	if err = decoder.Decode(filter, r.URL.Query()); err != nil {
		return
	}
	return
}

/*
Filter is interface to filter database querysets
*/
type Filter interface {

	// returns list of filter funcs
	Apply(queryset *gorm.DB) *gorm.DB
}

/*
NewUserListFilter returns new UserListFilter
*/
func NewUserListFilter(r *http.Request) Filter {
	return UserListFilter{
		IsActive: StringParseBool(r.URL.Query().Get("is_active")),
	}
}

/*
UserListFilter filters users from url
*/
type UserListFilter struct {
	IsActive *bool
}

/*
Apply applies filter to queryset
*/
func (u UserListFilter) Apply(queryset *gorm.DB) *gorm.DB {
	if u.IsActive != nil {
		queryset = ApplyFilterFuncs(queryset, FFWhere("is_active = ?", u.IsActive))
	}
	return queryset
}
