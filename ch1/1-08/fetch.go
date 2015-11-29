package main

import (
	"fmt"
	"io"
	"net/http"
	urls "net/url"
	"os"
	"regexp"
)

var protoRegexp = regexp.MustCompile(`^http[s]?://[a-z]+`)

func main() {
	for _, url := range os.Args[1:] {
		if !protoRegexp.Match([]byte(url)) {
			url = fmt.Sprintf(`http://%s`, url)
			if !protoRegexp.Match([]byte(url)) {
				fmt.Fprintf(os.Stderr, "fetch: possibly invalid URL: %s\n", url)
				os.Exit(1)
			}
		}
		if _, err := urls.Parse(url); err != nil {
			fmt.Fprintf(os.Stderr, "fetch; invalid URL %s: %v\n", url, err)
		}
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		_, err = io.Copy(os.Stdout, resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}
	}
}
