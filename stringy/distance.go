// Package stringy collects convenience functions related to the string.
//
// * operations related to the Levenshtein distance between two strings
package stringy

import (
	"github.com/agnivade/levenshtein"
)

const Tolerance = 1e-6

// HasHighStrictness determines if the Levenshtein similarity score is above 0.9 (i.e. high) with a tolerance of -1e-6
//   - simScore is the Levenshtein similarity score for the distance between two string
//
// Return true if the similarity score is above (0.9 - 1e-6)
func HasHighStrictness(simScore float32) bool {
	if simScore >= (0.9 - Tolerance) {
		return true

	}

	return false

}

// HasModerateStrictness determines if the Levenshtein similarity score is above 0.8 with a tolerance of -1e-6.
//   - simScore is the Levenshtein similarity score for the distance between two string
//
// NOTE: high similarity scores will also be moderately strict.
//
// Return true if the similarity score is above (0.80 - 1e-6)
func HasModerateStrictness(simScore float32) bool {
	if simScore >= (0.80 - Tolerance) {
		return true

	}

	return false

}

// HasLowStrictness determines if the Levenshtein similarity score is below 0.75 with a tolerance of -1e-6.
//   - simScore is the Levenshtein similarity score for the distance between two string
//
// Return true if the similarity score is below (0.75 - 1e-6)
func HasLowStrictness(simScore float32) bool {
	if simScore < (0.75 - Tolerance) {
		return true

	}

	return false

}

// IsHighlyStrict determines if the two strings compare closely at a highly strict level.
//   - str1 is the string to compare against str2
//   - str2 is the string to compare against str1
//
// Return true if the similiary score is above .9-tolerance (1e-6) otherwise return false
func IsHighlyStrict(str1, str2 string) bool {
	sim := Similarity(str1, str2)
	if sim >= (0.9 - Tolerance) {
		return true

	}

	return false

}

// IsModeratelyStrict determines if the two strings compare closely at a moderately strict level.
//   - str1 is the string to compare against str2
//   - str2 is the string to compare against str1
//
// Note: a moderately strict equality can also be highly strict.
//
// Return true if the similiary score is above .8-tolerance (1e-6) otherwise return false
func IsModeratelyStrict(str1, str2 string) bool {
	sim := Similarity(str1, str2)
	if sim >= (0.8 - Tolerance) {
		return true

	}

	return false

}

// IsLowlyStrict determines if the two strings compare closely at a low strict level.
//   - str1 is the string to compare against str2
//   - str2 is the string to compare against str1
//
// Return true if the similiary score is below .75-tolerance (1e-6) otherwise return false
func IsLowlyStrict(str1, str2 string) bool {
	sim := Similarity(str1, str2)
	if sim < (0.75 - Tolerance) {
		return true

	}

	return false

}

// Similarity determines where the Levenshtein distance falls within a similary scale.
//   - str1 is the string to compare against str2
//   - str2 is the string to compare against str1
//
// Return a number between 0 and 1
func Similarity(str1, str2 string) float32 {
	distance := float32(LDistance(str1, str2))
	return 1.0 - (distance / (float32(max(len(str1), len(str2)))))

}

// Distance computes the Levenshtein distance between two strings.
//   - str1 is the first or source string to compare
//   - str2 is the string to compare against the first
//
// Return the number of changes that would need to take place to change one string into another.
func LDistance(str1, str2 string) int {
	return levenshtein.ComputeDistance(str1, str2)

}
