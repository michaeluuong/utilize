// * string case conversion
package stringy

import (
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/transform"
)

type StringCase int

const (
	DefaultCase  StringCase = iota
	FoldCase                // fold to lower case
	LowerCase               // transform all letters to lowercase
	SentenceCase            // capitalize first letter of sentence
	TitleCase               // capitalize first letter of each word
	UpperCase               // transform all letters to uppercase
)

func (s StringCase) String() string {
	caseString := ""
	switch s {
	case DefaultCase:
		caseString = "DefaultCase"

	case FoldCase:
		caseString = "FoldCase"

	case LowerCase:
		caseString = "LowerCase"

	case SentenceCase:
		caseString = "SentenceCase"

	case TitleCase:
		caseString = "TitleCase"

	case UpperCase:
		caseString = "UpperCase"

	}

	return caseString

}

func (s StringCase) ToCase(userCase string) StringCase {
	var stringCase StringCase
	switch strings.ToLower(userCase) {
	case "fold":
		stringCase = FoldCase

	case "lower":
		stringCase = LowerCase

	case "sentence":
		stringCase = SentenceCase

	case "title":
		stringCase = TitleCase

	case "upper":
		stringCase = UpperCase

	default:
		stringCase = DefaultCase

	}

	return stringCase

}

func (s StringCase) Caser(langOpt ...language.Tag) cases.Caser {
	var lang language.Tag = language.English
	if len(langOpt) > 0 {
		lang = langOpt[0]

	}

	var caser cases.Caser

	switch s {
	case FoldCase:
		caser = cases.Fold()

	case LowerCase:
		caser = cases.Lower(lang)

	case TitleCase:
		caser = cases.Title(lang)

	case UpperCase:
		caser = cases.Upper(lang)

	}

	return caser

}

// SentenceCase capitalizes the first letter of a sentence.
//   - sentence is the string to capitalize the first letter of
//
// Return sentence with the first letter capitalized
func SentenceCasing(sentence string) string {
	if sentence == "" {
		return ""

	}

	runes := []rune(sentence)
	runes[0] = unicode.ToUpper(runes[0])
	beg := false
	for i, r := range runes[1:] {
		if beg == true && r != ' ' {
			runes[i+1] = unicode.ToUpper(runes[i+1])
			beg = false

		}

		if r == '.' || r == '!' || r == '?' {
			beg = true

		}

	}
	return string(runes)

}

// CaseString converts a string that has been transformed to stringCase.
//   - input is the string to transform
//   - stringCase is the case to transform to
//
// Return the transformed string
func CaseString(input string, stringCase StringCase) string {
	var casedString string
	if stringCase == SentenceCase {
		casedString, _, _ = transform.String(&SentenceTransformer{}, input)

	} else {
		var caser cases.Caser = stringCase.Caser()
		if caser != (cases.Caser{}) {
			casedString = caser.String(input)

		}

	}

	return casedString

}
