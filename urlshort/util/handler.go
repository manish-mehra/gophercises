package util

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v3"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if url, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, url, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

type PathMap struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func ParseYaml(yamlBytes []byte) ([]PathMap, error) {

	var pathEntries []PathMap
	err := yaml.Unmarshal(yamlBytes, &pathEntries)
	if err != nil {
		fmt.Println("Error parsing YAML:", err)
		return pathEntries, err
	}
	return pathEntries, nil
}

func buildMap(parsedYaml []PathMap) map[string]string {
	bMap := make(map[string]string)

	for _, entry := range parsedYaml {
		bMap[entry.Path] = entry.URL
	}
	return bMap
}

func YAMLHandler(yaml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := ParseYaml(yaml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}
