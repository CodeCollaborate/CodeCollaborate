package patching

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPatch_NewPatch(t *testing.T) {
	patchString := "v1:\n3:-8:deletion,\n2:+6:insert"

	diff1 := NewDiff(false, 3, "deletion")
	diff2 := NewDiff(true, 2, "insert")
	patch := NewPatch(1, []*Diff{diff1, diff2})
	require.Equal(t, patchString, patch.String())
}

func TestPatch_NewPatchFromString(t *testing.T) {
	patch, err := NewPatchFromString("v6:\n3:-8:deletion,\n2:+6:insert")
	require.Nil(t, err)
	require.Equal(t, int64(6), patch.BaseVersion)
	require.Equal(t, 2, len(patch.Changes))
	require.Equal(t, "3:-8:deletion", patch.Changes[0].String())
	require.Equal(t, "2:+6:insert", patch.Changes[1].String())

	// Test insertion from string
	patch, err = NewPatchFromString("v4:\n2:+1:a")
	require.Nil(t, err)
	require.Equal(t, int64(4), patch.BaseVersion)
	require.Equal(t, 1, len(patch.Changes))
	require.Equal(t, "2:+1:a", patch.Changes[0].String())

	// Test insertion from string
	patch, err = NewPatchFromString("v3:\n26:+2:ab,\n81:-3:cde")
	require.Nil(t, err)
	require.Equal(t, int64(3), patch.BaseVersion)
	require.Equal(t, 2, len(patch.Changes))
	require.Equal(t, "26:+2:ab", patch.Changes[0].String())
	require.Equal(t, "81:-3:cde", patch.Changes[1].String())
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
	patch, err := NewPatchFromString("v0:\n0:+5:test%0A")
	require.Nil(t, err)
	newPatch := patch.ConvertToCRLF("\r\ntest")
	require.Equal(t, 1, len(newPatch.Changes))
	require.Equal(t, "v0:\n0:+6:test%0D%0A", newPatch.String())

	patch, err = NewPatchFromString("v0:\n1:+5:test%0A")
	require.Nil(t, err)
	newPatch = patch.ConvertToCRLF("\r\ntest")
	require.Equal(t, 1, len(newPatch.Changes))
	require.Equal(t, "v0:\n2:+6:test%0D%0A", newPatch.String())

	patch, err = NewPatchFromString("v0:\n2:+5:test%0A")
	require.Nil(t, err)
	newPatch = patch.ConvertToCRLF("\r\ntest")
	require.Equal(t, 1, len(newPatch.Changes))
	require.Equal(t, "v0:\n3:+6:test%0D%0A", newPatch.String())

	patch, err = NewPatchFromString("v0:\n7:+5:test%0A")
	require.Nil(t, err)
	newPatch = patch.ConvertToCRLF("\r\ntes\r\nt")
	require.Equal(t, 1, len(newPatch.Changes))
	require.Equal(t, "v0:\n9:+6:test%0D%0A", newPatch.String())

	patch, err = NewPatchFromString("v0:\n2:+5:test%0A,\n7:+5:test%0A")
	require.Nil(t, err)
	newPatch = patch.ConvertToCRLF("\r\ntes\r\nt")
	require.Equal(t, 2, len(newPatch.Changes))
	require.Equal(t, "v0:\n3:+6:test%0D%0A,\n9:+6:test%0D%0A", newPatch.String())

	patch, err = NewPatchFromString("v0:\n2:+5:test%0A,\n7:+5:test%0A,\n0:+5:test%0A")
	require.Nil(t, err)
	newPatch = patch.ConvertToCRLF("\r\ntes\r\nt")
	require.Equal(t, 3, len(newPatch.Changes))
	require.Equal(t, "v0:\n3:+6:test%0D%0A,\n9:+6:test%0D%0A,\n0:+6:test%0D%0A", newPatch.String())
}

