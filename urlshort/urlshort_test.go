package urlshort_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/sethetter/gophercises/urlshort"
	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
)

var urlmap = map[string]string{
	"ow": "https://openwichita.org",
	"se": "https://seth.computer",
}

func fallbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "fallback")
}

func TestMapHandler(t *testing.T) {
	handler, err := urlshort.MapHandler(urlmap, http.HandlerFunc(fallbackHandler))
	assert.Nil(t, err, "error creating handler")
	testHandlerResponse(t, handler)
}

func TestMapHandlerFallback(t *testing.T) {
	handler, err := urlshort.MapHandler(urlmap, http.HandlerFunc(fallbackHandler))
	assert.Nil(t, err, "error creating handler")
	testHandlerResponse(t, handler)
}

func TestYAMLHander(t *testing.T) {
	yml := `
ow: https://openwichita.org
se: https://seth.computer
`
	handler, err := urlshort.YAMLHandler([]byte(yml), http.HandlerFunc(fallbackHandler))
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
	handler, err := urlshort.JSONHandler([]byte(json), http.HandlerFunc(fallbackHandler))
	assert.Nil(t, err, "error creating handler")
	testHandlerResponse(t, handler)
}

func TestDBHandler(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	defer os.Remove(db.Path())

	handler, err := urlshort.DBHandler(db, http.HandlerFunc(fallbackHandler))
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

func testDB(t *testing.T) *bolt.DB {
	t.Helper()

	f, err := ioutil.TempFile("", "")
	assert.Nil(t, err, "error creating temp file for test db")

	path := f.Name()
	f.Close()
	os.Remove(path)

	db, err := bolt.Open(path, 0600, nil)
	assert.Nil(t, err, "error creating test db")

	seedDB(t, db)

	return db
}

func seedDB(t *testing.T, db *bolt.DB) {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(urlshort.BucketName)
		assert.Nil(t, err)

		for k, v := range urlmap {
			err := bucket.Put([]byte(k), []byte(v))
			if err != nil {
				return err
			}
		}

		return nil
	})

	assert.Nil(t, err)
}
