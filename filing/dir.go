// Package filing collects functions that are reliant on the file system.
//
// * operations related to directories
package filing

import (
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	DirPerm os.FileMode = 0755
)

type AncestorLevel int

const (
	CurrentLevel = iota
	ParentLevel
	GrandparentLevel
	ParentGrandparentLevel
)

// Last returns the value in the final index of parts.
//   - parts is the slice to get the last index of
//
// Return the final index of []T or an empty T value.
func last[T any](parts []T) T {
	partsLen := len(parts)
	if parts == nil || partsLen == 0 {
		var empty T
		return empty

	}

	return parts[partsLen-1]

}

// lastStringPart gets the last index of the string separated by the separator.
//   - splitString is the string separated by sep to get the last part of
//   - sep is the separator charactore to separate splitString by
//
// Return the last part of splitString after separated by sep or splitString itself if it contains zero separator characters.
func lastStringPart(splitString, sep string) string {
	if splitString == "" || sep == "" {
		return splitString

	}

	stringParts := strings.Split(splitString, sep)
	return stringParts[len(stringParts)-1]

}

// RemoveLast gets rid of everything after the final separator.
//   - splitString is the string to get rid of everything after the final separator
//   - sep is the separator charactore to separate splitString by
//
// Return splitString with everything after the final separator removed or splitString itself it it contains zero separator characters.
func RemoveLast(splitString, sep string) string {
	if splitString == "" || sep == "" {
		return splitString

	}

	stringParts := strings.Split(splitString, sep)
	if len(stringParts) == 1 {
		return stringParts[0]

	}

	return strings.Join(stringParts[:len(stringParts)-1], sep)

}

// NextDir will try to find the next directory name so we don't clobber the existing one.
//   - dirPath is the directory to separate (i.e. absolute path, directory name, filename, etc.)
//   - partSepOpt is the optional directory name separator (default is _, e.g. old_directory_1)
//
// Return
//   - the next integer number if the directory already ended in an integer or 1 if this is the first
//   - an error if there was a problem listing files or finding the next suffix
func NextDir(dirPath string, partSepOpt ...string) (string, error) {
	newDir := dirPath

	var partSep string = "_"
	if len(partSepOpt) > 0 {
		partSep = partSepOpt[0]

	}

	if dirPath != "" && Exists(dirPath) {
		dirName := filepath.Dir(dirPath)
		baseDirName := filepath.Base(dirPath)

		baseDirNameStrip := baseDirName
		if matched, _ := regexp.MatchString(partSep+"[0-9]+$", baseDirName); matched {
			baseDirNameStripTmp := RemoveLast(baseDirName, partSep)
			// Try to prevent similar but not exact matches
			if matched, _ := regexp.MatchString(baseDirNameStripTmp+partSep+"[0-9]+$", baseDirName); matched {
				baseDirNameStrip = baseDirNameStripTmp

			}

		}

		dirnameRegex := "^" + regexp.QuoteMeta(baseDirNameStrip) + "(" + partSep + "[0-9]*)?$"
		dirList, err := Ls(dirName, dirnameRegex)
		if err != nil {
			return "", err

		}

		dirs := DirEntryName(dirList)
		// Sort by numerical suffix
		sort.Slice(dirs, func(i, j int) bool {
			iInt, _ := strconv.Atoi(lastStringPart(dirs[i], partSep))
			jInt, _ := strconv.Atoi(lastStringPart(dirs[j], partSep))
			return iInt < jInt

		})
		lastDir := last(dirs)
		suffix := lastStringPart(lastDir, partSep)

		var newSuffix int
		if fileSuffix, err := strconv.ParseInt(suffix, 10, 32); err == nil {
			newSuffix = int(fileSuffix)

		} else {
			newSuffix = 1

		}

		// Find an opening but backwards to be annoying
		if len(dirs) < newSuffix && newSuffix > 1 {
			dirNameSet := DirEntryNameSet(dirList)
			for i := newSuffix; i > 2; i-- {
				key := baseDirNameStrip + partSep + strconv.Itoa(i)
				if _, ok := dirNameSet[key]; ok {
					newSuffix--

				} else {
					break

				}

			}

		} else {
			newSuffix++

		}

		newDir = filepath.Join(dirName, baseDirNameStrip+"_"+strconv.Itoa(newSuffix))
		slog.Debug("vars|finding next directory", "baseDirNameStrip", baseDirNameStrip, "lastDir", lastDir, "suffix", suffix, "newDir", newDir)

	}

	return newDir, nil

}

// IsDir determines if the path is a directory.
//   - dirPath should be the absolute path to the file/folder
//
// Return
//   - true if dirPath is a directory and false if it is not
//   - error if os.Stat() fails
func IsDir(dirPath string) (bool, error) {
	isDir := false
	fileInfo, err := os.Stat(dirPath)
	if err == nil && fileInfo.IsDir() {
		isDir = true

	}

	return isDir, err

}
