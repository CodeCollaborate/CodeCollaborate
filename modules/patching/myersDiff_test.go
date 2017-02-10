package patching

import (
	"testing"

	"github.com/kr/pretty"
)

func TestMyersDiff(t *testing.T) {

	tests := []struct {
		desc string
		str1 string
		str2 string
	}{
		{
			desc: "Myers Diff paper example",
			str1: "abcabba",
			str2: "cbabac",
		},
		{
			desc: "Complete removal",
			str1: "abcde",
			str2: "",
		},
		{
			desc: "Prefix removal 1",
			str1: "abcde",
			str2: "bcde",
		},
		{
			desc: "Prefix removal 2",
			str1: "abcde",
			str2: "cde",
		},
		{
			desc: "Suffix removal 1",
			str1: "abcde",
			str2: "abcd",
		},
		{
			desc: "Suffix removal 2",
			str1: "abcde",
			str2: "abc",
		},
		{
			desc: "Suffix removal 3",
			str1: "abcde",
			str2: "a",
		},
		{
			desc: "Mid-string removal 1",
			str1: "abcdefg",
			str2: "abfg",
		},
		{
			desc: "Mid-string removal 2",
			str1: "abcdefg",
			str2: "aefg",
		},
		{
			desc: "Mid-string removal 3",
			str1: "abcdefg",
			str2: "abdfg",
		},
		{
			desc: "Complete insertion",
			str1: "",
			str2: "abcdefg",
		},
		{
			desc: "Complete insertion",
			str1: "",
			str2: "abcdefg",
		},
		{
			desc: "Prefix insertion",
			str1: "bcdefg",
			str2: "abcdefg",
		},
		{
			desc: "Prefix insertion 2",
			str1: "efg",
			str2: "abcdefg",
		},
		{
			desc: "Suffix insertion 1",
			str1: "abcdef",
			str2: "abcdefg",
		},
		{
			desc: "Suffix insertion 2",
			str1: "ab",
			str2: "abcdefg",
		},
		{
			desc: "Mid-string insertion 1",
			str1: "ag",
			str2: "abcdefg",
		},
		{
			desc: "Mid-string insertion 2",
			str1: "abcefg",
			str2: "abcdefg",
		},
		{
			desc: "Mid-string insertion 3",
			str1: "acdfg",
			str2: "abcdefg",
		},
	}

	for _, test := range tests {

		diffs, err := myersDiff(test.str1, test.str2)
		if err != nil {
			t.Errorf("TestMyersDiff[%s]: Unexpected error: %q", test.desc, err)
			continue
		}

		patch := NewPatch(0, diffs, 0)
		res, err := PatchText(test.str1, []*Patch{patch})
		if err != nil {
			t.Errorf("TestMyersDiff[%s]: Unexpected error: %q", test.desc, err)
			continue
		}

		if want, got := test.str2, res; want != got {
			t.Errorf("TestMyersDiff[%s]: Expected %s, but got %s. Diffs: %v", test.desc, want, got, pretty.Diff(want, got))
		}
	}
}
