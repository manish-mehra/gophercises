package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

// type loc struct {
// 	Value string `xml:"loc"`
// }

// type urlset struct {
// 	Urls  []loc  `xml:"url"`
// 	Xmlns string `xml:"xmlns,attr"`
// }

func main() {

	defaultLink := "https://gobyexample.com" // https://gobyexample.com
	rootLink := flag.String("link", defaultLink, "a flag for url")
	flag.Parse()

	// toXml := urlset{
	// 	Xmlns: xmlns,
	// }

	graph := BFS(rootLink)

	// for _, page := range graph {
	// 	for _, url := range page {
	// 		toXml.Urls = append(toXml.Urls, loc{url})
	// 	}
	// }

	// fmt.Print(xml.Header)
	// enc := xml.NewEncoder(os.Stdout)
	// enc.Indent("", "  ")
	// if err := enc.Encode(toXml); err != nil {
	// 	panic(err)
	// }
	// fmt.Println()
	PrintGraph(graph, *rootLink, 0)
}

func BFS(rootLink *string) map[string][]string {

	q := Queue{}
	visited := make(map[string]string)
	q.Enqueue(*rootLink)
	visited[*rootLink] = *rootLink

	graph := make(map[string][]string)

	base := BuildBaseURL(*rootLink)

	for !q.IsEmpty() {

		current, exists := q.Dequeue()

		if exists {
			resp, err := http.Get(current)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			link, err := Parse(resp.Body) // fetch all the href links from the page
			if err != nil {
				panic(err)
			}
			hrefs := FormatLinks(link, base)
			filteredLinks := FilterExternalLinks(hrefs, base)
			// enque all links
			for _, value := range filteredLinks {
				if _, isVisited := visited[value]; !isVisited {
					q.Enqueue(value)
					visited[value] = value
					graph[current] = append(graph[current], value)
				}
			}
		}
	}

	// PrintGraph(graph, *rootLink, 0)
	// indent := ""
	// for key, val := range graph {
	// 	fmt.Println(key)
	// 	for _, v := range val {
	// 		fmt.Println(indent, v)
	// 	}
	// 	indent += " "
	// }
	return graph
}

func PrintGraph(graph map[string][]string, node string, depth int) {
	indentation := strings.Repeat("\t", depth)
	fmt.Printf("%s%s\n", indentation, node)

	for _, neighbor := range graph[node] {
		PrintGraph(graph, neighbor, depth+1)
	}
}

func BuildBaseURL(rootLink string) string {
	resp, err := http.Get(rootLink)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// add base url
	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()
	return base
}

// format links i.e. add base domain, remove # or mailto
func FormatLinks(link []Link, base string) []string {

	var hrefs []string
	for _, l := range link {
		switch {
		case strings.HasPrefix(l.Href, "/"): // case: '/hello'
			hrefs = append(hrefs, base+l.Href)
		case strings.HasPrefix(l.Href, "http:") || strings.HasPrefix(l.Href, "https:"):
			hrefs = append(hrefs, l.Href)
		case strings.HasPrefix(l.Href, "") && !strings.HasPrefix(l.Href, "#") && !strings.HasPrefix(l.Href, "mailto:"):
			hrefs = append(hrefs, base+"/"+l.Href)
		}
	}
	return hrefs
}

// filter out external links
func FilterExternalLinks(hrefs []string, base string) map[string]string {
	urls := make(map[string]string)

	for _, href := range hrefs {
		isSameDomain, _ := sameDomain(base, href)
		if isSameDomain {
			if _, exists := urls[href]; !exists {
				urls[href] = href
			}
		}
	}
	return urls
}

func sameDomain(url1, url2 string) (bool, error) {
	// Parse the URLs
	parsedURL1, err := url.Parse(url1)
	if err != nil {
		return false, err
	}

	parsedURL2, err := url.Parse(url2)
	if err != nil {
		return false, err
	}

	// Extract Hostnames from the parsed URLs
	hostname1 := strings.ToLower(parsedURL1.Hostname())
	hostname2 := strings.ToLower(parsedURL2.Hostname())

	// Compare Hostnames to check if they are the same domain
	return hostname1 == hostname2, nil
}
