package main

import (
	"fmt"
	"os"
)

func main() {
	for i, a := range os.Args[1:] {
		fmt.Println(i, a)
	}
}
