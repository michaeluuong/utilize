package stringy

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLDistance_main(t *testing.T) {
	var str1, str2 string
	var dist int

	//--------------------------------------------------------------------------------------------
	// Happy
	str1, str2 = "test", "test"
	dist = LDistance(str1, str2)
	assert.Equal(t, 0, dist)

	//--------------------------------------------------------------------------------------------
	// 1 letter diff
	str1, str2 = "test", "testy"
	dist = LDistance(str1, str2)
	assert.Equal(t, 1, dist)

	//--------------------------------------------------------------------------------------------
	// 2 letter diff
	str1, str2 = "test", "tasty"
	dist = LDistance(str1, str2)
	assert.Equal(t, 2, dist)

	//--------------------------------------------------------------------------------------------
	// Completely different
	str1, str2 = "test ease", "hardly_tool"
	dist = LDistance(str1, str2)
	minLen := max(len(str1), len(str2))
	assert.Equal(t, minLen, dist)

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

func TestSimilarity_HasStrictness_main(t *testing.T) {
	var str1, str2 string
	var similarity float32
	var expected float32

	var hasHighStrictness, hasModerateStrictness, hasLowStrictness bool
	//--------------------------------------------------------------------------------------------
	// Happy
	str1, str2 = "test", "test"
	similarity = Similarity(str1, str2)
	fmt.Printf("similarity: %f\n", similarity)
	expected = 1.000000
	assert.Equal(t, expected, similarity)

	hasHighStrictness = HasHighStrictness(similarity)
	assert.Equal(t, true, hasHighStrictness)

	hasLowStrictness = HasLowStrictness(similarity)
	assert.Equal(t, false, hasLowStrictness)

	//--------------------------------------------------------------------------------------------
	// Completely different by case
	str1, str2 = "test", "TEST"
	similarity = Similarity(str1, str2)
	expected = 0
	assert.Equal(t, expected, similarity)

	hasHighStrictness = HasHighStrictness(similarity)
	assert.Equal(t, false, hasHighStrictness)

	hasModerateStrictness = HasModerateStrictness(similarity)
	assert.Equal(t, false, hasModerateStrictness)

	hasLowStrictness = HasLowStrictness(similarity)
	assert.Equal(t, true, hasLowStrictness)

	//--------------------------------------------------------------------------------------------
	// 1 letter diff
	str1, str2 = "test", "testy"
	similarity = Similarity(str1, str2)
	expected = .80
	assert.Equal(t, expected, similarity)

	hasHighStrictness = HasHighStrictness(similarity)
	assert.Equal(t, false, hasHighStrictness)

	hasModerateStrictness = HasModerateStrictness(similarity)
	assert.Equal(t, true, hasModerateStrictness)

	//--------------------------------------------------------------------------------------------
	// 2 letter diff
	str1, str2 = "test", "tasty"
	similarity = Similarity(str1, str2)
	expected = .6000000
	assert.Equal(t, expected, similarity)

	hasHighStrictness = HasHighStrictness(similarity)
	assert.Equal(t, false, hasHighStrictness)

	hasModerateStrictness = HasModerateStrictness(similarity)
	assert.Equal(t, false, hasModerateStrictness)

	hasLowStrictness = HasLowStrictness(similarity)
	assert.Equal(t, true, hasLowStrictness)

	//--------------------------------------------------------------------------------------------
	// Completely different
	str1, str2 = "test ease", "hardly_tool"
	similarity = Similarity(str1, str2)
	expected = 0
	assert.Equal(t, expected, similarity)

	hasLowStrictness = HasLowStrictness(similarity)
	assert.Equal(t, true, hasLowStrictness)

	//--------------------------------------------------------------------------------------------
	// Cleanup
}

func TestIsHighlyStrict_main(t *testing.T) {
	var str1, str2 string
	var isHighlyStrict bool

	//--------------------------------------------------------------------------------------------
	// Happy
	str1, str2 = "test", "test"
	isHighlyStrict = IsHighlyStrict(str1, str2)
	assert.Equal(t, true, isHighlyStrict)

	//--------------------------------------------------------------------------------------------
	// 1 letter diff
	str1, str2 = "test", "testy"
	isHighlyStrict = IsHighlyStrict(str1, str2)
	assert.Equal(t, false, isHighlyStrict)

	//--------------------------------------------------------------------------------------------
	// 2 letter diff
	str1, str2 = "test", "tasty"
	isHighlyStrict = IsHighlyStrict(str1, str2)
	assert.Equal(t, false, isHighlyStrict)

	//--------------------------------------------------------------------------------------------
	// Completely different
	str1, str2 = "test ease", "hardly_tool"
	isHighlyStrict = IsHighlyStrict(str1, str2)
	assert.Equal(t, false, isHighlyStrict)

	//--------------------------------------------------------------------------------------------
	// Cleanup
}

func TestIsModeratelyStrict_main(t *testing.T) {
	var str1, str2 string
	var isModeratelyStrict bool

	//--------------------------------------------------------------------------------------------
	// Happy
	str1, str2 = "test", "test"
	isModeratelyStrict = IsModeratelyStrict(str1, str2)
	assert.Equal(t, true, isModeratelyStrict)

	//--------------------------------------------------------------------------------------------
	// 1 letter diff
	str1, str2 = "test", "testy"
	isModeratelyStrict = IsModeratelyStrict(str1, str2)
	assert.Equal(t, true, isModeratelyStrict)

	//--------------------------------------------------------------------------------------------
	// 2 letter diff
	str1, str2 = "test", "tasty"
	isModeratelyStrict = IsModeratelyStrict(str1, str2)
	assert.Equal(t, false, isModeratelyStrict)

	//--------------------------------------------------------------------------------------------
	// Completely different
	str1, str2 = "test ease", "hardly_tool"
	isModeratelyStrict = IsModeratelyStrict(str1, str2)
	assert.Equal(t, false, isModeratelyStrict)

	//--------------------------------------------------------------------------------------------
	// Cleanup
}

func TestIsLowlyStrict_main(t *testing.T) {
	var str1, str2 string
	var isLowlyStrict bool

	//--------------------------------------------------------------------------------------------
	// Happy
	str1, str2 = "test", "test"
	isLowlyStrict = IsLowlyStrict(str1, str2)
	assert.Equal(t, false, isLowlyStrict)

	//--------------------------------------------------------------------------------------------
	// 1 letter diff
	str1, str2 = "test", "testy"
	isLowlyStrict = IsLowlyStrict(str1, str2)
	assert.Equal(t, false, isLowlyStrict)

	//--------------------------------------------------------------------------------------------
	// 2 letter diff
	str1, str2 = "test", "tasty"
	isLowlyStrict = IsLowlyStrict(str1, str2)
	assert.Equal(t, true, isLowlyStrict)

	//--------------------------------------------------------------------------------------------
	// Completely different
	str1, str2 = "test ease", "hardly_tool"
	isLowlyStrict = IsLowlyStrict(str1, str2)
	assert.Equal(t, true, isLowlyStrict)

	//--------------------------------------------------------------------------------------------
	// Cleanup
}
