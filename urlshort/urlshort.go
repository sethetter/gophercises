package urlshort

import (
	"net/http"
	"strings"

	"github.com/go-yaml/yaml"
)

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var parsed map[string]string
	err := yaml.Unmarshal(yml, &parsed)
	if err != nil {
		return nil, err
	}
	return MapHandler(parsed, fallback)
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
