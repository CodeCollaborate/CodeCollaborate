package patching

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPatch_NewPatch(t *testing.T) {

	// Test addition
	diff := NewDiff(true, 2, "test")
	patch := NewPatch(2, []*Diff{diff})
	require.Equal(t, 2, patch.BaseVersion)
	require.Equal(t, 1, len(patch.Changes))
	require.Equal(t, diff, patch.Changes[0])

	// Test removal
	diff = NewDiff(false, 42, "string")
	patch = NewPatch(10, []*Diff{diff})
	require.Equal(t, 10, patch.BaseVersion)
	require.Equal(t, 1, len(patch.Changes))
	require.Equal(t, diff, patch.Changes[0])
}

func TestPatch_NewPatchFromString(t *testing.T) {

	// Test insertion from string
	patch, err := NewPatchFromString("v4:\n2:+1:a")
	require.Nil(t, err)
	require.Equal(t, 4, patch.BaseVersion)
	require.Equal(t, 1, len(patch.Changes))
	require.Equal(t, "2:+1:a", patch.Changes[0].String())

	// Test insertion from string
	patch, err = NewPatchFromString("v3:\n26:+2:ab,\n81:-3:cde")
	require.Nil(t, err)
	require.Equal(t, 3, patch.BaseVersion)
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
	patch, err := NewPatchFromString("v2:\n0:+5:test%0A")
	require.Nil(t, err)
	newPatch := patch.ConvertToCRLF("\r\ntest")
	require.Equal(t, "v2:\n0:+6:test%0D%0A", newPatch.String())

	patch, err = NewPatchFromString("v6:\n1:+5:test%0A,\n7:+4:test")
	require.Nil(t, err)
	newPatch = patch.ConvertToCRLF("\r\ntes\r\nt")
	require.Equal(t, "v6:\n2:+6:test%0D%0A,\n9:+4:test", newPatch.String())
}

func TestPatch_ConvertToLF(t *testing.T) {
	patch, err := NewPatchFromString("v2:\n0:+6:test%0D%0A")
	require.Nil(t, err)
	newPatch := patch.ConvertToLF("\ntest")
	require.Equal(t, "v2:\n0:+5:test%0A", newPatch.String())

	patch, err = NewPatchFromString("v6:\n2:+6:test%0D%0A,\n9:+4:test")
	require.Nil(t, err)
	newPatch = patch.ConvertToLF("\ntes\nt")
	require.Equal(t, "v6:\n1:+5:test%0A,\n7:+4:test", newPatch.String())
}

func TestPatch_GetUndo(t *testing.T) {
	patch, err := NewPatchFromString("v2:\n0:+4:test")
	require.Nil(t, err)
	newPatch := patch.Undo()
	require.Equal(t, "v2:\n0:-4:test", newPatch.String())

	patch, err = NewPatchFromString("v6:\n2:+6:test%0D%0A,\n9:+4:test")
	require.Nil(t, err)
	newPatch = patch.Undo()
	require.Equal(t, "v6:\n9:-4:test,\n2:-6:test%0D%0A", newPatch.String())
}

func TestPatch_Transform(t *testing.T) {

	// Test set 1
	patch1String := "v1:\n0:-1:a"
	patch2String := "v0:\n3:-8:deletion,\n3:+6:insert"

	patch1, err := NewPatchFromString(patch1String)
	require.Nil(t, err)
	patch2, err := NewPatchFromString(patch2String)
	require.Nil(t, err)

	expectedString := "v0:\n2:-8:deletion,\n2:+6:insert"
	result := patch2.Transform([]*Patch{patch1})
	require.Equal(t, expectedString, result.String())

	// Test set 2
	patch1String = "v1:\n0:-1:a"
	patch2String = "v2:\n0:-1:b"
	patch3String := "v0:\n3:-8:deletion,\n3:+6:insert"

	patch1, err = NewPatchFromString(patch1String)
	require.Nil(t, err)
	patch2, err = NewPatchFromString(patch2String)
	require.Nil(t, err)
	patch3, err := NewPatchFromString(patch3String)
	require.Nil(t, err)

	expectedString = "v0:\n1:-8:deletion,\n1:+6:insert"
	result = patch3.Transform([]*Patch{patch1, patch2})
	require.Equal(t, expectedString, result.String())
}
