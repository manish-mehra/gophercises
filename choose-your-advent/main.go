package main

import (
	"encoding/json"
	"fmt"
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

func ReqestHandler(fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimLeft(r.URL.Path, "/")
		story := parseJSON("./story.json")

		if data, ok := story[path]; ok {
			fmt.Fprintf(w, data.Title)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

func main() {
	mux := defaultMux()

	http.ListenAndServe(":8080", ReqestHandler(mux))
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", NotFound)
	return mux
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Not Found!")
}
