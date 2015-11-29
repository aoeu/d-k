package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	c := make(counts)
	files := os.Args[1:]
	if len(files) == 0 {
		c.countLines(os.Stdin)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			c.countLines(f)
			f.Close()
		}
	}
	for line, r := range c {
		if r.count > 1 {
			fmt.Printf("%d\t%s", r.count, line)
			for fileName, _ := range r.files {
				fmt.Printf("\t%s", fileName)
			}
			fmt.Printf("\n")
		}
	}
}

type counts map[string]struct {
	count int
	files map[string]bool
}

func (c counts) countLines(f *os.File) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		s := input.Text()
		// Keep in mind: http://code.google.com/p/go/issues/detail?id=3117
		t := c[s]
		if t.count == 0 {
			t.files = make(map[string]bool)
		}
		t.count++
		t.files[f.Name()] = true
		c[s] = t
	}
}
