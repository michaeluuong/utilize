package jsoning

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TestEase struct {
	Artist       string `json:"artist,omitempty"`
	Album        string `json:"album,omitempty"`
	Album_artist string `json:"album_artist,omitempty"`
	Genre        string `json:"genre,omitempty"`
	Year         int    `json:"year,omitzero"`
}

var testEase TestEase
var testJSONFile = "music.json"
var testJSONWriteFile = "music_write.json"

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()

	os.Exit(code)

}

func setup() {
	testEase = TestEase{
		Artist:       "Deftones",
		Album:        "Around the Fur",
		Year:         1997,
		Album_artist: "Around the Fur",
	}

	err := NewJsonify().WriteObjToJSONFile(testJSONFile, testEase)
	if err != nil {
		panic(err)

	}

}

func teardown() {
	os.Remove(testJSONFile)
	os.Remove(testJSONWriteFile)

}

// go test -benchtime=1s -bench . -cpuprofile cpu.prof -memprofile mem.prof
// go tool pprof cpu.prof
func BenchmarkUnmarshall(b *testing.B) {
	var testEaseTemp TestEase
	newJsonify := NewJsonify()
	for b.Loop() {
		test2AnyUnmarshal2Struct(testJSONFile, &testEaseTemp)

		newJsonify.WriteObjToJSONFile(testJSONWriteFile, testEase)

	}

}

func test2AnyUnmarshal2Struct(jsonFilename string, testEase *TestEase) error {
	jsonify := NewJsonify()
	var testEaseTemp any = testEase
	return jsonify.Unmarshal2Struct(jsonFilename, &testEaseTemp)

}

// go test -v -cover ./...
// go test -v -cover ./jsonify_test.go ./jsonify.go
func TestUnmarshal2Struct_main(t *testing.T) {
	var testEaseActual TestEase
	var err error

	//--------------------------------------------------------------------------------------------
	// Happy
	err = test2AnyUnmarshal2Struct(testJSONFile, &testEaseActual)
	fmt.Printf("testEase: %#v\n", testEase)
	fmt.Printf("testEaseActual: %#v\n", testEaseActual)

	assert.Equal(t, testEase, testEaseActual)
	assert.Nilf(t, err, "testEase and testEaseActual should be equal")

	//--------------------------------------------------------------------------------------------
	// Empty file
	testEaseActual = TestEase{}
	err = test2AnyUnmarshal2Struct("nofile.json", &testEaseActual)
	assert.Emptyf(t, testEaseActual, "Empty file, object should be empty: %v", testEaseActual)
	assert.NotNil(t, err)

	// No file
	err = test2AnyUnmarshal2Struct("", &testEaseActual)
	assert.Emptyf(t, testEaseActual, "Empty file, object should be empty: %v", testEaseActual)
	assert.NotNil(t, err)

	//--------------------------------------------------------------------------------------------
	// Improper JSON
	badJsonFile := "bad.json"
	err = os.WriteFile(badJsonFile, []byte("json"), 0644)
	if err != nil {
		panic(err)
	}

	err = test2AnyUnmarshal2Struct(badJsonFile, &testEaseActual)
	assert.ErrorContainsf(t, err, "Error unmarshaling JSON: invalid character", "Should not be able to unmarshal file %s", badJsonFile)

	os.Remove(badJsonFile)

	//--------------------------------------------------------------------------------------------
	// Empty field (genre)
	testEaseActual = TestEase{}
	err = test2AnyUnmarshal2Struct(testJSONFile, &testEaseActual)
	genre := testEaseActual.Genre

	assert.Emptyf(t, genre, "Genre should be empty: %s", genre)
	assert.Nil(t, err)

	//--------------------------------------------------------------------------------------------
	// nil object
	err = NewJsonify().Unmarshal2Struct(testJSONFile, nil)
	assert.EqualError(t, err, "obj cannot be nil.")

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

type jsonReaderMock struct {
	mock.Mock
}

func (r *jsonReaderMock) ReadAll(file *os.File) ([]byte, error) {
	args := r.Called(file)
	return args.Get(0).([]byte), args.Error(1)
}

func TestUnmarshal2Struct_ReadAll(t *testing.T) {
	//--------------------------------------------------------------------------------------------
	// Happy
	jsonReader := new(jsonReaderMock)
	mockReadError := errors.New("Error reading all")
	jsonReader.On("ReadAll", mock.Anything).Return([]byte{}, mockReadError)
	newJsonify := Jsonify{jsonReader}

	var testEaseTemp any = testEase
	err := newJsonify.Unmarshal2Struct(testJSONFile, &testEaseTemp)
	assert.ErrorIsf(t, err, mockReadError, "Mocked ReadAll should return \"%v\"", mockReadError)

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

func TestWriteObjToJSONFile_main(t *testing.T) {
	newJsonify := NewJsonify()

	//--------------------------------------------------------------------------------------------
	// Happy
	err := newJsonify.WriteObjToJSONFile(testJSONWriteFile, testEase)
	assert.Nilf(t, err, "Struct testEase should get written to jsonFilename without error")

	var testEaseActual TestEase
	err = test2AnyUnmarshal2Struct(testJSONWriteFile, &testEaseActual)
	assert.Nilf(t, err, "Should not have an error when unmarshaling %s to testEaseActual", testJSONWriteFile)
	assert.Equalf(t, testEase, testEaseActual, "testEaseActual was unmarshaled from the result of writing testEase to %s", testJSONWriteFile)

	//--------------------------------------------------------------------------------------------
	// nil obj
	err = newJsonify.WriteObjToJSONFile(testJSONWriteFile, nil)
	assert.EqualError(t, err, "obj cannot be nil.")

	//--------------------------------------------------------------------------------------------
	// Can't create file
	outFile := "/does/not/exist/" + testJSONWriteFile
	err = newJsonify.WriteObjToJSONFile(outFile, testEase)
	assert.NotNilf(t, err, "A file with a non-existant path should cause an error %s", outFile)

	//--------------------------------------------------------------------------------------------
	// Failed encode
	type fail struct {
		Numb    complex128 `json:"numb"`
		NewChan chan int   `json:"new_chan"`
	}

	failure := &fail{
		Numb:    complex(3.5, 2.1),
		NewChan: make(chan int),
	}

	err = newJsonify.WriteObjToJSONFile(testJSONWriteFile, &failure)
	theError := "json: unsupported type: complex128"
	assert.EqualErrorf(t, err, theError, "Expected \"%s\"", theError)

	//--------------------------------------------------------------------------------------------
	// Cleanup

}
