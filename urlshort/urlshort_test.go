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
	"ow": "https://openwichita.org",
	"se": "https://seth.computer",
}

func fallbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "fallback")
}

func TestMapHandler(t *testing.T) {
	handler, err := MapHandler(urlmap, http.HandlerFunc(fallbackHandler))
	assert.Nil(t, err, "error creating handler")
	testHandlerResponse(t, handler)
}

func TestMapHandlerFallback(t *testing.T) {
	handler, err := MapHandler(urlmap, http.HandlerFunc(fallbackHandler))
	assert.Nil(t, err, "error creating handler")
	testHandlerResponse(t, handler)
}

func TestYAMLHander(t *testing.T) {
	yml := `
ow: https://openwichita.org
se: https://seth.computer
`
	handler, err := YAMLHandler([]byte(yml), http.HandlerFunc(fallbackHandler))
	assert.Nil(t, err, "error creating handler")
	testHandlerResponse(t, handler)
}

func TestJSONHander(t *testing.T) {
	json := `
{
	"ow": "https://openwichita.org",
	"se": "https://seth.computer"
}
`
	handler, err := JSONHandler([]byte(json), http.HandlerFunc(fallbackHandler))
	assert.Nil(t, err, "error creating handler")
	testHandlerResponse(t, handler)
}
func testHandlerResponse(t *testing.T, handler http.Handler) {
	paths := []string{"ow", "se", "does-not-exist"}

	for _, path := range paths {
		r := httptest.NewRequest("GET", "/"+path, nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, r)
		w.Flush()
		result := w.Result()

		if to, ok := urlmap[path]; ok {
			assert.Equal(t, 302, result.StatusCode, "status code mismatch")
			assert.Equal(t, to, result.Header.Get("location"), "incorect redirect url")
		} else {
			assert.Equal(t, 200, result.StatusCode, "status code mismatch")
			body, err := ioutil.ReadAll(result.Body)
			assert.Nil(t, err, "error reading response body")
			assert.Equal(t, "fallback", string(body), "incorrect response body")
		}
	}
}