func TestPatch_ConvertToLF(t *testing.T) {
	patch, err := NewPatchFromString("v0:\n0:+6:test%0D%0A")
	require.Nil(t, err)
	newPatch := patch.ConvertToLF("\r\ntest")
	require.Equal(t, 1, len(newPatch.Changes))
	require.Equal(t, "v0:\n0:+5:test%0A", newPatch.String())

	patch, err = NewPatchFromString("v0:\n2:+6:test%0D%0A")
	require.Nil(t, err)
	newPatch = patch.ConvertToLF("\r\ntest")
	require.Equal(t, 1, len(newPatch.Changes))
	require.Equal(t, "v0:\n1:+5:test%0A", newPatch.String())

	patch, err = NewPatchFromString("v0:\n3:+6:test%0D%0A")
	require.Nil(t, err)
	newPatch = patch.ConvertToLF("\r\ntest")
	require.Equal(t, 1, len(newPatch.Changes))
	require.Equal(t, "v0:\n2:+5:test%0A", newPatch.String())

	patch, err = NewPatchFromString("v0:\n9:+6:test%0D%0A")
	require.Nil(t, err)
	newPatch = patch.ConvertToLF("\r\ntes\r\nt")
	require.Equal(t, 1, len(newPatch.Changes))
	require.Equal(t, "v0:\n7:+5:test%0A", newPatch.String())

	patch, err = NewPatchFromString("v0:\n3:+6:test%0D%0A,\n9:+6:test%0D%0A")
	require.Nil(t, err)
	newPatch = patch.ConvertToLF("\r\ntes\r\nt")
	require.Equal(t, 2, len(newPatch.Changes))
	require.Equal(t, "v0:\n2:+5:test%0A,\n7:+5:test%0A", newPatch.String())

	patch, err = NewPatchFromString("v0:\n3:+6:test%0D%0A,\n9:+6:test%0D%0A,\n0:+6:test%0D%0A")
	require.Nil(t, err)
	newPatch = patch.ConvertToLF("\r\ntes\r\nt")
	require.Equal(t, 3, len(newPatch.Changes))
	require.Equal(t, "v0:\n2:+5:test%0A,\n7:+5:test%0A,\n0:+5:test%0A", newPatch.String())
}

func TestPatch_Undo(t *testing.T) {
	patch, err := NewPatchFromString("v0:\n0:+5:test%0A")
	require.Nil(t, err)
	newPatch := patch.Undo()
	require.Equal(t, 1, len(newPatch.Changes))
	require.Equal(t, "v0:\n0:-5:test%0A", newPatch.String())

	patch, err = NewPatchFromString("v0:\n1:-5:test%0A")
	require.Nil(t, err)
	newPatch = patch.Undo()
	require.Equal(t, 1, len(newPatch.Changes))
	require.Equal(t, "v0:\n1:+5:test%0A", newPatch.String())

	patch, err = NewPatchFromString("v0:\n2:+5:test%0A")
	require.Nil(t, err)
	newPatch = patch.Undo()
	require.Equal(t, 1, len(newPatch.Changes))
	require.Equal(t, "v0:\n2:-5:test%0A", newPatch.String())

	patch, err = NewPatchFromString("v0:\n7:-5:test%0A")
	require.Nil(t, err)
	newPatch = patch.Undo()
	require.Equal(t, 1, len(newPatch.Changes))
	require.Equal(t, "v0:\n7:+5:test%0A", newPatch.String())

	patch, err = NewPatchFromString("v0:\n2:-5:test%0A,\n7:+5:test%0A")
	require.Nil(t, err)
	newPatch = patch.Undo()
	require.Equal(t, 2, len(newPatch.Changes))
	require.Equal(t, "v0:\n7:-5:test%0A,\n2:+5:test%0A", newPatch.String())

	patch, err = NewPatchFromString("v0:\n2:+5:test%0A,\n7:-5:test%0A,\n0:-5:test%0A")
	require.Nil(t, err)
	newPatch = patch.Undo()
	require.Equal(t, 3, len(newPatch.Changes))
	require.Equal(t, "v0:\n0:+5:test%0A,\n7:+5:test%0A,\n2:-5:test%0A", newPatch.String())
}

func TestPatch_Transform(t *testing.T) {
	// Test set 1
	patch1, err := NewPatchFromString("v1:\n0:-1:a")
	require.Nil(t, err)
	patch2, err := NewPatchFromString("v0:\n3:-8:deletion,\n3:+6:insert")
	require.Nil(t, err)
	newPatch := patch2.Transform([]*Patch{patch1})
	require.Equal(t, 2, len(newPatch.Changes))
	require.Equal(t, "v1:\n2:-8:deletion,\n2:+6:insert", newPatch.String())

	// Test set 2
	patch1, err = NewPatchFromString("v1:\n0:-1:a")
	require.Nil(t, err)
	patch2, err = NewPatchFromString("v2:\n0:-1:b")
	require.Nil(t, err)
	patch3, err := NewPatchFromString("v0:\n3:-8:deletion,\n3:+6:insert")
	require.Nil(t, err)
	newPatch = patch3.Transform([]*Patch{patch1, patch2})
	require.Equal(t, 2, len(newPatch.Changes))
	require.Equal(t, "v2:\n1:-8:deletion,\n1:+6:insert", newPatch.String())
}
