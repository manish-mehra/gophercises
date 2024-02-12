package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	port := flag.Int("port", 3000, "server port")
	filename := flag.String("file", "story.json", "JSON file with story")
	flag.Parse()

	file, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}
	story, err := JSONStory(file)
	if err != nil {
		panic(err)
	}

	content, err := os.ReadFile("./story.html")
	if err != nil {
		panic(err)
	}

	htmlString := string(content)

	tpl := template.Must(template.New("").Parse(htmlString))
	h := NewHandler(
		story,
		WithTemplate(tpl),
		WithPathFunc(pathFn),
	)
	mux := http.NewServeMux()
	mux.Handle("/story/", h)
	mux.Handle("/", NewHandler(story))
	fmt.Printf("Starting the server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}

func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}
	return path[len("/story/"):]
}
