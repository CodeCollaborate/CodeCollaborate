package patching

import (
	"testing"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/require"
)

func TestPatch_NewPatch(t *testing.T) {
	patchString := "v1:\n3:-8:deletion,\n2:+6:insert:\n11"

	diff1 := NewDiff(false, 3, "deletion")
	diff2 := NewDiff(true, 2, "insert")
	patch := NewPatch(1, Diffs{diff1, diff2}, 11)
	require.Equal(t, patchString, patch.String())
}

func TestPatch_NewPatchFromString(t *testing.T) {
	patch, err := NewPatchFromString("v6:\n3:-8:deletion,\n2:+6:insert:\n11")
	require.Nil(t, err)
	require.Equal(t, int64(6), patch.BaseVersion)
	require.Equal(t, 2, len(patch.Changes))
	require.Equal(t, "3:-8:deletion", patch.Changes[0].String())
	require.Equal(t, "2:+6:insert", patch.Changes[1].String())
	require.Equal(t, 11, patch.DocLength)

	// Test insertion from string
	patch, err = NewPatchFromString("v4:\n2:+1:a:\n12")
	require.Nil(t, err)
	require.Equal(t, int64(4), patch.BaseVersion)
	require.Equal(t, 1, len(patch.Changes))
	require.Equal(t, "2:+1:a", patch.Changes[0].String())
	require.Equal(t, 12, patch.DocLength)

	// Test insertion from string
	patch, err = NewPatchFromString("v3:\n26:+2:ab,\n81:-3:cde:\n13")
	require.Nil(t, err)
	require.Equal(t, int64(3), patch.BaseVersion)
	require.Equal(t, 2, len(patch.Changes))
	require.Equal(t, "26:+2:ab", patch.Changes[0].String())
	require.Equal(t, "81:-3:cde", patch.Changes[1].String())
	require.Equal(t, 13, patch.DocLength)
}

func TestPatch_NewPatchFromStringInvalidFormats(t *testing.T) {

	_, err := NewPatchFromString("test")
	require.NotNil(t, err, "Did not throw an error on invalid format")

	_, err = NewPatchFromString("v1:\n0:@1:test")
	require.NotNil(t, err, "Did not throw an error on invalid diff operation type")

	_, err = NewPatchFromString("0:+1:test")
	require.NotNil(t, err, "Did not throw an error on wrong changes length")

	_, err = NewPatchFromString("a:+4:test")
	require.NotNil(t, err, "Did not throw an error on invalid offset")

	_, err = NewPatchFromString("0:+err:test")
	require.NotNil(t, err, "Did not throw an erorr on invalid length.")

	_, err = NewPatchFromString("v:\n0:+2:te")
	require.NotNil(t, err, "Did not throw an error on invalid baseVersion")

	_, err = NewPatchFromString("va:\n0:+2:te")
	require.NotNil(t, err, "Did not throw an error on invalid baseVersion")

	_, err = NewPatchFromString("v1:\n")
	require.NotNil(t, err, "Did not throw an error on empty changes")
}

func TestPatch_ConvertToCRLF(t *testing.T) {
	patch, err := NewPatchFromString("v0:\n0:+5:test%0A:\n12")
	require.Nil(t, err)
	newPatch := patch.ConvertToCRLF("\r\ntest")
	require.Equal(t, 1, len(newPatch.Changes))
	require.Equal(t, "v0:\n0:+6:test%0D%0A:\n7", newPatch.String())

	patch, err = NewPatchFromString("v0:\n1:+5:test%0A:\n12")
	require.Nil(t, err)
	newPatch = patch.ConvertToCRLF("\r\ntest")
	require.Equal(t, 1, len(newPatch.Changes))
	require.Equal(t, "v0:\n2:+6:test%0D%0A:\n7", newPatch.String())

	patch, err = NewPatchFromString("v0:\n2:+5:test%0A:\n12")
	require.Nil(t, err)
	newPatch = patch.ConvertToCRLF("\r\ntest")
	require.Equal(t, 1, len(newPatch.Changes))
	require.Equal(t, "v0:\n3:+6:test%0D%0A:\n7", newPatch.String())

	patch, err = NewPatchFromString("v0:\n7:+5:test%0A:\n12")
	require.Nil(t, err)
	newPatch = patch.ConvertToCRLF("\r\ntes\r\nt")
	require.Equal(t, 1, len(newPatch.Changes))
	require.Equal(t, "v0:\n9:+6:test%0D%0A:\n10", newPatch.String())

	patch, err = NewPatchFromString("v0:\n2:+5:test%0A,\n7:+5:test%0A:\n12")
	require.Nil(t, err)
	newPatch = patch.ConvertToCRLF("\r\ntes\r\nt")
	require.Equal(t, 2, len(newPatch.Changes))
	require.Equal(t, "v0:\n3:+6:test%0D%0A,\n9:+6:test%0D%0A:\n10", newPatch.String())

	patch, err = NewPatchFromString("v0:\n2:+5:test%0A,\n7:+5:test%0A,\n0:+5:test%0A:\n12")
	require.Nil(t, err)
	newPatch = patch.ConvertToCRLF("\r\ntes\r\nt")
	require.Equal(t, 3, len(newPatch.Changes))
	require.Equal(t, "v0:\n3:+6:test%0D%0A,\n9:+6:test%0D%0A,\n0:+6:test%0D%0A:\n10", newPatch.String())
}

