package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
)

const (
	indexPath = "./templates/index.html"
	testUrl   = "UnexistingLink"
)

func TestDefaultHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DefaultHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v expected %v", status, http.StatusOK)
	}

	temp, err := template.ParseFiles(filepath.Join(
		filepath.Dir(base),
		indexPath))

	if err != nil {
		fmt.Println(err)
	}

	var buf bytes.Buffer
	if err := temp.Execute(&buf, idxParams); err != nil {
		fmt.Println(err)
	}

	expected := buf.String()

	if rr.Body.String() != expected {
		t.Errorf("handler invalid body: got %v expected %v",
			rr.Body.String(), expected)
	}

}

func TestShortenedURLHandler(t *testing.T) {
	f := url.Values{}
	f.Set(textAreaName, testUrl)
	req, err := http.NewRequest("POST", "http://localhost:8080"+ActionName,
		strings.NewReader(f.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	ShortenedURLHandler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if !strings.Contains(string(body), req.Host) {
		t.Errorf("Invalid host %v expected %v",
			string(body), req.Host)
	}
}
