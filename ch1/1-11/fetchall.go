package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	args := struct {
		out string
		in  string
	}{}
	flag.StringVar(&args.out, "out", "stdout",
		"The name of the file to print results to, defaulting to standard output.")
	flag.StringVar(&args.in, "in", "",
		"The name of a file of newline separated URLs to fetch.")
	flag.Parse()
	out := os.Stdout
	var err error
	if args.out != "stdout" && args.out != "" {
		out, err = os.Create(args.out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating %s: %v", args.out, err)
			os.Exit(1)
		}
		defer out.Close()
	}
	URLs := flag.Args()
	if args.in != "" {
		b, err := ioutil.ReadFile(args.in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading %s: %v", args.in, err)
			os.Exit(1)
		}
		bb := bytes.Split(b, []byte("\n"))
		URLs = make([]string, len(bb))
		for i, u := range bb {
			URLs[i] = string(u)
		}
	}
	start := time.Now()
	ch := make(chan string)
	for _, url := range URLs {
		go fetch(url, ch)
	}
	for range URLs {
		fmt.Fprintln(out, <-ch)
	}
	fmt.Fprintf(out, "%.2f seconds elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	s := time.Now()
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
	secs := time.Since(s).Seconds()
	ch <- fmt.Sprintf("%.2f seconds, %7d bytes, %s", secs, nbytes, url)
}
