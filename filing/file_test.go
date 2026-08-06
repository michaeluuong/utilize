package filing

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
		ExecutablePath()
		Exists("")
		Cat("")
		Ls("")
	}

}

func TestCat_main(t *testing.T) {
	//--------------------------------------------------------------------------------------------
	// Happy
	originalFilename := "file.go.txt"
	catFilename := originalFilename + "_cat_output"
	f, err := os.Create(catFilename)
	assert.Nilf(t, err, "os.Create() should not have errored out: %s", catFilename)
	Cat(originalFilename, f)

	originalFilenameBites, err := os.ReadFile(originalFilename)
	catFilenameBites, err := os.ReadFile(catFilename)
	assert.Equalf(t, originalFilenameBites, catFilenameBites, "Contents of originalFilename(%s) should match catFilename(%s)", originalFilename, catFilename)

	//--------------------------------------------------------------------------------------------
	// Cleanup
	os.Remove(catFilename)

}

func TestCopyFile_main(t *testing.T) {
	var err error
	sourceFilename := "file.go.txt"
	emptySourceFilename := ""
	destinationFilename := "destination.txt"
	emptyDestinationFilename := ""

	//--------------------------------------------------------------------------------------------
	// Happy see TestReadFileToString_main

	//--------------------------------------------------------------------------------------------
	err = CopyFile(emptySourceFilename, destinationFilename)
	assert.EqualError(t, err, "open : no such file or directory")

	CopyFile(sourceFilename, emptyDestinationFilename)
	assert.EqualError(t, err, "open : no such file or directory")

	// Same as empySourceFile CopyFile("", "")

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

func TestDirEntryNameIter_DirEntryIter_main(t *testing.T) {
	testDir := "TestDirEntryNameIter"
	testFilename := "TestDirEntryNameIter"
	testFileExt := ".txt"
	os.Mkdir(testDir, 0755)

	var expectedFiles []string
	for i := range 5 {
		fullTestFilename := testFilename + "_" + strconv.Itoa(i+1) + testFileExt
		expectedFiles = append(expectedFiles, fullTestFilename)
		fullTestFile := filepath.Join(testDir, fullTestFilename)
		os.WriteFile(fullTestFile, []byte(fullTestFile), 0644)

	}

	//--------------------------------------------------------------------------------------------
	// Happy
	dirList, err := Ls(testDir)
	assert.Nil(t, err)
	var actualFiles []string
	for file := range DirEntryNameIter(dirList) {
		actualFiles = append(actualFiles, file)

	}
	assert.Equal(t, expectedFiles, actualFiles)

	//--------------------------------------------------------------------------------------------
	for file := range DirEntryNameIter(dirList) {
		if file == testFilename+"_3"+testFileExt {
			break

		}

	}

	//--------------------------------------------------------------------------------------------
	// DirEntryIter
	var actualInfoNames []string
	for fileInfo := range DirEntryIter(dirList) {
		actualInfoNames = append(actualInfoNames, fileInfo.Name())

	}
	assert.Equal(t, expectedFiles, actualInfoNames)

	for fileInfo := range DirEntryIter(dirList) {
		if fileInfo.Name() == testFilename+"_3"+testFileExt {
			break

		}

	}

	//--------------------------------------------------------------------------------------------
	// Cleanup
	os.RemoveAll(testDir)

}

func TestExecutablePath_main(t *testing.T) {
	var execPath string
	var err error

	//--------------------------------------------------------------------------------------------
	// Happy
	execPath, err = ExecutablePath()
	assert.Emptyf(t, execPath, "There should not be an executable path when running from test")

	//--------------------------------------------------------------------------------------------
	// os.Executable failed
	execError := errors.New("Can't find executable")
	executableSav := executable
	executable = func() (string, error) {
		return "", execError

	}
	execPath, err = ExecutablePath()

	executable = executableSav

	assert.ErrorIsf(t, err, execError, "Mocked executable should return \"%s\"", execError)

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

func TestExists_main(t *testing.T) {
	testFilename := "test_exists.txt"

	var exists bool

	//--------------------------------------------------------------------------------------------
	// Happy
	os.WriteFile(testFilename, []byte("test"), 0644)
	exists = Exists(testFilename)
	assert.Equalf(t, exists, true, "%s should exist", testFilename)

	os.Remove(testFilename)

	//--------------------------------------------------------------------------------------------
	// DNE
	exists = Exists(testFilename)
	assert.Equalf(t, exists, false, "%s should not exist", testFilename)

	//--------------------------------------------------------------------------------------------
	// Permission Denied
	workDir, _ := os.Getwd()
	tempDir, _ := os.MkdirTemp(workDir, "utilize")
	tempTestFilename := filepath.Join(tempDir, testFilename)
	os.WriteFile(tempTestFilename, []byte("test"), 0000)
	os.Chmod(tempDir, 0000)

	exists = Exists(tempTestFilename)
	os.Chmod(tempDir, DirPerm)
	os.Chmod(tempTestFilename, FilePerm)

	//--------------------------------------------------------------------------------------------
	// Cleanup
	os.RemoveAll(tempDir)

}

func TestLsEntryName_DirEntryName_main(t *testing.T) {
	testDir := "TestLsEntryName"
	os.Mkdir(testDir, 0755)

	testFilename := "TestLsEntryName.txt"
	testFilePath := filepath.Join(testDir, testFilename)
	os.WriteFile(testFilePath, []byte(filepath.Join(testDir+"-"+testFilePath)), 0644)

	var nameSlice []string

	//--------------------------------------------------------------------------------------------
	// Happy
	nameSlice = LsEntryName(testDir)
	assert.Equal(t, testFilename, nameSlice[0])

	//--------------------------------------------------------------------------------------------
	// Directory
	nameSlice = LsEntryName(".", testDir, "-d")
	assert.Equal(t, testDir, nameSlice[0])

	//--------------------------------------------------------------------------------------------
	// Empty
	nameSlice = LsEntryName("", testDir, "-d")
	assert.Equal(t, testDir, nameSlice[0])

	//--------------------------------------------------------------------------------------------
	// addDirOk
	nameSlice = LsEntryName(testDir, "add_dir")
	assert.Equal(t, testFilePath, nameSlice[0])

	//--------------------------------------------------------------------------------------------
	// Cleanup
	os.Remove(testFilePath)
	os.Remove(testDir)

}

func TestLsEntryNameSet_DirEntryNameSet_main(t *testing.T) {
	testDir := "TestLsEntryNameSet"
	os.Mkdir(testDir, 0755)

	testFilename := "TestLsEntryNameSet.txt"
	testFilePath := filepath.Join(testDir, testFilename)
	os.WriteFile(testFilePath, []byte(filepath.Join(testDir+"-"+testFilePath)), 0644)

	//--------------------------------------------------------------------------------------------
	// Happy
	nameSet := LsEntryNameSet(testDir)
	assert.True(t, nameSet[testFilename])

	//--------------------------------------------------------------------------------------------
	// Cleanup
	os.Remove(testFilePath)
	os.Remove(testDir)

}

func TestLsGlob_main(t *testing.T) {
	filename := "file.go.txt"

	//--------------------------------------------------------------------------------------------
	files, err := LsGlob("f[iI]*.g[oO].txt")
	assert.Nil(t, err)
	assert.Equal(t, filename, files[0])
	fmt.Printf("files: %v\n", files)

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

func TestLs_PrintDirEntry_longestDirEntryName_main(t *testing.T) {
	testDir := "TEST"
	os.Mkdir(testDir, 0755)
	var dirList []os.DirEntry
	var regex string

	//--------------------------------------------------------------------------------------------
	// Here
	//dirList, _ := Ls(".", "(file|file_test).go$|"+testDir, "-d")
	dirList, _ = Ls(".", "^TEST$", "-d")
	info, _ := dirList[0].Info()
	assert.Equal(t, testDir, info.Name())
	assert.True(t, info.IsDir())

	// PrintDirEntry() Stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	PrintDirEntry(dirList)
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	regex = fmt.Sprintf("^drwxr-xr-x   %-20s[0-9]+$", testDir)
	assert.Regexp(t, regex, strings.TrimSpace(buf.String()))

	// PrintDirEntry() file
	testOutFilename := "ls_print_dir_entry.txt"
	testOutFile := filepath.Join(testDir, testOutFilename)
	testOutFileWriter, err := os.Create(testOutFile)
	assert.Nil(t, err)
	defer testOutFileWriter.Close()
	dirList, _ = Ls(testDir)
	PrintDirEntry(dirList, testOutFileWriter)

	actual, err := ReadFileToString(testOutFile)
	assert.Nil(t, err)
	testOutFilenameVerb := "%-" + strconv.Itoa(len(testOutFilename)+5) + "s"
	regex = fmt.Sprintf("^-rw-r--r--   "+testOutFilenameVerb+"[0-9]+$", testOutFilename)
	assert.Regexp(t, regex, strings.TrimSpace(actual))

	//--------------------------------------------------------------------------------------------
	// Empty dir
	dirList, _ = Ls("", "-d")
	info, _ = dirList[0].Info()
	assert.Equal(t, testDir, info.Name())
	assert.True(t, info.IsDir())

	//--------------------------------------------------------------------------------------------
	// err != nil
	os.Chmod(testDir, 0000)
	dirList2, err := Ls(testDir, "-d")
	assert.Empty(t, dirList2)
	assert.EqualError(t, err, "could not read directory: "+testDir)

	os.Chmod(testDir, 0755)

	//--------------------------------------------------------------------------------------------
	// Cleanup
	os.Remove(testOutFile)
	os.Remove(testDir)

}

func TestNextFile_main(t *testing.T) {
	testDir := "TEST_NEXT_FILE"
	testFileNoExt := "test_next_file"
	testFileExt := ".txt"
	testFilename := "test_next_file" + testFileExt
	var newTestFilename, filePath, expectedFilename string

	//--------------------------------------------------------------------------------------------
	newTestFilename = NextFile(testFilename)
	assert.Equal(t, testFilename, newTestFilename)
	os.WriteFile(newTestFilename, []byte(newTestFilename), FilePerm)
	os.Remove(newTestFilename)

	//--------------------------------------------------------------------------------------------
	// Directory
	os.Mkdir(testDir, DirPerm)

	filePath = filepath.Join(testDir, testFilename)
	newTestFilename = NextFile(filePath)
	expectedFilename = filepath.Join(testDir, testFilename)
	assert.Equal(t, expectedFilename, newTestFilename)
	os.WriteFile(newTestFilename, []byte(newTestFilename), FilePerm)

	newTestFilename = NextFile(newTestFilename)
	expectedFilename = filepath.Join(testDir, testFileNoExt+"_2"+testFileExt)
	assert.Equal(t, expectedFilename, newTestFilename)
	os.WriteFile(newTestFilename, []byte(newTestFilename), FilePerm)

	//--------------------------------------------------------------------------------------------
	// Empty
	newTestFilename = NextFile("")
	assert.Equal(t, "", newTestFilename)

	//--------------------------------------------------------------------------------------------
	// Cleanup
	os.RemoveAll(testDir)

}

func TestNormalizeFilename_main(t *testing.T) {
	testFilename := "testing_file.txt"
	workDir, _ := os.Getwd()
	testFilePath := filepath.Join(workDir, testFilename)

	var normalizedFilename string

	//--------------------------------------------------------------------------------------------
	// Happy
	normalizedFilename = NormalizeFilename(testFilename)
	assert.Equal(t, testFilename, normalizedFilename)

	//--------------------------------------------------------------------------------------------
	// Path
	normalizedFilename = NormalizeFilename(testFilePath)
	assert.NotEqual(t, testFilePath, normalizedFilename)

	//--------------------------------------------------------------------------------------------
	// Colon
	normalizedFilename = NormalizeFilename(testFilename + ":")
	assert.Equal(t, testFilename+" ", normalizedFilename)

	//--------------------------------------------------------------------------------------------
	// Curly
	normalizedFilename = NormalizeFilename(testFilename + "’")
	assert.Equal(t, testFilename+"'", normalizedFilename)

	//--------------------------------------------------------------------------------------------
	// Allowed
	allowedFile := "azAZ09().'+-_!@&"
	normalizedFilename = NormalizeFilename(allowedFile)
	assert.Equal(t, allowedFile, normalizedFilename)

	//--------------------------------------------------------------------------------------------
	// Underscores
	spaceFile := " test file .txt"
	normalizedFilename = NormalizeFilename(spaceFile, true)
	assert.Equal(t, strings.ReplaceAll(spaceFile, " ", "_"), normalizedFilename)

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

func TestReadFileToString_main(t *testing.T) {
	sourceFile := "file.go.txt"
	readTestFile := "read_test.txt"
	CopyFile(sourceFile, readTestFile)

	//--------------------------------------------------------------------------------------------
	// Happy
	expected, err := os.ReadFile(readTestFile)
	assert.Nil(t, err)

	readTestFileContents, err := ReadFileToString(readTestFile)
	assert.Equal(t, string(expected), readTestFileContents)

	os.Remove(readTestFile)

	// Empty
	readEmptyContents, err := ReadFileToString("")
	assert.Empty(t, readEmptyContents)
	assert.EqualError(t, err, "open : no such file or directory")

	//--------------------------------------------------------------------------------------------
	// Cleanup

}
