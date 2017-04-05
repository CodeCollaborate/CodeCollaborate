package patching

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetNextDiff(t *testing.T) {
	tests := []struct {
		desc          string
		patch         *Patch
		currIndex     int
		wasNoOp       bool
		expectedDiff  string
		expectedIndex int
	}{
		{
			desc:          "Starting noOp",
			patch:         GetPatchesOrDie(t, "v0:2:\n1:+1:A:\n5")[0],
			currIndex:     -1,
			wasNoOp:       false,
			expectedDiff:  "0:+0:",
			expectedIndex: -1,
		},
		{
			desc:          "First patch after noOp",
			patch:         GetPatchesOrDie(t, "v0:2:\n1:+1:A:\n5")[0],
			currIndex:     -1,
			wasNoOp:       true,
			expectedDiff:  "1:+1:A",
			expectedIndex: 0,
		},
		{
			desc:          "Ending noOp",
			patch:         GetPatchesOrDie(t, "v0:2:\n1:+1:A:\n5")[0],
			currIndex:     0,
			wasNoOp:       false,
			expectedDiff:  "1:+0:",
			expectedIndex: 0,
		},
		{
			desc:          "Insertion noOp",
			patch:         GetPatchesOrDie(t, "v0:2:\n1:+1:A:\n5")[0],
			currIndex:     0,
			wasNoOp:       false,
			expectedDiff:  "1:+0:",
			expectedIndex: 0,
		},
		{
			desc:          "Removal noOp",
			patch:         GetPatchesOrDie(t, "v0:2:\n1:-1:A:\n5")[0],
			currIndex:     0,
			wasNoOp:       false,
			expectedDiff:  "2:+0:",
			expectedIndex: 0,
		},
		{
			desc:          "Beyond ending noOp",
			patch:         GetPatchesOrDie(t, "v0:2:\n1:+1:A:\n5")[0],
			currIndex:     0,
			wasNoOp:       true,
			expectedDiff:  "nil",
			expectedIndex: -1,
		},
		{
			desc:          "Starting noOp with first diff at 0",
			patch:         GetPatchesOrDie(t, "v0:2:\n0:+1:A:\n5")[0],
			currIndex:     -1,
			wasNoOp:       false,
			expectedDiff:  "0:+0:",
			expectedIndex: -1,
		},
		{
			desc:          "Ending noOp after noOp with first diff at 0",
			patch:         GetPatchesOrDie(t, "v0:2:\n0:+1:A:\n5")[0],
			currIndex:     0,
			wasNoOp:       false,
			expectedDiff:  "0:+0:",
			expectedIndex: 0,
		},
		{
			desc:          "Ending noOp with no space after",
			patch:         GetPatchesOrDie(t, "v0:2:\n0:+1:A:\n0")[0],
			currIndex:     0,
			wasNoOp:       false,
			expectedDiff:  "0:+0:",
			expectedIndex: 0,
		},
		{
			desc:          "Beyond ending noOp",
			patch:         GetPatchesOrDie(t, "v0:2:\n0:+1:A:\n5")[0],
			currIndex:     0,
			wasNoOp:       true,
			expectedDiff:  "nil",
			expectedIndex: -1,
		},
		{
			desc:          "Adjacent patch insert-insert 1",
			patch:         GetPatchesOrDie(t, "v0:2:\n0:+1:A,\n1:+1:B:\n5")[0],
			currIndex:     0,
			wasNoOp:       false,
			expectedDiff:  "0:+0:",
			expectedIndex: 0,
		},
		{
			desc:          "Adjacent patch insert-insert 2",
			patch:         GetPatchesOrDie(t, "v0:2:\n0:+1:A,\n1:+1:B:\n5")[0],
			currIndex:     0,
			wasNoOp:       true,
			expectedDiff:  "1:+1:B",
			expectedIndex: 1,
		},
		{
			desc:          "Adjacent patch insert-remove 1",
			patch:         GetPatchesOrDie(t, "v0:2:\n0:+1:A,\n1:-1:B:\n5")[0],
			currIndex:     0,
			wasNoOp:       false,
			expectedDiff:  "0:+0:",
			expectedIndex: 0,
		},
		{
			desc:          "Adjacent patch insert-remove 2",
			patch:         GetPatchesOrDie(t, "v0:2:\n0:+1:A,\n1:-1:B:\n5")[0],
			currIndex:     0,
			wasNoOp:       true,
			expectedDiff:  "1:-1:B",
			expectedIndex: 1,
		},
		{
			desc:          "Adjacent patch remove-insert",
			patch:         GetPatchesOrDie(t, "v0:2:\n0:-1:A,\n1:+1:B:\n5")[0],
			currIndex:     0,
			wasNoOp:       false,
			expectedDiff:  "1:+1:B",
			expectedIndex: 1,
		},
	}

	for _, test := range tests {
		diff, index := getNextDiff(test.patch, test.currIndex, test.wasNoOp)

		if test.expectedDiff == "nil" {
			assert.Nil(t, diff, "TestGetNextDiff[%s]: Unexpected nil diff", test.desc)
		} else {
			assert.Equal(t, test.expectedDiff, diff.String(), "TestGetNextDiff[%s]: Unexpected diff", test.desc)
		}
		assert.Equal(t, test.expectedIndex, index, "TestGetNextDiff[%s]: Unexpected index", test.desc)
	}
}

type overallConsolidationTest struct {
	desc     string
	baseText string
	patches  []*Patch
	error    error
}

func TestConsolidatePatch(t *testing.T) {
	tests := []overallConsolidationTest{
		{
			desc:     "Simple Add-Only test",
			baseText: "",
			patches:  GetPatchesOrDie(t, "v0:1:\n0:+7:testing:\n0", "v1:2:\n1:+2:AB,\n4:+4:CDEF,\n7:+3:GHI:\n7"),
		},
		{
			desc:     "Simple deletion-addition test 1",
			baseText: "testing",
			patches:  GetPatchesOrDie(t, "v0:1:\n2:-2:st,\n5:-1:n:\n7", "v1:2:\n1:+2:AB,\n4:+4:CDEF:\n4"),
		},
		{
			desc:     "Simple deletion-addition test 2",
			baseText: "testing",
			patches:  GetPatchesOrDie(t, "v0:1:\n2:-2:st,\n5:-1:n:\n7", "v1:2:\n0:-2:te,\n4:+4:CDEF:\n4"),
		},
		{
			desc:     "Mixed deletion-addition test 1",
			baseText: "testing",
			patches:  GetPatchesOrDie(t, "v0:1:\n1:+2:AB,\n5:-1:n:\n7", "v1:2:\n0:-2:tA,\n4:+4:CDEF:\n8"),
		},
		{
			desc:     "Mixed deletion-addition test 2",
			baseText: "testing",
			patches:  GetPatchesOrDie(t, "v0:1:\n1:-3:est,\n5:+2:AB:\n6", "v1:2:\n0:-2:ti,\n3:+4:CDEF:\n9"),
		},
	}

	for _, test := range tests {
		patchedText, err := PatchText(test.baseText, test.patches)
		require.Nil(t, err)

		consolidatedPatch, err := ConsolidatePatches(test.patches)
		require.Nil(t, err)

		consolidatedPatchedText, err := PatchText(test.baseText, []*Patch{consolidatedPatch})
		require.Equal(t, patchedText, consolidatedPatchedText, "TestConsolidatePatch[%s]: Expected %s but got %s", test.desc, patchedText, consolidatedPatchedText)
	}
}

// TODO(wongb): Add more extensive testing
