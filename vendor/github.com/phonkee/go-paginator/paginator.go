/*
paginator imlementation for golang
*/
package paginator

import (
	"net/http"

	"fmt"
	"net/url"

	"context"

	"strconv"
	"strings"

	"math"

	"github.com/phonkee/go-signals"
	"github.com/jinzhu/gorm"
)

/*
Paginator interface
*/
type Paginator interface {
	Count(int) Paginator
	GetCount() int

	Page(int) Paginator
	GetPage() int

	PerPage(int) Paginator
	GetPerPage() int

	// Normalize normalizes paginator (runs all options on it)
	Normalize() Paginator

	// return num of pages
	GetNumPages() int

	// read data from request, db
	From(...interface{}) Paginator
	FromDB(DB *gorm.DB, normalize ...bool) Paginator
	//FromJSON([]byte) Paginator
	FromRequest(request *http.Request, normalize ...bool) Paginator
	FromURLValues(values url.Values, normalize ...bool) Paginator

	// signal to be called when `From` with request is called or `FromRequest`
	OnRequest(signals.Receiver, ...string)

	PerPageParam(string) Paginator
	PageParam(string) Paginator

	GetLimitOffset(add ...bool) (int, int)

	// json marshal interface
	MarshalJSON() ([]byte, error)
	//// json unmarshal interface
	//UnmarshalJSON(data []byte) error

	UpdateURLValues(values url.Values) Paginator
}

/*
defaultPaginator returns default instantiated paginator
*/
func defaultPaginator(options ...Option) Paginator {
	return &paginator{
		options:      options[:],
		perPage:      DEFAULT_PER_PAGE,
		onRequest:    signals.New(),
		perPageParam: DEFAULT_PER_PAGE_PARAM,
		pageParam:    DEFAULT_PAGE_PARAM,
		page:         DEFAULT_PAGE,
	}
}

/*
Paginator imlementation
*/
type paginator struct {
	options []Option

	count   int
	perPage int
	page    int

	perPageParam string
	pageParam    string

	onRequest signals.Signal
}

/*
Set count of results
*/
func (p *paginator) Count(count int) Paginator {
	p.count = count
	return p
}

/*
Returns count of results
*/
func (p *paginator) GetCount() int {
	return p.count
}

/*
From can get information from request, db
*/
func (p *paginator) From(from ...interface{}) Paginator {
	for _, item := range from {

		// decide by type switch which method to call
		switch item := item.(type) {
		case *gorm.DB:
			p.FromDB(item)
		case *http.Request:
			p.FromRequest(item)
		case url.Values:
			p.FromURLValues(item)
		default:
			if DEBUG {
				panic(fmt.Sprintf("paginator: unsupported from value %+v", item))
			}
		}
	}

	// normalize after all operations if from was given
	if len(from) > 0 {
		p.Normalize()
	}

	return p
}

/*
FromDB calls `Count` method on given gorm.DB, so you need `queryset`
*/
func (p *paginator) FromDB(db *gorm.DB, normalize ...bool) Paginator {
	count := 0
	if err := db.Count(&count).Error; err == nil {
		p.Count(count)
	}
	if len(normalize) > 0 && normalize[0] {
		p.Normalize()
	}
	return p
}

/*
FromRequest reads request and sets GET values
*/
func (p *paginator) FromRequest(request *http.Request, normalize ...bool) Paginator {
	// @TODO: read request

	// dispatch onRequest signal
	ctx := context.WithValue(context.Background(), CONTEXT_ON_REQUEST_REQUEST, request)
	p.onRequest.Dispatch(ctx, true)

	if len(normalize) > 0 && normalize[0] {
		p.Normalize()
	}
	return p
}

/*
Reads url.Values and sets values to paginator
*/
func (p *paginator) FromURLValues(values url.Values, normalize ...bool) Paginator {

	if p.pageParam != "" {
		value := strings.TrimSpace(values.Get(p.pageParam))

		if iv, err := strconv.Atoi(value); err == nil {
			p.Page(iv)
		}
	}

	if p.perPageParam != "" {
		value := strings.TrimSpace(values.Get(p.perPageParam))
		if iv, err := strconv.Atoi(value); err == nil {
			p.PerPage(iv)
		}
	}

	if len(normalize) > 0 && normalize[0] {
		p.Normalize()
	}

	return p
}

/*
GetLimitOffset return limit offset numbers to be passed to db, if add is set offset will be added to limit
*/
func (p *paginator) GetLimitOffset(add ...bool) (limit int, offset int) {
	limit = p.perPage
	offset = (p.page - 1) * p.perPage

	if IsEnabled(add) {
		limit += offset
	}

	return
}

func (p *paginator) PageParam(param string) Paginator {
	p.pageParam = strings.TrimSpace(param)
	return p
}

func (p *paginator) Page(page int) Paginator {
	p.page = page
	return p
}

func (p *paginator) GetPage() int {
	return p.page
}

func (p *paginator) PerPage(perPage int) Paginator {
	p.perPage = perPage
	return p
}

func (p *paginator) GetPerPage() int {
	return p.perPage
}

func (p *paginator) GetNumPages() int {
	return int(math.Ceil(float64(p.GetCount()) / float64(p.GetPerPage())))
}

func (p *paginator) PerPageParam(param string) Paginator {
	p.perPageParam = strings.TrimSpace(param)
	return p
}

func appendKeyValueInt(source []byte, name string, value int) (target []byte) {
	target = source
	if len(source) > 1 {
		target = append(target, ", "...)
	}
	target = strconv.AppendQuote(target, name)
	target = append(target, ": "...)
	target = strconv.AppendInt(target, int64(value), 10)
	return
}

/*
custom (faster version) of json marshal
*/
func (p *paginator) MarshalJSON() (result []byte, err error) {
	result = []byte{}
	result = append(result, "{"...)
	if p.pageParam != "" {
		result = appendKeyValueInt(result, p.pageParam, p.page)
	}
	if p.perPageParam != "" {
		result = appendKeyValueInt(result, p.perPageParam, p.perPage)
	}
	result = appendKeyValueInt(result, "count", p.GetCount())
	result = appendKeyValueInt(result, "num_pages", p.GetNumPages())
	result = append(result, "}"...)
	return
}

/*
Normalize calls all options to re-set values
*/
func (p *paginator) Normalize() Paginator {
	runOptions(p, false, p.options...)
	return p
}

/*
Connect receiver to onRequest signal
*/
func (p *paginator) OnRequest(receiver signals.Receiver, id ...string) {
	p.onRequest.Connect(receiver, id...)
}

/*
UpdateURLValues update url vlues with paginator data, this will be reusable for MarshalJSON
*/
func (p *paginator) UpdateURLValues(values url.Values) Paginator {

	return p
}
