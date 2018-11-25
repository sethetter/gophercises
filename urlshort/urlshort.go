package urlshort

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-yaml/yaml"
	bolt "go.etcd.io/bbolt"
)

var BucketName = []byte("urls")

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var parsed map[string]string
	err := yaml.Unmarshal(yml, &parsed)
	if err != nil {
		return nil, err
	}
	return MapHandler(parsed, fallback)
}

func JSONHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var parsed map[string]string
	err := json.Unmarshal(yml, &parsed)
	if err != nil {
		return nil, err
	}
	return MapHandler(parsed, fallback)
}

// TODO: make bucket name configurable
func DBHandler(db *bolt.DB, fallback http.Handler) (http.HandlerFunc, error) {
	urlmap := make(map[string]string)

	if err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			return nil
		}
		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			urlmap[string(k)] = string(v)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return MapHandler(urlmap, fallback)
}

func MapHandler(urlmap map[string]string, fallback http.Handler) (http.HandlerFunc, error) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if to, ok := urlmap[path]; ok {
			http.Redirect(w, r, to, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	})
	return handler, nil
}
