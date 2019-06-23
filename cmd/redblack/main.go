package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/nathanjcochran/redblack"
)

func main() {
	t := &redblack.Tree{}
	for _, arg := range parseArgs() {
		t.Insert(arg)
		fmt.Println(t)
	}
}

func usage() {
	fmt.Printf("Usage: %s <int>...\n", os.Args[0])
	os.Exit(2)
}

func parseArgs() []int {
	if len(os.Args) == 1 {
		usage()
	}

	var ints []int
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-h", "--help":
			usage()
		}

		i, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Printf("Invalid int: %s\n", arg)
			usage()
		}
		ints = append(ints, i)
	}
	return ints
}
