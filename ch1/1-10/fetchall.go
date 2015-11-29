package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	var outFilename string
	flag.StringVar(&outFilename, "out", "stdout",
		"The name of the file to print results to, defaulting to standard output.")
	flag.Parse()
	out := os.Stdout
	var err error
	if outFilename != "stdout" && outFilename != "" {
		out, err = os.Create(outFilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating %s: %v", outFilename, err)
			os.Exit(1)
		}
		defer out.Close()
	}
	start := time.Now()
	ch := make(chan string)
	for _, url := range flag.Args() {
		go fetch(url, ch)
	}
	for range flag.Args() {
		fmt.Fprintln(out, <-ch)
	}
	fmt.Fprintf(out, "%.2f seconds elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	defer resp.Body.Close()
	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		ch <- fmt.Sprintf("Error while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2f seconds, %7d bytes, %s", secs, nbytes, url)
}
