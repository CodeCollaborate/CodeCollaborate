package patching

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// GetPatchOrDie parses a patch from the given string, or throws a fatal error if it fails
func GetPatchOrDie(t *testing.T, patchStr string) *Patch {
	patch, err := NewPatchFromString(patchStr)
	require.Nil(t, err)

	return patch
}

// GetPatchesOrDie parses a set of patches from the given array, or throws a fatal error if it fails
func GetPatchesOrDie(t *testing.T, patchStrs ...string) []*Patch {
	patches := []*Patch{}

	for _, str := range patchStrs {
		patch, err := NewPatchFromString(str)
		if err != nil {
			t.Fatalf("Failed to build patch from string %s", str)
		}
		patches = append(patches, patch)
	}

	return patches
}
