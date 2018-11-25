package urlshort

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var urlmap = map[string]string{
	"/ow": "https://openwichita.org",
	"/se": "https://seth.computer",
}

func fallbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "fallback")
}

func TestMapHandler(t *testing.T) {
	path := "/ow"

	r := httptest.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()

	handler, err := MapHandler(urlmap, http.HandlerFunc(fallbackHandler))
	assert.Nil(t, err, "error creating handler")

	handler(w, r)
	w.Flush()
	result := w.Result()

	assert.Equal(t, 302, result.StatusCode, "status code mismatch")
	assert.Equal(t, urlmap[path], result.Header.Get("location"), "incorect redirect url")
}

func TestMapHandlerFallback(t *testing.T) {
	path := "/does-not-exist"

	r := httptest.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()

	handler, err := MapHandler(urlmap, http.HandlerFunc(fallbackHandler))
	assert.Nil(t, err, "error creating handler")

	handler(w, r)
	w.Flush()
	result := w.Result()

	assert.Equal(t, 200, result.StatusCode, "status code mismatch")

	body, err := ioutil.ReadAll(result.Body)
	assert.Nil(t, err, "error reading response body")
	assert.Equal(t, "fallback", string(body), "incorrect response body")
}
