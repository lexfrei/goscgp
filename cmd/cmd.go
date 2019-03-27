package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/lexfrei/goscgp/parser"
)

func main() {
	siteURL, err := url.Parse("http://www.starcitygames.com/results?&numpage=25&switch_display=1&name=forest")
	if err != nil {
		os.Exit(1)
	}

	result, err := parser.QueryWalker(*siteURL)
	if err != nil {
		os.Exit(1)
	}

	fmt.Println(len(result))

	// json, err := json.Marshal(result)
	// if err != nil {
	// 	os.Exit(1)
	// }

	// fmt.Printf("%s\n", json)
}
