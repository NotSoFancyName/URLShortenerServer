// Package handlers implements two handlers
// First one provides form to enter long URL and handles provided short URLs 
// Second one provides the short URL for previosly entered long one
package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/NotSoFancyName/URLShortenerServer/persistance"
	"github.com/NotSoFancyName/URLShortenerServer/shortener"
)

const (
	//Form`s action name
	ActionName   = "/shortened_url/"
	textAreaName = "body"

	tmplPath = "./templates/*"
	idxTmplName = "index.html"
	empTmplName = "empty.html"
	urlTmplName = "shortenedurl.html"

	oneMinute = 60000000000
	oneHour   = 60 * oneMinute
	oneDay    = 24 * oneHour
	oneWeek   = 7 * oneDay
)

type longUrlEntry struct {
	longUrl  string
	expTimer *time.Timer
}

var idxParams = struct {
	Action       string
	TextAreaName string
}{
	ActionName,
	textAreaName,
}

var ( 
	_, base, _, _ = runtime.Caller(0)
	templates = template.Must(template.ParseGlob(filepath.Join(
		filepath.Dir(base),
		tmplPath)))
)

func init() {
	shortener.SetCounter(persistance.GetMostRecentUpdatedEntryID())
}

// Serves a simple index.html with a form to enter a long URL
// It also handles short URLs and redirects to their long counterparts  
func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) > 1 {
		if elem, ok := cachedURLs[r.URL.Path[1:]]; ok {
			postponeURLEntryDeletion(&elem, r.URL.Path[1:])
			http.Redirect(w, r, elem.longUrl, http.StatusFound)
		} else {
			if url := persistance.GetURLFromDB(r.URL.Path[1:], false); url == "" {
				http.NotFound(w, r)
			} else {
				http.Redirect(w, r, url, http.StatusFound)
			}
		}
	} else {
		templates.ExecuteTemplate(w, idxTmplName, idxParams)
	}
}

// Provides a short URL that corresponds the entered long one
// First it tries to find short URL in cachedURLs, if it fails it makes a 
// query to DB, if it fails again it creates a new one
func ShortenedURLHandler(w http.ResponseWriter, r *http.Request) {
	enteredURL := strings.TrimSpace(r.FormValue(textAreaName))
	if enteredURL == "" || strings.Contains(enteredURL, r.Host) {
		templates.ExecuteTemplate(w, empTmplName, struct{}{})
		return
	}

	shortURL := getCachedShortURL(enteredURL)

	if shortURL == "" {
		shortURL = persistance.GetURLFromDB(enteredURL, true)
		if shortURL != "" {
			cachedURLs[shortURL] = longUrlEntry{
				enteredURL,
				time.AfterFunc(oneWeek, deleteURLEntryFunc(shortURL))}
		}
	}

	if shortURL == "" {
		shortURL = shortener.ShortURLString()
		cachedURLs[shortURL] = longUrlEntry{
			enteredURL,
			time.AfterFunc(oneWeek, deleteURLEntryFunc(shortURL))}
		persistance.Save(shortURL, enteredURL)
	}

	templates.ExecuteTemplate(w, urlTmplName, struct {
		ShortenedURL string
	}{
		r.Host + "/" + shortURL,
	})
}
