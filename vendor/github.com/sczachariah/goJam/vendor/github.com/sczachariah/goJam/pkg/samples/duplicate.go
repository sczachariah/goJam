package samples

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func StdinDup() {
	counts := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)

	for input.Scan() {
		if input.Text() == "" {
			break
		}
		counts[input.Text()]++
	}

	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

func FileDup() {
	counts := make(map[string]int)
	filenames := []string{"D:\\FMW\\IDE\\IntelliJ\\goJam\\main.go"}

	for _, filename := range filenames[0:] {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "FileDup: %v\n", err)
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			counts[line]++
		}
	}

	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}
