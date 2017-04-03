package patching

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Patch represents a set of changes to a versioned document
type Patch struct {
	// BaseVersion is the version that this patch was created on.
	BaseVersion int64

	// Changes is the list of changes that were applied to the document.
	// When patching, changes MUST be applied in order.
	Changes Diffs

	// DocLength is the length of the document prior to the application of this patch
	DocLength int
}

// GetPatches creates an array of patches, given the array of strings
func GetPatches(patchStrs []string) ([]*Patch, error) {
	patches := make([]*Patch, len(patchStrs))
	for index, patchStr := range patchStrs {
		patch, err := NewPatchFromString(patchStr)
		if err != nil {
			return nil, err
		}
		patches[index] = patch
	}
	return patches, nil
}

// NewPatch creates a new patch with the given parameters
func NewPatch(baseVersion int64, changes Diffs, docLength int) *Patch {
	patch := &Patch{
		BaseVersion: baseVersion,
		Changes:     changes,
		DocLength:   docLength,
	}

	return patch.simplify()
}

// NewPatchFromString parses a patch from its given string representation
func NewPatchFromString(str string) (*Patch, error) {
	var err error
	patch := Patch{}

	parts := strings.Split(str, ":\n")
	if len(parts) < 3 {
		return nil, errors.New("Invalid patch format")
	}

	if len(parts[0]) <= 1 {
		return nil, errors.New("Invalid base version")
	}

	patch.BaseVersion, err = strconv.ParseInt(string(parts[0][1:]), 10, 64)
	if err != nil {
		return nil, err
	}

	docLen64, err := strconv.ParseInt(string(parts[2]), 10, 0)
	if err != nil {
		return nil, err
	}
	patch.DocLength = int(docLen64)

	diffStrs := strings.Split(parts[1], ",\n")

	for _, diffStr := range diffStrs {
		if len(diffStr) == 0 {
			continue
		}

		newDiff, err := NewDiffFromString(diffStr)
		if err != nil {
			return nil, err
		}
		patch.Changes = append(patch.Changes, newDiff)
	}

	return patch.simplify(), nil
}

// ConvertToCRLF converts this patch from using LF to CRLF line separators given the base text to patch.
func (patch *Patch) ConvertToCRLF(base string) *Patch {
	newChanges := Diffs{}

	for _, diff := range patch.Changes {
		newChanges = append(newChanges, diff.ConvertToCRLF(base))
	}

	return NewPatch(patch.BaseVersion, newChanges, utf8.RuneCountInString(strings.Replace(base, "\n", "\r\n", -1)))
}

// ConvertToLF converts this patch from using CRLF to LF line separators given the base text to patch.
func (patch *Patch) ConvertToLF(base string) *Patch {
	newChanges := Diffs{}

	for _, diff := range patch.Changes {
		newChanges = append(newChanges, diff.ConvertToLF(base))
	}

	return NewPatch(patch.BaseVersion, newChanges, utf8.RuneCountInString(strings.Replace(base, "\r\n", "\n", -1)))
}

func (patch *Patch) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("v")
	buffer.WriteString(fmt.Sprintf("%d", patch.BaseVersion))
	buffer.WriteString(":\n")
	if patch.Changes.Len() > 0 {
		buffer.WriteString(patch.Changes[0].String())
		for _, diff := range patch.Changes[1:] {
			buffer.WriteString(",\n")
			buffer.WriteString(diff.String())
		}
	}
	buffer.WriteString(fmt.Sprintf(":\n%d", patch.DocLength))

	return buffer.String()
}

func (patch *Patch) simplify() *Patch {
	patch.Changes = patch.Changes.Simplify()

	return patch
}
