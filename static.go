package core

import (
	"net/http"
	"os"
)

/*
FallbackFileSystem wraps http Filesystem which checks if requested file is not found, it falls back to given default.
This is useful for SPA with html5 mode so it opens any url in index.html
 */
type FallbackFileSystem struct {
	original http.FileSystem
	fallback string
}

/*
Open opens and reads file
 */
func (f *FallbackFileSystem) Open(name string) (result http.File, err error) {
	// Try to check if file exists
	if result, err = f.original.Open(name); err == os.ErrNotExist {

		// If not found point to index
		result, err = f.original.Open(f.fallback)
	}
	return
}
