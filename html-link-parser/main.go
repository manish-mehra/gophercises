package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {

	file, err := os.ReadFile("./ex2.html")
	if err != nil {
		panic(err)
	}

	htmlString := string(file)

	r := strings.NewReader(htmlString)
	links, err := Parse(r)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", links)

}
