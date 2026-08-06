package stringy

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestCat_main(t *testing.T) {
	var sentence, expected, actual string

	//--------------------------------------------------------------------------------------------
	// Happy
	// SentenceCase, first letter of first sentence is capitalized
	sentence = "hello world where the existance of self is universal. it is the Universe that is singular."
	expected = "Hello world where the existance of self is universal. it is the Universe that is singular."
	actual = CaseString(sentence, SentenceCase)
	assert.Equal(t, expected, actual)

	// DefaultCase, nothing happens
	expected = ""
	actual = CaseString(sentence, DefaultCase)
	assert.Equal(t, expected, actual)

	// FoldCase, fold to lower case
	expected = "hello world where the existance of self is universal. it is the universe that is singular."
	actual = CaseString(sentence, FoldCase)
	assert.Equal(t, expected, actual)

	// LowerCase, convert complete sentence to lowercase
	sentence2 := "HELLO world where the existance of self is universal. it is the UNIVERSE that is singular."
	expected = "hello world where the existance of self is universal. it is the universe that is singular."
	actual = CaseString(sentence2, LowerCase)
	assert.Equal(t, expected, actual)

	// TitleCase, transform first letter of each word to upper case
	expected = "Hello World Where The Existance Of Self Is Universal. It Is The Universe That Is Singular."
	actual = CaseString(sentence, TitleCase)
	assert.Equal(t, expected, actual)

	expected = "HELLO WORLD WHERE THE EXISTANCE OF SELF IS UNIVERSAL. IT IS THE UNIVERSE THAT IS SINGULAR."
	actual = CaseString(sentence, UpperCase)
	assert.Equal(t, expected, actual)

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

func TestString_ToCase_main(t *testing.T) {
	var expected, actual string

	//--------------------------------------------------------------------------------------------
	// Happy
	// TitleCase
	expected = "TitleCase"
	actual = fmt.Sprintf("%s", TitleCase)
	assert.Equal(t, expected, actual)

	actual = fmt.Sprintf("%s", DefaultCase.ToCase("title"))
	assert.Equal(t, expected, actual)

	// DefaultCase
	expected = "DefaultCase"
	actual = fmt.Sprintf("%s", DefaultCase)
	assert.Equal(t, expected, actual)

	actual = fmt.Sprintf("%s", DefaultCase.ToCase("hello"))
	assert.Equal(t, expected, actual)

	// FoldCase
	expected = "FoldCase"
	actual = fmt.Sprintf("%s", FoldCase)
	assert.Equal(t, expected, actual)

	actual = fmt.Sprintf("%s", DefaultCase.ToCase("fold"))
	assert.Equal(t, expected, actual)

	// LowerCase
	expected = "LowerCase"
	actual = fmt.Sprintf("%s", LowerCase)
	assert.Equal(t, expected, actual)

	actual = fmt.Sprintf("%s", DefaultCase.ToCase("lower"))
	assert.Equal(t, expected, actual)

	// SentenceCase
	expected = "SentenceCase"
	actual = fmt.Sprintf("%s", SentenceCase)
	assert.Equal(t, expected, actual)

	actual = fmt.Sprintf("%s", DefaultCase.ToCase("sentence"))
	assert.Equal(t, expected, actual)

	// UpperCase
	expected = "UpperCase"
	actual = fmt.Sprintf("%s", UpperCase)
	assert.Equal(t, expected, actual)

	actual = fmt.Sprintf("%s", DefaultCase.ToCase("upper"))
	assert.Equal(t, expected, actual)

	//--------------------------------------------------------------------------------------------
	// Cleanup

}
func TestCaser_main(t *testing.T) {
	//--------------------------------------------------------------------------------------------
	// Happy
	titler := TitleCase.Caser(language.Afrikaans)
	fmt.Printf("%s\n", titler.String("hello world! goodbye cruelty?"))

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

func TestSentenceCasing_main(t *testing.T) {
	var sentence, expected, actual string

	//--------------------------------------------------------------------------------------------
	// Happy
	sentence = "happy am I. are you not? why aren't you! dude."
	expected = "Happy am I. Are you not? Why aren't you! Dude."
	actual = SentenceCasing(sentence)
	assert.Equal(t, expected, actual)

	sentence = ""
	expected = ""
	actual = SentenceCasing(sentence)
	assert.Equal(t, expected, actual)

	//--------------------------------------------------------------------------------------------
	// Cleanup

}
