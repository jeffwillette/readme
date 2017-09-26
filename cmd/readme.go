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

	// regex for inline code
	ic = regexp.MustCompile("`(.*?)`")

	header = color.New(color.FgMagenta, color.Bold, color.Underline, color.Italic)
	code   = color.New(color.FgBlue, color.Faint)
)

// EnableColors causes the colors to be enabled by default. The color library usually
// tests to see if the stdout is a terminal first and disables them if not. With colors
// always enabled, the user can pipe to `less -R` and still get color output
func EnableColors() {
	header.EnableColor()
	code.EnableColor()
}

// startCodeBlock returns true if it was not previously a code block and the regex matches
// the beginning of a code block.
func startCodeBlock(was, match bool) bool {
	return !was && match
}

// endCodeBlock returns true if it was previously a code block and the regex matches the
// code block closign tag.
func endCodeBlock(was, match bool) bool {
	return was && match
}

// PrintInlineCode prints when there is code in between backticks (`code`)
func PrintInlineCode(s string) {
	print := code.SprintFunc()
	inlineCodedString := ic.ReplaceAllStringFunc(s, func(foo string) string {
		return print(foo)
	})

	fmt.Println(inlineCodedString)
}

// Readme gets and prints the file to the terminal
func Readme(s []string) {

	EnableColors()

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

		wasCodeBlock := isCodeBlock
		if startCodeBlock(wasCodeBlock, c.MatchString(t)) {
			isCodeBlock = true
		}

		switch {
		case isCodeBlock:
			code.Println(t)
		case r.MatchString(t):
			header.Print(t)
			fmt.Print("\n\n")
		case ic.MatchString(t):
			PrintInlineCode(t)
		default:
			fmt.Println(t)
		}

		if endCodeBlock(wasCodeBlock, c.MatchString(t)) {
			isCodeBlock = false
		}
	}

	fmt.Print("\n\n")
}
