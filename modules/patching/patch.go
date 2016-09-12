package patching

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

type Patch struct {
	BaseVersion int
	Changes     []*Diff
}

func NewPatch(baseVersion int, changes []*Diff) *Patch {
	return &Patch{
		BaseVersion: baseVersion,
		Changes:     changes,
	}
}

func NewPatchFromString(str string) (*Patch, error) {
	var err error
	patch := Patch{}

	parts := strings.Split(str, ":\n")

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

func (patch *Patch) ConvertToCRLF(base string) *Patch {
	newChanges := []*Diff{}

	for _, diff := range patch.Changes {
		newChanges = append(newChanges, diff.ConvertToCRLF(base))
	}

	return NewPatch(patch.BaseVersion, newChanges)
}

func (patch *Patch) ConvertToLF(base string) *Patch {
	newChanges := []*Diff{}

	for _, diff := range patch.Changes {
		newChanges = append(newChanges, diff.ConvertToLF(base))
	}

	return NewPatch(patch.BaseVersion, newChanges)
}

func (patch *Patch) GetUndo() *Patch {
	newChanges := []*Diff{}

	for _, diff := range patch.Changes {
		newChanges = append(newChanges, diff.GetUndo())
	}

	return NewPatch(patch.BaseVersion, newChanges)
}

func (patch *Patch) Transform(others []*Patch) *Patch {
	intermediateDiffs := patch.Changes
	maxVersionSeen := patch.BaseVersion

	for _, otherPatch := range others {

		newIntermediateDiffs := []*Diff{}

		for _, diff := range intermediateDiffs {
			newIntermediateDiffs = append(newIntermediateDiffs, diff.transform(otherPatch.Changes)...)
		}

		intermediateDiffs = newIntermediateDiffs
		if maxVersionSeen < patch.BaseVersion {
			maxVersionSeen = patch.BaseVersion
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
