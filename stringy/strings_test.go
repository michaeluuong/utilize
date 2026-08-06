package stringy

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()

	os.Exit(code)

}

func setup() {

}

func teardown() {

}

// go test -benchtime=1s -bench . -cpuprofile cpu.prof
// go tool pprof cpu.prof
func BenchmarkFile(b *testing.B) {
	for b.Loop() {
	}

}

func TestSepSpace_SplitNoSpace_main(t *testing.T) {
	var sepString, expected, actual string

	//--------------------------------------------------------------------------------------------
	// Happy SepSpace()
	sepString = "   hello   |    goodbye     "
	expected = "hello|goodbye"
	actual = SepSpace(sepString, "|")
	assert.Equal(t, expected, actual)

	sepString = "   hello   ,    goodbye     ,      "
	expected = "hello,goodbye,"
	actual = SepSpace(sepString, ",")
	assert.Equal(t, expected, actual)

	//--------------------------------------------------------------------------------------------
	// Happy SplitNoSpace()
	sepString = "   hello   |    goodbye     "
	expectedSlice := []string{"hello", "goodbye"}
	actualSlice := SplitNoSpace(sepString, "|")
	assert.Equal(t, expectedSlice, actualSlice)

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

func TestReplaceUserString_main(t *testing.T) {
	var userString, expected, actual string

	var variables map[string]string = map[string]string{
		"Track Number": "1",
		"Title":        "Paint It Black",
		"Artist":       "The Bowling Crones",
	}

	//--------------------------------------------------------------------------------------------
	// Happy
	userString = "%02s{Track Number}. %{title} - %{artist}"
	expected = "01. Paint It Black - The Bowling Crones"
	actual = ReplaceUserString(userString, variables, TitleCase)
	assert.Equal(t, expected, actual)

	//--------------------------------------------------------------------------------------------
	// Unhappy
	userString = "%02s{Track Number}. %{ttle} - %{artist}"
	expected = "01. %{ttle} - The Bowling Crones"
	actual = ReplaceUserString(userString, variables, TitleCase)
	assert.Equal(t, expected, actual)

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

func TestItemSentence_main(t *testing.T) {
	var sepString, expected, actual string

	//--------------------------------------------------------------------------------------------
	// Happy
	sepString = "Led Zeppelin, Black Sabbath, Judas Priest"
	expected = "Led Zeppelin, Black Sabbath & Judas Priest"
	actual = ItemSentence(sepString, ",", "", "")
	assert.Equal(t, expected, actual)

	sepString = "Led Zeppelin, Black Sabbath, Judas Priest"
	expected = "->Led Zeppelin, Black Sabbath & Judas Priest<-"
	actual = ItemSentence(sepString, ",", "->", "<-")
	assert.Equal(t, expected, actual)

	sepString = "Led Zeppelin, Black Sabbath, Judas Priest"
	expected = "Led Zeppelin & Judas Priest"
	actual = ItemSentence(sepString, ",", "", "", map[string]bool{"Black Sabbath": true})
	assert.Equal(t, expected, actual)

	sepString = "Led Zeppelin, Black Sabbath, Judas Priest, & Deep Purple"
	expected = "Judas Priest & Deep Purple"
	actual = ItemSentence(sepString, ",", "", "", map[string]bool{"Led Zeppelin": true, "Black Sabbath": true, " Judas Priest": true})
	assert.Equal(t, expected, actual)

	sepString = "Led Zeppelin, Black Sabbath, Judas Priest, Deep Purple"
	expected = "Led Zeppelin, Black Sabbath & Judas Priest"
	actual = ItemSentence(sepString, ",", "", "", map[string]bool{"Deep Purple": true})
	assert.Equal(t, expected, actual)

	sepString = "Led Zeppelin, Black Sabbath, Judas Priest, the & Deep Purple"
	expected = "Led Zeppelin, Black Sabbath, Judas Priest, the & Deep Purple"
	actual = ItemSentence(sepString, ",", "", "", map[string]bool{"Deep Purple": true})
	assert.Equal(t, expected, actual)

	//--------------------------------------------------------------------------------------------
	// Cleanup
}

func TestPrintDataSlice_PrintData_main(t *testing.T) {
	var buf bytes.Buffer
	var expectedBuf *bytes.Buffer
	var expected string

	//--------------------------------------------------------------------------------------------
	// Slice
	data := [][]string{{"row1 col1", "row1 col2"}, {"row2 col1", "row2 col2"}}
	expected = "test\nrow1 col1 | row1 col2 \nrow2 col1 | row2 col2 \n\n"
	expectedBuf = bytes.NewBufferString(expected)
	PrintDataSlice("test", data, &buf)
	assert.Equal(t, *expectedBuf, buf)

	// No Header
	buf.Reset()
	expected = "row1 col1 | row1 col2 \nrow2 col1 | row2 col2 \n\n"
	expectedBuf = bytes.NewBufferString(expected)
	PrintDataSlice("", data, &buf)
	assert.Equal(t, *expectedBuf, buf)

	// Stderr
	PrintDataSlice("test stderr", data, os.Stderr)

	// Header only
	buf.Reset()
	expected = "test\n\n"
	expectedBuf = bytes.NewBufferString(expected)
	PrintDataSlice("test", [][]string{}, &buf)
	assert.Equal(t, *expectedBuf, buf)

	//--------------------------------------------------------------------------------------------
	// Map
	buf.Reset()
	dataMap := map[int][]string{0: {"row1 col1", "row2 col2"}, 1: {"row2 col1", "row2 col2"}}
	expected = "test\nrow1 col1 | row2 col2 \nrow2 col1 | row2 col2 \n\n"
	expectedBuf = bytes.NewBufferString(expected)
	PrintData("test", dataMap, &buf)
	assert.Equal(t, *expectedBuf, buf)

	PrintData("test stderr", dataMap, os.Stderr)

	// No header
	buf.Reset()
	expected = "row1 col1 | row2 col2 \nrow2 col1 | row2 col2 \n\n"
	expectedBuf = bytes.NewBufferString(expected)
	PrintData("", dataMap, &buf)
	assert.Equal(t, *expectedBuf, buf)

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

func TestPrintLine_main(t *testing.T) {
	//var expected, actual string

	//--------------------------------------------------------------------------------------------
	// Happy
	PrintLine("-", 10)
	PrintLine("-", 10, "stderr")
	PrintLine("-", 10, "//", "\\", "stderr")
	PrintLine("", 10, "//", "\\")

	//--------------------------------------------------------------------------------------------
	// Cleanup
}
