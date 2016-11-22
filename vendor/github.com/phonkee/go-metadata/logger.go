package metadata

import "github.com/phonkee/go-logger"

var (
	// global debug option
	debugEnabled = false
)

var (
	loggerInfo    = logger.Info("metadata")
	loggerError   = logger.Error("metadata")
	loggerDebug   = logger.Debug("metadata")
	loggerWarning = logger.Warning("metadata")
)


/*
Debug enables global debug for go-metadata
 */
func Debug() {
	debugEnabled = true
}