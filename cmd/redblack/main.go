package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nathanjcochran/redblack"
)

func main() {
	t := &redblack.Tree{}
	fmt.Println(t)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		if len(text) == 0 {
			continue
		}

		var del bool
		if text[0] == 'd' {
			del = true
			text = strings.TrimSpace(text[1:])
		}

		val, err := strconv.Atoi(text)
		if err != nil {
			fmt.Println("Error: invalid integer")
			continue
		}

		switch del {
		case false:
			t.Insert(val)
		case true:
			t.Remove(val)
		}

		fmt.Println(t)
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading from stdin: %s", err)
	}
}
