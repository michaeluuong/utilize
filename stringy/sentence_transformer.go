package stringy

import (
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/transform"
)

// SenetenceTransformer implements transform.SpanningTransformer in order to capitalize the first letter of a sentence.
type SentenceTransformer struct {
	transform.NopResetter
}

// Transfer capitalizes the first letter of a sentence.
//   - dst is the buffer that will contain the result of capitalizing the firt letter of src
//   - src is the sentence buffer to capitalize the first letter of
//   - atEOF is unused since we only need the first character and don't need to flush any buffers
//
// Return
//   - nDst is the number of bytes written into the dst buffer
//   - nSrc is the number of bytes consumed from the src buffer to produce nDst
//   - err contains ansform.ErrShortDst if the dst buffer was too small to hold the transformation
func (t *SentenceTransformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	if len(src) == 0 {
		return 0, 0, nil

	}

	r, size := utf8.DecodeRune(src)
	capitalized := string(unicode.ToTitle(r))
	capitalizedBytes := []byte(capitalized)

	if len(dst) < len(capitalizedBytes) {
		return 0, 0, transform.ErrShortDst

	}

	nDst += copy(dst, capitalizedBytes)
	nSrc += size

	if len(src[size:]) > 0 {
		if len(dst[nDst:]) < len(src[size:]) {
			return nDst, nSrc, transform.ErrShortDst

		}

		nDst += copy(dst[nDst:], src[size:])
		nSrc += len(src[size:])

	}

	return nDst, nSrc, nil

}

// Span determines how much of the source already meets the transformation.
//   - src is the buffer containing the sentence to tranform
//   - atEOF is unused because we only care about the first letter
//
// Return
//   - n is the size of the transformation which is a single rune or no runes if already capitalized
//   - err will never contain an error (i.e. always nil)
func (c *SentenceTransformer) Span(src []byte, atEOF bool) (n int, err error) {
	if len(src) == 0 {
		return 0, nil

	}

	r, size := utf8.DecodeRune(src)
	if unicode.IsTitle(r) || unicode.IsUpper(r) {
		return size, nil

	}

	return 0, nil

}
