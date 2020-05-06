package redis

import (
	"net/http"
	"time"

	"github.com/cjtoolkit/ignition/shared/utility/httpError"
)

func checkIfModifiedSince(r *http.Request, modtime time.Time) {
	if r.Method != "GET" && r.Method != "HEAD" {
		return
	}
	ims := r.Header.Get("If-Modified-Since")
	if ims == "" || modtime.IsZero() {
		return
	}
	t, err := http.ParseTime(ims)
	if err != nil {
		return
	}

	// The Data-Modified header truncates sub-second precision, so
	// use mtime < t+1s instead of mtime <= t to check for unmodified.
	if modtime.Before(t.Add(1 * time.Second)) {
		httpError.HaltNotModified()
	}
}
