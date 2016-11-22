package paginator

import (
	"testing"

	"net/url"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPaginator_GetCount(t *testing.T) {

	Convey("Test set/get", t, func() {
		factory := NewFactory()

		p := factory()
		So(p.GetCount(), ShouldEqual, 0)
		p.Count(10)
		So(p.GetCount(), ShouldEqual, 10)
	})
}

func TestPaginator_FromURLValues(t *testing.T) {

	Convey("Test set/get", t, func() {
		factory := NewFactory()
		values := url.Values{
			DEFAULT_PAGE_PARAM:     []string{"130"},
			DEFAULT_PER_PAGE_PARAM: []string{"42"},
		}
		p := factory(values)
		So(p.GetPage(), ShouldEqual, 130)
		So(p.GetPerPage(), ShouldEqual, 42)
	})

	Convey("Test disable per page param", t, func() {
		factory := NewFactory(
			DisablePerPage(),
		)
		values := url.Values{
			DEFAULT_PAGE_PARAM:     []string{"130"},
			DEFAULT_PER_PAGE_PARAM: []string{"42"},
		}
		p := factory(values)
		So(p.GetPage(), ShouldEqual, 130)
		So(p.GetPerPage(), ShouldEqual, DEFAULT_PER_PAGE)
	})

	Convey("Test GetLimitOffset", t, func() {
		factory := NewFactory()

		var tests = []struct {
			page         int
			perpage      int
			add          bool
			expectlimit  int
			expectoffset int
		}{
			{2, 10, false, 10, 10},
			{5, 2, true, 10, 8},
		}

		for _, item := range tests {
			p := factory()
			p.Page(item.page)
			p.PerPage(item.perpage)
			limit, offset := p.GetLimitOffset(item.add)

			So(limit, ShouldEqual, item.expectlimit)
			So(offset, ShouldEqual, item.expectoffset)
		}

	})

}
