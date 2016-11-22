/*
timeutil is set of functions to work with time.Time such as: alignment, stripping of time values etc..
*/
package core

import "time"

/*
TimeStripTime strips time part from time
*/
func TimeStripTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

/*
TimeAlignWeek aligns date to first day of week and strips time and timezone
*/
func TimeAlignWeek(t time.Time) time.Time {
	return TimeStripTime(t).AddDate(0, 0, -int(t.Day()))
}

/*
TimeAlignWeek aligns date to first day of month and strips time and timezone
*/
func TimeAlignMonth(t time.Time) time.Time {
	return TimeStripTime(t).AddDate(0, 0, -t.Day()+1)
}

/*
TimeAlignWeek aligns date to first day of year and strips time and timezone
*/
func TimeAlignYear(t time.Time) time.Time {
	return TimeAlignMonth(t).AddDate(0, -int(t.Month())+1, 0)
}
