package firebase_test

import (
	"path/filepath"
	"runtime"
)

// BasePath returns the catalog of the module this function resides in
func BasePath() string {
	_, b, _, _ := runtime.Caller(0)
	base := filepath.Dir(b)
	return base
}
