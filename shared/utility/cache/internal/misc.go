package internal

import (
	"net/http"
	"time"

	ctx "github.com/cjtoolkit/ctx/v2"
	"github.com/cjtoolkit/ctx/v2/ctxHttp"
	"github.com/cjtoolkit/ignition/shared/utility/httpError"
)

func CheckIfModifiedSince(r *http.Request, modtime time.Time) {
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

func CheckModifiedTime(modifiedTime time.Time, context ctx.Context) {
	if !modifiedTime.IsZero() {
		ctxHttp.Response(context).Header().Set("Last-Modified", modifiedTime.UTC().Format(http.TimeFormat))
	}
}
