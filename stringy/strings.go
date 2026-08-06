// * common string operations
package stringy

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/michaeluuong/utilize/filing"
	"golang.org/x/text/cases"
)

// SepSpace removes leading and trailing space characters, and spaces around the separation character.
//   - sepString is the string to remove spaces from
//   - sep is the separation character to remove spaces around (cannot use |)
func SepSpace(sepString, sep string) string {
	regSep := sep
	if sep == "|" {
		regSep = "\\|"

	}

	regex := "\\s*" + regSep + "\\s*"
	spaceRegexp := regexp.MustCompile(regex)

	newSepString := spaceRegexp.ReplaceAllString(sepString, sep)

	return strings.TrimSpace(newSepString)

}

// SplitNoSpace splits a string into substrings after removing leading and trailing spaces around a separation character,
// including leading and trailing spaces.
//   - sepString is the string to split
//   - sep is the separation character to remove spaces around
func SplitNoSpace(sepString, sep string) []string {
	return strings.Split(SepSpace(sepString, sep), sep)

}

// ReplaceUserString replaces strings within user variable declarations "%{string}" with values from the variables map.
//   - template is the string containing user variables to replace
//   - variables is keyed by user variables and valued by the replacement string
//   - keyCaseOpt use "title" for keys in title case, "lower" for all lowercase, or "upper" for all uppercase
//
// Return template with variables replaced by values from the variable map.
func ReplaceUserString(template string, variables map[string]string, keyCaseOpt ...StringCase) string {
	re := regexp.MustCompile(`%(0[0-9]+[dfs])?\{[^}]*\}`)

	var titler cases.Caser
	if len(keyCaseOpt) > 0 {
		titler = keyCaseOpt[0].Caser()

	}

	result := re.ReplaceAllStringFunc(template, func(match string) string {
		index := strings.Index(match, "{")
		verb := match[0:index]
		varName := match[index+1 : len(match)-1]
		if titler != (cases.Caser{}) {
			varName = titler.String(varName)

		}

		//fmt.Printf("verb: %s, varName: %s, index: %d\n", verb, varName, index)
		if val, exists := variables[varName]; exists {
			if verb != "%" {
				val = fmt.Sprintf(verb, val)

			}

			return val

		}

		return match

	})

	return result

}

// ItemSentence separates a string into substrings and joins them into an itemized sentence (e.g. 1, 2, 3 = 1, 2 & 3).
//   - sepString the string containing the separator character to split
//   - sep the separator character
//   - prefix is a word or phrase to place at the beginning of the resulting string
//   - suffix is a word or phrase to place at the end of the resulting string
//   - excludeOpt is a set of words to exclude from the final sentence
//
// Return a string of items separated by commas with the last item separated by &.
func ItemSentence(sepString, sep, prefix, suffix string, excludeOpt ...map[string]bool) string {
	sepStringParts := SplitNoSpace(sepString, sep)
	sepStringPartsLen := len(sepStringParts)

	var exclude map[string]bool
	if len(excludeOpt) > 0 {
		exclude = excludeOpt[0]

	}

	var itemBuilder strings.Builder
	for i, sepStringPart := range sepStringParts {
		if _, ok := exclude[sepStringPart]; !ok {
			if itemBuilder.Len() > 0 {
				var excludeNext bool = false
				if i == sepStringPartsLen-2 {
					if _, ok = exclude[sepStringParts[i+1]]; ok {
						excludeNext = true

					}

				}

				if i == sepStringPartsLen-1 || excludeNext {
					if !strings.Contains(sepStringPart, "&") {
						itemBuilder.WriteString(" & ")

					} else if strings.HasPrefix(sepStringPart, "&") {
						itemBuilder.WriteString(" ")

					} else {
						itemBuilder.WriteString(", ")

					}

				} else {
					itemBuilder.WriteString(", ")

				}

			}

			itemBuilder.WriteString(sepStringPart)

		}

	}

	var itemString string
	if itemBuilder.Len() > 0 {
		itemString = prefix + itemBuilder.String() + suffix

	}

	return itemString

}

