package patching

import (
	"testing"

	"github.com/kr/pretty"
)

func TestFilePatcher_ApplyPatch(t *testing.T) {

	tests := []struct {
		desc     string
		patches  []*Patch
		text     string
		expected string
		error    string
	}{
		{
			desc:     "Single Patch, Single insertion",
			patches:  GetPatchesOrDie(t, "v0:1:\n2:+1:a:\n10"),
			text:     "test",
			expected: "teast",
		},
		{
			desc:     "Single Patch, Single deletion",
			patches:  GetPatchesOrDie(t, "v0:1:\n2:-1:s:\n10"),
			text:     "test",
			expected: "tet",
		},
		{
			desc:    "Single Patch, Single deletion, Incorrect base text",
			patches: GetPatchesOrDie(t, "v0:1:\n2:-1:s:\n10"),
			text:    "aaaa",
			error:   "PatchText: Deleted text [a] does not match changes in diff: [s]",
		},
		{
			desc:     "Single Patch, Double insertion",
			patches:  GetPatchesOrDie(t, "v0:1:\n2:+1:m,\n3:+2:ab:\n10"),
			text:     "test",
			expected: "temsabt",
		},
		{
			desc:     "Single Patch, Double deletion",
			patches:  GetPatchesOrDie(t, "v0:1:\n0:-1:t,\n2:-2:st:\n10"),
			text:     "test",
			expected: "e",
		},
		{
			desc:     "Single Patch, Insert+Delete",
			patches:  GetPatchesOrDie(t, "v0:1:\n1:+1:z,\n2:-2:st:\n10"),
			text:     "test",
			expected: "tze",
		},
		{
			desc:     "Single Patch, Delete+Insert",
			patches:  GetPatchesOrDie(t, "v0:1:\n1:-2:es,\n4:+2:lm:\n10"),
			text:     "test",
			expected: "ttlm",
		},
		{
			desc:     "Double Patch, Single Insertions, 1 first",
			patches:  GetPatchesOrDie(t, "v0:1:\n1:+1:a:\n10", "v1:2:\n1:+1:b:\n10"),
			text:     "test",
			expected: "tbaest",
		},
		{
			desc:     "Double Patch, Single Insertions, 2 first",
			patches:  GetPatchesOrDie(t, "v0:1:\n1:+1:a:\n10", "v1:2:\n2:+1:b:\n10"),
			text:     "test",
			expected: "tabest",
		},
		{
			desc:     "Double Patch, Single Deletions, 1 first",
			patches:  GetPatchesOrDie(t, "v0:1:\n1:-1:e:\n10", "v1:2:\n1:-1:s:\n10"),
			text:     "test",
			expected: "tt",
		},
		{
			desc:     "Double Patch, Insert-Deletes",
			patches:  GetPatchesOrDie(t, "v0:1:\n1:+1:z,\n2:-2:st:\n10", "v1:2:\n0:+2:aa,\n2:-1:e:\n10"),
			text:     "test",
			expected: "aatz",
		},
		{
			desc:     "Double Patch, Delete-Inserts",
			patches:  GetPatchesOrDie(t, "v0:1:\n1:-2:es,\n3:+2:lm:\n10", "v1:2:\n0:-2:tl,\n3:+2:kk:\n10"),
			text:     "test",
			expected: "mkkt",
		},
		{
			desc:     "Double Patch, Insert-Delete, Delete-Insert",
			patches:  GetPatchesOrDie(t, "v0:1:\n1:+1:z,\n2:-2:st:\n10", "v1:2:\n0:-2:tz,\n3:+2:ab:\n10"),
			text:     "test",
			expected: "eab",
		},
		{
			desc:     "Double Patch, Delete-Insert, Insert-Delete",
			patches:  GetPatchesOrDie(t, "v0:1:\n1:-2:es,\n3:+2:lm:\n10", "v1:2:\n0:-2:tl,\n3:+2:kk:\n10"),
			text:     "test",
			expected: "mkkt",
		},
	}

	for _, test := range tests {
		result, err := PatchText(test.text, test.patches)

		if test.error != "" {
			if err == nil {
				t.Errorf("TestApplyPatch[%s]: Expected error: %q", test.desc, test.error)
				continue
			}
			if want, got := test.error, err.Error(); want != got {
				t.Error(pretty.Sprintf("TestApplyPatch[%s]: Expected %q, got %q. Diffs: %v", test.desc, want, got, pretty.Diff(want, got)))
				continue
			}
		} else if err != nil {
			t.Errorf("TestApplyPatch[%s]: Unexpected error: %q", test.desc, err)
			continue
		} else {
			if want, got := test.expected, result; want != got {
				t.Errorf("TestApplyPatch[%s]: Expected %s, but got %s. Diffs: %v", test.desc, want, got, pretty.Diff(want, got))
			}
		}
	}
}
