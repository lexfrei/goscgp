package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/lexfrei/goscgp/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <card name>\n", os.Args[0])
		return
	}

	siteURL, err := url.Parse("http://www.starcitygames.com/results?&switch_display=1")
	if err != nil {
		os.Exit(1)
	}

	q := siteURL.Query()
	q.Set("name", strings.Join(os.Args[1:], " "))
	siteURL.RawQuery = q.Encode()

	c := &http.Client{}

	result, err := parser.DoRequest(*siteURL, c)
	if err != nil {
		os.Exit(1)
	}

	json, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		os.Exit(1)
	}

	fmt.Printf("%s\n", json)
}
