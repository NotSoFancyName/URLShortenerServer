package handlers

import (
	"sync"
	"time"
)

var mtx sync.Mutex

func deleteURLEntryFunc(s string) func() {
	return func() {
		mtx.Lock()
		defer mtx.Unlock()
		delete(cachedURLs, s)
	}
}

func postponeURLEntryDeletion(e *longUrlEntry, s string) {
	e.expTimer.Stop()
	e.expTimer = time.AfterFunc(oneWeek, deleteURLEntryFunc(s))
}

func getCachedShortURL(l string) string {
	res := ""
	for k, v := range cachedURLs {
		if v.longUrl == l {
			res = k
			postponeURLEntryDeletion(&v, k)
			break
		}
	}

	return res
}