func PrintDataSlice(header string, data [][]string, outputOpt ...io.Writer) {
	var out io.Writer
	out = os.Stdout
	if len(outputOpt) > 0 {
		out = outputOpt[0]

	}

	maxColumnLengths := maxColumnLength2(data)

	if header != "" {
		fmt.Fprintf(out, "%s\n", header)

	}
	var rowBuilder strings.Builder
	for rowI := 0; rowI < len(data); rowI++ {
		const columnSpace = 1
		for colI := 0; colI < len(data[rowI]); colI++ {
			verb := "%-" + strconv.Itoa(maxColumnLengths[colI]+columnSpace) + "s"
			if colI < len(data[rowI])-1 {
				verb += "| "

			}
			fmt.Fprintf(&rowBuilder, verb, data[rowI][colI])

		}
		fmt.Fprintf(out, "%s\n", rowBuilder.String())

		rowBuilder.Reset()

	}
	fmt.Fprintf(out, "\n")

}

// func PrintData(header string, data map[int][]string, outputOpt ...*os.File) {
func PrintData(header string, data map[int][]string, outputOpt ...io.Writer) {
	var out io.Writer
	out = os.Stdout
	if len(outputOpt) > 0 {
		out = outputOpt[0]

	}

	maxColumnLengths := maxColumnLength(data)

	if header != "" {
		fmt.Fprintf(out, "%s\n", header)

	}
	var rowBuilder strings.Builder
	for rowI := 0; rowI < len(data); rowI++ {
		const columnSpace = 1
		for colI := 0; colI < len(data[rowI]); colI++ {
			verb := "%-" + strconv.Itoa(maxColumnLengths[colI]+columnSpace) + "s"
			if colI < len(data[rowI])-1 {
				verb += "| "

			}
			fmt.Fprintf(&rowBuilder, verb, data[rowI][colI])

		}
		fmt.Fprintf(out, "%s\n", rowBuilder.String())

		rowBuilder.Reset()

	}
	fmt.Fprintf(out, "\n")

}

func maxColumnLength2(data [][]string) []int {
	if len(data) == 0 {
		return []int{}

	}

	maxLengths := make([]int, len(data[0]))

	for rowI := 0; rowI < len(data); rowI++ {
		for colI := 0; colI < len(data[rowI]); colI++ {
			maxLengths[colI] = max(len(data[rowI][colI]), maxLengths[colI])

		}

	}

	return maxLengths

}

func maxColumnLength(data map[int][]string) []int {
	var maxLengths []int

	for rowI := 0; rowI < len(data); rowI++ {
		if maxLengths == nil {
			maxLengths = make([]int, len(data[rowI]))

		}

		for colI := 0; colI < len(data[rowI]); colI++ {
			maxLengths[colI] = max(len(data[rowI][colI]), maxLengths[colI])

		}

	}

	return maxLengths

}

// PrintLine prints a line to stdout.
//   - lineChar is the character(s) used to print the line
//   - length is the number of time to repeat lineChar
//   - prePostOpts index 0 should be the string to prefix the line with while index 1 should be the string to suffix the line with or use "stderr" to print to stderr
func PrintLine(lineChar string, length int, prePostOpts ...string) {
	var w io.Writer = os.Stdout
	if len(prePostOpts) > 0 {
		var stderrFound bool
		prePostOpts, stderrFound = filing.FindRegexOpt("stderr", prePostOpts...)
		if stderrFound {
			w = os.Stderr

		}

	}

	if lineChar == "" {
		lineChar = "-"

	}

	var pre, post string
	prePostOptsLen := len(prePostOpts)
	if prePostOptsLen > 0 {
		pre = prePostOpts[0]
		if prePostOptsLen > 1 {
			post = prePostOpts[0]

		}

	}

	fmt.Fprintf(w, "%s%s%s\n", pre, strings.Repeat(lineChar, length), post)

}