func TestPatch_ConvertToLF(t *testing.T) {
	patch, err := NewPatchFromString("v0:\n0:+6:test%0D%0A:\n6")
	require.Nil(t, err)
	newPatch := patch.ConvertToLF("\r\ntest")
	require.Equal(t, 1, len(newPatch.Changes))
	require.Equal(t, "v0:\n0:+5:test%0A:\n5", newPatch.String())

	patch, err = NewPatchFromString("v0:\n2:+6:test%0D%0A:\n6")
	require.Nil(t, err)
	newPatch = patch.ConvertToLF("\r\ntest")
	require.Equal(t, 1, len(newPatch.Changes))
	require.Equal(t, "v0:\n1:+5:test%0A:\n5", newPatch.String())

	patch, err = NewPatchFromString("v0:\n3:+6:test%0D%0A:\n6")
	require.Nil(t, err)
	newPatch = patch.ConvertToLF("\r\ntest")
	require.Equal(t, 1, len(newPatch.Changes))
	require.Equal(t, "v0:\n2:+5:test%0A:\n5", newPatch.String())

	patch, err = NewPatchFromString("v0:\n9:+6:test%0D%0A:\n8")
	require.Nil(t, err)
	newPatch = patch.ConvertToLF("\r\ntes\r\nt")
	require.Equal(t, 1, len(newPatch.Changes))
	require.Equal(t, "v0:\n7:+5:test%0A:\n6", newPatch.String())

	patch, err = NewPatchFromString("v0:\n3:+6:test%0D%0A,\n9:+6:test%0D%0A:\n8")
	require.Nil(t, err)
	newPatch = patch.ConvertToLF("\r\ntes\r\nt")
	require.Equal(t, 2, len(newPatch.Changes))
	require.Equal(t, "v0:\n2:+5:test%0A,\n7:+5:test%0A:\n6", newPatch.String())

	patch, err = NewPatchFromString("v0:\n3:+6:test%0D%0A,\n9:+6:test%0D%0A,\n0:+6:test%0D%0A:\n8")
	require.Nil(t, err)
	newPatch = patch.ConvertToLF("\r\ntes\r\nt")
	require.Equal(t, 3, len(newPatch.Changes))
	require.Equal(t, "v0:\n2:+5:test%0A,\n7:+5:test%0A,\n0:+5:test%0A:\n6", newPatch.String())
}

func TestPatch_Transform(t *testing.T) {
	// Test set 1
	patch1, err := NewPatchFromString("v1:\n0:-1:a:\n10")
	require.Nil(t, err)
	patch2, err := NewPatchFromString("v0:\n3:-8:deletion,\n3:+6:insert:\n10")
	require.Nil(t, err)
	newPatch := patch2.Transform([]*Patch{patch1}, true)
	require.Equal(t, 2, len(newPatch.Changes))
	require.Equal(t, "v2:\n2:-8:deletion,\n2:+6:insert:\n9", newPatch.String())

	// Test set 2
	patch1, err = NewPatchFromString("v1:\n0:-1:a:\n10")
	require.Nil(t, err)
	patch2, err = NewPatchFromString("v2:\n0:-1:b:\n10")
	require.Nil(t, err)
	patch3, err := NewPatchFromString("v0:\n3:-8:deletion,\n3:+6:insert:\n10")
	require.Nil(t, err)
	newPatch = patch3.Transform([]*Patch{patch1, patch2}, true)
	require.Equal(t, 2, len(newPatch.Changes))
	require.Equal(t, "v3:\n1:-8:deletion,\n1:+6:insert:\n8", newPatch.String())
}

