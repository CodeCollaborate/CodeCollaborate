package patching

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

// PatchTextFromString applies the provided patches onto the given text. The patches are applied strictly in the order given.
// This method completes in O(n*m) time, where n is the base text length, and m is the number of patches.
func PatchTextFromString(text string, patchesStr []string) (string, error) {
	patches := make([]*Patch, len(patchesStr))
	for i, patch := range patchesStr {
		patch, err := NewPatchFromString(patch)
		if err != nil {
			return "", err
		}
		patches[i] = patch
	}

	return PatchText(text, patches)
}

// ErrorIllegalLocation is the error thrown if a diff attempts to insert in an invalid location, such as between an \r and \n
var ErrorIllegalLocation = errors.New("Attempted to apply diff at an illegal lcoation")

// PatchText applies the provided patches onto the given text. The patches are applied strictly in the order given.
// This method completes in O(n*m) time, where n is the base text length, and m is the number of patches.
func PatchText(text string, patches []*Patch) (string, error) {
	useCRLF := strings.Contains(text, "\r\n")

	for _, patch := range patches {
		noOpLength := 0
		prevEndIndex := 0
		var prevDiff *Diff
		var buffer bytes.Buffer

		if useCRLF {
			patch.ConvertToCRLF(text)
		}

		for _, diff := range patch.Changes {
			if diff.StartIndex > 0 && diff.StartIndex < utf8.RuneCountInString(text) &&
				text[diff.StartIndex-1] == '\r' && text[diff.StartIndex] == '\n' {
				return "", ErrorIllegalLocation
			}

			noOpLength = diff.StartIndex
			if prevDiff != nil {
				// Min of (diff.StartIndex - prevDiff.StartIndex) and (diff.StartIndex - prevEndIndex)
				noOpLength = diff.StartIndex - prevDiff.StartIndex
				if diff.StartIndex-prevEndIndex < noOpLength {
					noOpLength = diff.StartIndex - prevEndIndex
				}
				// Max of 0, noOpLength
				if noOpLength < 0 {
					noOpLength = 0
				}
			}

			// Copy any text that is untouched
			if noOpLength > 0 {
				buffer.WriteString(text[prevEndIndex : prevEndIndex+noOpLength])
			}

			if diff.Insertion {
				// Commit insertion
				buffer.WriteString(diff.Changes)

				// End index is incremented only by the no-op length;
				// insertions do not change the index in the original text
				prevEndIndex += noOpLength
			} else {
				// Move to start of deletion
				prevEndIndex += noOpLength

				if text[prevEndIndex:prevEndIndex+diff.Length()] != diff.Changes {
					return "", fmt.Errorf("PatchText: Deleted text [%s] does not match changes in diff: [%s]", text[prevEndIndex:prevEndIndex+diff.Length()], diff.Changes)
				}
				// Skip past the text that is deleted
				prevEndIndex += diff.Length()
			}
			prevDiff = diff
		}

		// Copy the remainder
		if prevEndIndex < len(text) {
			buffer.WriteString(text[prevEndIndex:])
		}
		text = buffer.String()
	}

	return text, nil
}
