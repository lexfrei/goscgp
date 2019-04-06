package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/lexfrei/goscgp/parser"
)

func main() {
	siteURL, err := url.Parse("http://www.starcitygames.com/results?&numpage=25&switch_display=1&name=Skewer+the+Critics")
	if err != nil {
		os.Exit(1)
	}

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
