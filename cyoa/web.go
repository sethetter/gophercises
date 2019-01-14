package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

func RunWeb(adventure Adventure) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := RenderArc(w, adventure, r.URL.Path); err != nil {
			switch err.Error() {
			case "Invalid arc key":
				http.NotFound(w, r)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	})

	fmt.Println("Listening on 8080")
	return http.ListenAndServe(":8080", nil)
}

func RenderArc(w http.ResponseWriter, adventure Adventure, key string) error {
	if key[0] == '/' {
		key = key[1:]
	}

	arc, ok := adventure.arcs[key]
	if !ok {
		return errors.New("Invalid arc key")
	}

	tmpl, err := template.New("arc").Parse(`
<!doctype html>
<html>
<head>
	<title>Adventure!</title>
</head>
<body>
	{{range $i, $p := .Story}}
		<p>{{$p}}</p>
	{{end}}
	<ul>
		{{range $j, $o := .Options}}
			<li><a href="/{{$o.Arc}}">{{$o.Text}}</a></li>
		{{end}}
	</ul>
</body>
</html>
	`)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, arc)
}
