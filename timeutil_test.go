package core

import (
	"testing"
	"time"
)

const layout = "Jan 2 2006 15:04:05"

func cmp(t *testing.T, in, out string, fn func(time.Time) time.Time) {
	parsedIn, errIn := time.Parse(layout, in)
	if errIn != nil {
		t.Errorf("parse error: %v", errIn.Error())
	}

	parsedOut, errOut := time.Parse(layout, out)
	if errOut != nil {
		t.Errorf("parse error: %v", errOut.Error())
	}

	result := fn(parsedIn)
	if !result.Equal(parsedOut) {
		t.Errorf("fn(%v) returned %v and not %v", parsedIn, result, parsedOut)
	}
}

func TestTimeAlignWeek(t *testing.T) {
	tc := []struct {
		in  string
		out string
	}{
		{"Nov 4 2016 15:04:05", "Oct 31 2016 00:00:00"},
	}

	for _, tt := range tc {
		t.Run("", func(st *testing.T) {
			cmp(t, tt.in, tt.out, TimeAlignWeek)
		})
	}
}

func TestTimeAlignMonth(t *testing.T) {
	tc := []struct {
		in  string
		out string
	}{
		{"Nov 4 2016 15:04:05", "Nov 1 2016 00:00:00"},
		{"Dec 31 2016 15:04:05", "Dec 1 2016 00:00:00"},
	}

	for _, tt := range tc {
		t.Run("", func(st *testing.T) {
			cmp(t, tt.in, tt.out, TimeAlignMonth)
		})
	}
}

func TestTimeAlignYear(t *testing.T) {
	tc := []struct {
		in  string
		out string
	}{
		{"Nov 4 2016 15:04:05", "Jan 1 2016 00:00:00"},
	}

	for _, tt := range tc {
		t.Run("", func(st *testing.T) {
			cmp(t, tt.in, tt.out, TimeAlignYear)
		})
	}
}
