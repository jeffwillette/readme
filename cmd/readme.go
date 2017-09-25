package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

var (
	// this is the regex for finding headings in the readme
	r = regexp.MustCompile(`^#{1,6}`)

	// this regex is for turning on/off heading highlighting, since comments in code
	// blocks can look a lot like mardown headers.
	c = regexp.MustCompile("```")

	header = color.New(color.FgMagenta, color.Bold, color.Underline, color.Italic)
	code   = color.New(color.FgBlue, color.Faint)
)

// Readme gets and prints the file to the terminal
func Readme(s []string) {
	split := strings.Split(s[0], "/")

	domain := "http://raw.githubusercontent.com"
	readme := "README.md"
	branch := "master"

	everything := []string{
		domain, split[0], split[1], branch, readme,
	}

	url := strings.Join(everything, "/")

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	isCodeBlock := false

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		t := scanner.Text()

		if c.MatchString(t) {
			isCodeBlock = !isCodeBlock
		}

		switch {
		case r.MatchString(t) && !isCodeBlock:
			header.Print(t)
			fmt.Print("\n\n")
		case isCodeBlock:
			code.Println(t)
		default:
			fmt.Println(t)
		}
	}

	fmt.Print("\n\n")
}
