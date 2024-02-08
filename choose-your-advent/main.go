package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

func readJSONFile(file string) []byte {
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("something went wrong reading file")
		os.Exit(1)
	}

	return []byte(data)
}

func parseJSON(filepath string) map[string]Chapter {
	jsonByt := readJSONFile(filepath)
	var story map[string]Chapter

	if err := json.Unmarshal(jsonByt, &story); err != nil {
		panic(err)
	}
	return story
}

func ReqestHandler(fallback http.Handler, template *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimLeft(r.URL.Path, "/")
		story := parseJSON("./story.json")

		if data, ok := story[path]; ok {
			tmplt, _ := template.ParseFiles("story.html")
			err := tmplt.Execute(w, data)
			if err != nil {
				return
			}
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

func main() {
	mux := defaultMux()
	var tmplt *template.Template

	http.ListenAndServe(":8080", ReqestHandler(mux, tmplt))
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", NotFound)
	return mux
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Not Found!")
}
