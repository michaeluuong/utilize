// Package jsoning provides convenience operations for JSON.
package jsoning

// Toggle coverage: control+option+command+t

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/michaeluuong/utilize/filing"
)

type reader interface {
	ReadAll(file *os.File) ([]byte, error)
}

type JsonReader struct{}

func (j JsonReader) ReadAll(file *os.File) ([]byte, error) {
	return io.ReadAll(file)

}

// Jsonify reads a file by wrapping io.ReadAll().
type Jsonify struct {
	reader
}

// NewJsonify provides a Jsonify object.
func NewJsonify() *Jsonify {
	return &Jsonify{reader: JsonReader{}}

}

// Unmarshal2Struct unmarshals a JSON object from a file into an object.
//   - jsonFilename is the full path to the file with the JSON object
//   - obj is the object to populate with the JSON object
//
// Return error if
//   - Unable to open file
//   - Unable to read file
//   - Unable to unmarshal JSON object
func (j Jsonify) Unmarshal2Struct(jsonFilename string, obj *any) error {
	if obj == nil {
		return errors.New("obj cannot be nil.")

	}

	jsonFile, err := os.Open(jsonFilename)
	if err != nil {
		return err

	}
	defer jsonFile.Close()

	// Global closure for better coverage
	byteValue, err := j.ReadAll(jsonFile)
	if err != nil {
		slog.Error("ReadAll()|could not read JSON file", "jsonFilename", jsonFilename, "err", err)
		return fmt.Errorf("Error reading file %s: %w", jsonFilename, err)

	}
	err = json.Unmarshal(byteValue, &obj)
	if err != nil {
		slog.Error("Unmarshall()|could not unmarshall JSON file to object", "jsonFilename", jsonFilename, "err", err)
		return fmt.Errorf("Error unmarshaling JSON: %w", err)

	}

	return nil

}

// WriteObjToJSONFile marshals an object to JSON and writes it to a file.
//   - jsonFilename is the name of the file to write the JSON object to
//   - obj is the object to encode into jsonFilename
//
// Return error if
//   - obj is nil
//   - Unable to create file
//   - Unable to encode obj
func (j Jsonify) WriteObjToJSONFile(jsonFilename string, obj any) error {
	if obj == nil {
		return errors.New("obj cannot be nil.")

	}

	jsonDirName := filepath.Dir(jsonFilename)
	if !filing.Exists(jsonDirName) {
		os.MkdirAll(jsonDirName, 0755)
		slog.Debug("os.MkdirAll()|created directory", "jsonDirName", jsonDirName, "jsonFilename", jsonFilename)

	}

	file, err := os.Create(jsonFilename)
	if err != nil {
		slog.Error("os.Create()|could not create file", "jsonFilename", jsonFilename, "err", err)
		return err

	}
	defer file.Close()
	slog.Debug("os.Create()|created file", "jsonFilename", jsonFilename)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(obj)
	if err != nil {
		slog.Error("encoder.Encode()|could not encode object as JSON to file", "jsonFilename", jsonFilename, "err", err)
		return err

	}

	return nil

}
