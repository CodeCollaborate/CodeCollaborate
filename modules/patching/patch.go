package patching

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

// Patch represents a set of changes to a versioned document
type Patch struct {
	// BaseVersion is the version that this patch was created on.
	BaseVersion int

	// Changes is the list of changes that were applied to the document.
	// When patching, changes MUST be applied in order.
	Changes []*Diff
}

// NewPatch creates a new patch with the given parameters
func NewPatch(baseVersion int, changes []*Diff) *Patch {
	return &Patch{
		BaseVersion: baseVersion,
		Changes:     changes,
	}
}

// NewPatchFromString parses a patch from its given string representation
func NewPatchFromString(str string) (*Patch, error) {
	var err error
	patch := Patch{}

	parts := strings.Split(str, ":\n")
	if len(parts) < 2{
		return nil, errors.New("Invalid patch format")
	}

	if len(parts[0]) <= 1 {
		return nil, errors.New("Invalid base version")
	}
	patch.BaseVersion, err = strconv.Atoi(string(parts[0][1]))

	if err != nil {
		return nil, err
	}

	diffStrs := strings.Split(parts[1], ",\n")

	if len(diffStrs) == 0 {
		return nil, errors.New("No changes attached to patch")
	}

	for _, diffStr := range diffStrs {
		newDiff, err := NewDiffFromString(diffStr)
		if err != nil {
			return nil, err
		}
		patch.Changes = append(patch.Changes, newDiff)
	}

	return &patch, nil
}

// ConvertToCRLF converts this patch from using LF to CRLF line separators given the base text to patch.
func (patch *Patch) ConvertToCRLF(base string) *Patch {
	newChanges := []*Diff{}

	for _, diff := range patch.Changes {
		newChanges = append(newChanges, diff.ConvertToCRLF(base))
	}

	return NewPatch(patch.BaseVersion, newChanges)
}

// ConvertToLF converts this patch from using CRLF to LF line separators given the base text to patch.
func (patch *Patch) ConvertToLF(base string) *Patch {
	newChanges := []*Diff{}

	for _, diff := range patch.Changes {
		newChanges = append(newChanges, diff.ConvertToLF(base))
	}

	return NewPatch(patch.BaseVersion, newChanges)
}

// Undo reverses this patch, producing a patch to undo the changes done by applying the patch.
func (patch *Patch) Undo() *Patch {
	newChanges := []*Diff{}

	// This needs to be in reverse order, since all the diffs in a package will have been applied in order.
	// The last diff will have been computed relative to the previous few.
	for i := len(patch.Changes) - 1; i >= 0; i-- {
		newChanges = append(newChanges, patch.Changes[i].Undo())
	}

	return NewPatch(patch.BaseVersion, newChanges)
}

// Transform does an Operational Transform against the other patches, creating a set
// of changes relative to previously applied changes.
func (patch *Patch) Transform(others []*Patch) *Patch {
	intermediateDiffs := patch.Changes
	maxVersionSeen := patch.BaseVersion

	for _, otherPatch := range others {
		newIntermediateDiffs := []*Diff{}

		for _, diff := range intermediateDiffs {
			newIntermediateDiffs = append(newIntermediateDiffs, diff.transform(otherPatch.Changes)...)
		}

		intermediateDiffs = newIntermediateDiffs
		if maxVersionSeen < otherPatch.BaseVersion {
			maxVersionSeen = otherPatch.BaseVersion
		}
	}

	return NewPatch(maxVersionSeen, intermediateDiffs)
}

func (patch *Patch) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("v")
	buffer.WriteString(strconv.Itoa(patch.BaseVersion))
	buffer.WriteString(":\n")
	buffer.WriteString(patch.Changes[0].String())
	for _, diff := range patch.Changes[1:] {
		buffer.WriteString(",\n")
		buffer.WriteString(diff.String())
	}

	return buffer.String()
}
