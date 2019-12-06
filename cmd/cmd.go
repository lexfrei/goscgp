package main

import (
	"fmt"
	"os"

	"github.com/lexfrei/goscgp/qgen"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <card name>\n", os.Args[0])
		return
	}

	u := qgen.NewQurl()
	u.QUrlWithName("negate")
	fmt.Println(u.URL.String())
}
