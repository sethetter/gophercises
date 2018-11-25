package urlshort

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var urlmap = map[string]string{
	"/ow": "https://openwichita.org",
	"/se": "https://seth.computer",
}

func fallback(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "fallback")
}

func TestMapHandler(t *testing.T) {
	path := "/ow"

	r := httptest.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()

	handler, err := MapHandler(urlmap, http.HandlerFunc(fallback))
	if err != nil {
		t.Errorf("err: %v", err)
	}

	handler(w, r)
	w.Flush()
	result := w.Result()

	if result.StatusCode != 302 {
		t.Errorf("status code: expected %d, got %d", 302, result.StatusCode)
	}

	l := result.Header.Get("location")
	if l != urlmap[path] {
		t.Errorf("location header: expected %s, got %s", urlmap[path], l)
	}
}

func TestMapHandlerFallback(t *testing.T) {
	path := "/does-not-exist"

	r := httptest.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()

	handler, err := MapHandler(urlmap, http.HandlerFunc(fallback))
	if err != nil {
		t.Errorf("err: %v", err)
	}

	handler(w, r)
	w.Flush()
	result := w.Result()

	if result.StatusCode != 200 {
		t.Errorf("status code: expected %d, got %d", 302, result.StatusCode)
	}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Errorf("err: %v", err)
	}

	if string(body) != "fallback" {
		t.Errorf("fallback failure: expected %s, got %s", "fallback", string(body))
	}
}
