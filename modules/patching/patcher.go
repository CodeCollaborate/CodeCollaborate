package patching

import (
	"bytes"
	"fmt"
)

// PatchText applies the provided patches onto the given text. The patches are applied strictly in the order given.
// This method completes in O(n*m) time, where n is the base text length, and m is the number of patches.
func BuildAndPatchText(text string, patchesStr []string) (string, error) {
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

// PatchText applies the provided patches onto the given text. The patches are applied strictly in the order given.
// This method completes in O(n*m) time, where n is the base text length, and m is the number of patches.
func PatchText(text string, patches []*Patch) (string, error) {
	for _, patch := range patches {
		var buffer bytes.Buffer
		startIndex := 0
		for _, diff := range patch.Changes {
			// Copy anything before the changes
			if startIndex < diff.StartIndex {
				buffer.WriteString(text[startIndex:diff.StartIndex])
			}
			if diff.Insertion {
				// insert item
				buffer.WriteString(diff.Changes)

				// If the diff's startIndex is greater, move it up.
				// Otherwise, a previous delete may have deleted over the start index.
				if startIndex < diff.StartIndex {
					startIndex = diff.StartIndex
				}
			} else {
				// validate that we're deleting the right characters
				if want, got := diff.Changes, text[diff.StartIndex:diff.StartIndex+diff.Length()]; want != got {
					return "", fmt.Errorf("PatchText: Deleted text %q does not match changes in diff: %q", got, want)
				}

				// shift the start index of the next round
				startIndex = diff.StartIndex + diff.Length()
			}
		}
		// Copy the remainder
		buffer.WriteString(text[startIndex:])
		text = buffer.String()
	}

	return text, nil
}