func TestPatch_Simplify(t *testing.T) {
	tests := []struct {
		desc     string
		patchStr string
		expected string
	}{
		{
			desc:     "Double-Insert, Adjacent",
			patchStr: "v1:\n0:+1:a,\n1:+1:b:\n10",
			expected: "v1:\n0:+1:a,\n1:+1:b:\n10",
		},
		{
			desc:     "Double-Remove, Adjacent",
			patchStr: "v1:\n0:-1:a,\n1:-1:b:\n10",
			expected: "v1:\n0:-2:ab:\n10",
		},
		{
			desc:     "Insert-Remove, Adjacent",
			patchStr: "v1:\n0:+1:a,\n1:-1:b:\n10",
			expected: "v1:\n0:+1:a,\n1:-1:b:\n10",
		},
		{
			desc:     "Remove-Insert, Adjacent",
			patchStr: "v1:\n0:-1:a,\n1:+1:b:\n10",
			expected: "v1:\n0:-1:a,\n1:+1:b:\n10",
		},
		{
			desc:     "Double-Insert, Not adjacent",
			patchStr: "v1:\n0:+1:a,\n2:+1:b:\n10",
			expected: "v1:\n0:+1:a,\n2:+1:b:\n10",
		},
		{
			desc:     "Double-Remove, Not adjacent",
			patchStr: "v1:\n0:-1:a,\n2:-1:b:\n10",
			expected: "v1:\n0:-1:a,\n2:-1:b:\n10",
		},
		{
			desc:     "Insert-Remove, Not adjacent",
			patchStr: "v1:\n0:+1:a,\n2:-1:b:\n10",
			expected: "v1:\n0:+1:a,\n2:-1:b:\n10",
		},
		{
			desc:     "Remove-Insert, Not adjacent",
			patchStr: "v1:\n0:-1:a,\n2:+1:b:\n10",
			expected: "v1:\n0:-1:a,\n2:+1:b:\n10",
		},
		{
			desc:     "Triple-Insert, Adjacent",
			patchStr: "v1:\n0:+1:a,\n1:+1:b,\n2:+1:c:\n10",
			expected: "v1:\n0:+1:a,\n1:+1:b,\n2:+1:c:\n10",
		},
		{
			desc:     "Triple-Remove, Adjacent",
			patchStr: "v1:\n0:-1:a,\n1:-1:b,\n2:-1:c:\n10",
			expected: "v1:\n0:-3:abc:\n10",
		},
		{
			desc:     "Double-Insert, Single Remove, Adjacent",
			patchStr: "v1:\n0:+1:a,\n1:+1:b,\n2:-1:c:\n10",
			expected: "v1:\n0:+1:a,\n1:+1:b,\n2:-1:c:\n10",
		},
		{
			desc:     "Single-Remove, Double-Insert, Adjacent",
			patchStr: "v1:\n0:-1:a,\n1:+1:b,\n2:+1:c:\n10",
			expected: "v1:\n0:-1:a,\n1:+1:b,\n2:+1:c:\n10",
		},
		{
			desc:     "Double-Remove, Single Insert, Adjacent",
			patchStr: "v1:\n0:-1:a,\n1:-1:b,\n2:+1:c:\n10",
			expected: "v1:\n0:-2:ab,\n2:+1:c:\n10",
		},
		{
			desc:     "Single-Insert, Double-Remove, Adjacent",
			patchStr: "v1:\n0:+1:a,\n1:-1:b,\n2:-1:c:\n10",
			expected: "v1:\n0:+1:a,\n1:-2:bc:\n10",
		},
		{
			desc:     "Double-Insert, Single Remove, Not adjacent",
			patchStr: "v1:\n0:+1:a,\n2:+1:b,\n3:-1:c:\n10",
			expected: "v1:\n0:+1:a,\n2:+1:b,\n3:-1:c:\n10",
		},
		{
			desc:     "Single-Remove, Double-Insert, Not adjacent",
			patchStr: "v1:\n0:-1:a,\n1:+1:b,\n3:+1:c:\n10",
			expected: "v1:\n0:-1:a,\n1:+1:b,\n3:+1:c:\n10",
		},
		{
			desc:     "Double-Remove, Single Insert, Not adjacent",
			patchStr: "v1:\n0:-1:a,\n2:-1:b,\n3:+1:c:\n10",
			expected: "v1:\n0:-1:a,\n2:-1:b,\n3:+1:c:\n10",
		},
		{
			desc:     "Single-Insert, Double-Remove, Not adjacent",
			patchStr: "v1:\n0:+1:a,\n1:-1:b,\n3:-1:c:\n10",
			expected: "v1:\n0:+1:a,\n1:-1:b,\n3:-1:c:\n10",
		},
		{
			desc:     "Interleaved Insert-Delete-Insert, Adjacent",
			patchStr: "v1:\n0:+1:a,\n1:-1:b,\n2:+1:c:\n10",
			expected: "v1:\n0:+1:a,\n1:-1:b,\n2:+1:c:\n10",
		},
		{
			desc:     "Interleaved Delete-Insert-Delete, Adjacent",
			patchStr: "v1:\n0:-1:a,\n1:+1:b,\n2:-1:c:\n10",
			expected: "v1:\n0:-1:a,\n1:+1:b,\n2:-1:c:\n10",
		},
	}

	for _, test := range tests {
		patch, err := NewPatchFromString(test.patchStr)
		require.Nil(t, err)

		if want, got := test.expected, patch.String(); want != got {
			t.Errorf("TestPatchSimplify[%s]: Expected %s, but got %s. Diffs: %v", test.desc, want, got, pretty.Diff(want, got))
		}
	}
}
