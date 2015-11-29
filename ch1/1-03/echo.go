package ch1

import (
	"os"
	"strings"
)

func ConcatArgs(args []string) string {
	s, sep := "", ""
	for _, arg := range args {
		s += sep + arg
		sep = " "
	}
	return s
}

func JoinArgs(args []string) string {
	return strings.Join(os.Args[1:], " ")
}
