package urlshort

import (
	"net/http"
)

// func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
// 	// parse yaml to map
// 	return MapHandler(urlmap, fallback)
// }

func MapHandler(urlmap map[string]string, fallback http.Handler) (http.HandlerFunc, error) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if to, ok := urlmap[r.URL.Path]; ok {
			http.Redirect(w, r, to, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	})
	return handler, nil
}
