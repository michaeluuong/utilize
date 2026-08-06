package filing

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkDir(b *testing.B) {
	for b.Loop() {
	}

}

func TestIsDir_main(t *testing.T) {
	var isDir bool
	var err error
	testDir := "TestIsDir"
	os.Mkdir(testDir, 0755)

	testFile := filepath.Join(testDir, "test_is_dir.txt")
	os.WriteFile(testFile, []byte("test is dir\n"), FilePerm)

	//--------------------------------------------------------------------------------------------
	// Happy
	isDir, err = IsDir(testDir)
	assert.Nil(t, err)
	assert.True(t, isDir)

	//--------------------------------------------------------------------------------------------
	// File
	isDir, err = IsDir(testFile)
	assert.Nil(t, err)
	assert.False(t, isDir)

	//--------------------------------------------------------------------------------------------
	// Empty & error
	isDir, err = IsDir("")
	assert.EqualError(t, err, "stat : no such file or directory")
	assert.False(t, isDir)

	//--------------------------------------------------------------------------------------------
	// Here
	isDir, err = IsDir(".")
	assert.Nil(t, err)
	assert.True(t, isDir)

	//--------------------------------------------------------------------------------------------
	// There
	isDir, err = IsDir("..")
	assert.Nil(t, err)
	assert.True(t, isDir)

	//--------------------------------------------------------------------------------------------
	// Error
	testDirEnDir := filepath.Join(testDir, testDir)
	os.Mkdir(testDirEnDir, 0000)
	err = os.Chmod(testDir, 0000)
	isDir, err = IsDir(testDirEnDir)

	os.Chmod(testDir, 0755)
	os.Chmod(testDirEnDir, 0755)
	os.Remove(testDirEnDir)

	//--------------------------------------------------------------------------------------------
	// Cleanup
	os.Remove(testFile)
	os.Remove(testDir)

}

func TestLast_lastPart_RemoveLast_main(t *testing.T) {
	parts := []string{"once", "i", "was", "last"}
	splitString := strings.Join(parts, "_")
	var testLastPart, testLastStringPart, testRemoveLast string
	testNoSep := strings.Join(parts, "")

	//--------------------------------------------------------------------------------------------
	// last
	testLastPart = last(parts)
	assert.Equal(t, parts[len(parts)-1], testLastPart)

	// Empty parts
	testLastPart = last([]string{})
	assert.Empty(t, testLastPart)

	//--------------------------------------------------------------------------------------------
	// lastStringPart
	testLastStringPart = lastStringPart(splitString, "_")
	assert.Equal(t, parts[len(parts)-1], testLastStringPart)

	// No separator in string
	testLastStringPart = lastStringPart(testNoSep, "_")
	assert.Equal(t, testNoSep, testLastStringPart)

	// No separator
	testLastStringPart = lastStringPart(splitString, "")
	assert.Equal(t, splitString, testLastStringPart)

	// Empty string
	testLastStringPart = lastStringPart("", "_")
	assert.Empty(t, testLastStringPart)

	//--------------------------------------------------------------------------------------------
	// RemoveLast
	testRemoveLast = RemoveLast(splitString, "_")
	expected := strings.Join(parts[:len(parts)-1], "_")
	assert.Equal(t, expected, testRemoveLast)

	// No separator in string
	testRemoveLast = RemoveLast(testNoSep, "_")
	assert.Equal(t, testNoSep, testRemoveLast)

	// No separator
	testRemoveLast = RemoveLast(testNoSep, "")
	assert.Equal(t, testNoSep, testRemoveLast)

	// Empty string
	testRemoveLast = RemoveLast("", "")
	assert.Empty(t, testRemoveLast)

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

func TestNextDir_main(t *testing.T) {
	//workDir, _ := os.Getwd()

	testDir := "TestNextDir"
	testDirEnDir := filepath.Join(testDir, testDir)
	os.Mkdir(testDir, 0755)

	var newDir string
	var err error

	//--------------------------------------------------------------------------------------------
	// Happy
	// TestNextDir/TestNextDir
	newDir, err = NextDir(testDirEnDir)
	assert.Nil(t, err)

	os.Mkdir(newDir, 0755)
	assert.Equal(t, filepath.Join(testDir, testDir), newDir)

	// TestNextDir/TestNextDir_2
	newDir, err = NextDir(newDir)
	assert.Nil(t, err)
	os.Mkdir(newDir, 0755)
	assert.Equal(t, filepath.Join(testDir, testDir+"_2"), newDir)

	// TestNextDir/TestNextDir_3
	newDir, err = NextDir(newDir)
	assert.Nil(t, err)
	os.Mkdir(newDir, 0755)
	assert.Equal(t, filepath.Join(testDir, testDir+"_3"), newDir)

	//--------------------------------------------------------------------------------------------
	// Error
	// TestNextDir/TestNextDir_4
	// It seemed like permissions were caching so need a different directory than before
	errDir := "TestErrDir"
	newErrDir := filepath.Join(errDir, errDir+"_1")
	fmt.Printf("errDir: %v, newErrDir: %v\n", errDir, newErrDir)
	os.Mkdir(errDir, 0000)

	newErrDir, err = NextDir(newErrDir)
	fmt.Printf("errDir: %v, newErrDir: %v, err: %v\n", errDir, newErrDir, err)
	os.Mkdir(newErrDir, 0755)
	assert.ErrorContains(t, err, "could not read directory: ")

	os.RemoveAll(errDir)

	//--------------------------------------------------------------------------------------------
	// Find an opening
	// TestNextDir/TestNextDir_3
	testOpenDir3 := filepath.Join(testDir, testDir+"_3")
	testOpenDir4 := filepath.Join(testDir, testDir+"_4")
	os.Mkdir(testOpenDir4, 0755)
	os.Remove(testOpenDir3)
	newDir, err = NextDir(testOpenDir4, "_")
	assert.Equal(t, testOpenDir3, newDir)

	//--------------------------------------------------------------------------------------------
	/*cmd := exec.Command("ls", "-lat")
	out, err := cmd.Output()
	fmt.Printf("out: %s, err: %v\n", out, err)*/

	//--------------------------------------------------------------------------------------------
	// Cleanup
	os.Chmod(testDir, 0777)
	os.RemoveAll(testDir)

}
