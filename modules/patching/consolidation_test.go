package patching

import (
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

//type diffConsolidationTest struct {
//	desc      string
//	diffA     *Diff
//	diffB     *Diff
//	expectedA Diffs
//	expectedB *Diff
//	offset    int
//	error     string
//}
//
//func runTests(t *testing.T, testName string, tests []diffConsolidationTest) {
//	for _, test := range tests {
//		diffsA, diffB, offset, err := ConsolidateDiffs(test.diffA, test.diffB)
//
//		if test.error != "" {
//			if err == nil {
//				t.Errorf("%s[%s]: Expected error: %q", testName, test.desc, test.error)
//				continue
//			}
//			if want, got := test.error, err.Error(); !strings.Contains(got, want) {
//				t.Error(pretty.Sprintf("%s[%s]: Expected %q, got %q. Diffs: %v", testName, test.desc, want, got, pretty.Diff(want, got)))
//				continue
//			}
//		} else if err != nil {
//			t.Errorf("%s[%s]: Unexpected error: %q", testName, test.desc, err)
//			continue
//		} else {
//			require.Equal(t, test.offset, offset)
//			require.Equal(t, test.expectedA.Len(), diffsA.Len(),
//				"%s[%s]: DiffsA unequal length; Expected %d, but got %d", testName, test.desc, test.expectedA.Len(), diffsA.Len())
//			for i := 0; i < test.expectedA.Len(); i++ {
//				require.Equal(t, test.expectedA[i].String(), diffsA[i].String(),
//					"%s[%s]: Expected %s, but got %s", testName, test.desc, test.expectedA[i].String(), diffsA[i].String())
//			}
//			if test.expectedB != nil {
//				if diffB == nil {
//					t.Errorf("%s[%s]: DiffB was nil, expected %s", testName, test.desc, test.expectedB.String())
//				} else {
//					require.Equal(t, test.expectedB.String(), diffB.String(),
//						"%s[%s]: Expected %s, but got %s", testName, test.desc, test.expectedB.String(), diffB.String())
//				}
//			} else {
//				if diffB != nil {
//					t.Errorf("%s[%s]: DiffB was not nil, got %s", testName, test.desc, diffB.String())
//				}
//			}
//		}
//	}
//}
//
//func TestConsolidate_Diff1A(t *testing.T) {
//	tests := []diffConsolidationTest{
//		{
//			desc:      "Non-overlapping",
//			diffA:     getDiffOrDie(t, "0:+2:ab"),
//			diffB:     getDiffOrDie(t, "3:+2:cd"),
//			expectedA: Diffs{getDiffOrDie(t, "0:+2:ab")},
//			expectedB: getDiffOrDie(t, "3:+2:cd"),
//			offset:    0,
//		},
//		{
//			desc:      "Adjacent",
//			diffA:     getDiffOrDie(t, "0:+2:ab"),
//			diffB:     getDiffOrDie(t, "2:+2:cd"),
//			expectedA: Diffs{getDiffOrDie(t, "0:+4:abcd")},
//			expectedB: nil,
//			offset:    2,
//		},
//		{
//			desc:      "Overlapping 1",
//			diffA:     getDiffOrDie(t, "0:+2:ad"),
//			diffB:     getDiffOrDie(t, "1:+2:bc"),
//			expectedA: Diffs{getDiffOrDie(t, "0:+4:abcd")},
//			expectedB: nil,
//			offset:    2,
//		},
//		{
//			desc:      "Overlapping 2",
//			diffA:     getDiffOrDie(t, "0:+2:ag"),
//			diffB:     getDiffOrDie(t, "1:+5:bcdef"),
//			expectedA: Diffs{getDiffOrDie(t, "0:+7:abcdefg")},
//			expectedB: nil,
//			offset:    5,
//		},
//	}
//
//	runTests(t, "TestConsolidate_Diff1A", tests)
//}
//
//func TestConsolidate_Diff1B(t *testing.T) {
//	tests := []diffConsolidationTest{
//		{
//			desc:      "Non-overlapping",
//			diffA:     getDiffOrDie(t, "0:+2:ab"),
//			diffB:     getDiffOrDie(t, "3:-2:cd"),
//			expectedA: Diffs{getDiffOrDie(t, "0:+2:ab")},
//			expectedB: getDiffOrDie(t, "3:-2:cd"),
//			offset:    0,
//		},
//		{
//			desc:      "Adjacent",
//			diffA:     getDiffOrDie(t, "0:+2:ab"),
//			diffB:     getDiffOrDie(t, "2:-2:cd"),
//			expectedA: Diffs{getDiffOrDie(t, "0:+2:ab")},
//			expectedB: getDiffOrDie(t, "2:-2:cd"),
//			offset:    0,
//		},
//		{
//			desc:      "Overlapping, DiffA longer",
//			diffA:     getDiffOrDie(t, "0:+4:abcd"),
//			diffB:     getDiffOrDie(t, "1:-2:bc"),
//			expectedA: Diffs{getDiffOrDie(t, "0:+2:ad")},
//			expectedB: nil,
//			offset:    -2,
//		},
//		{
//			desc:      "Overlapping, same ending point",
//			diffA:     getDiffOrDie(t, "0:+4:abcd"),
//			diffB:     getDiffOrDie(t, "1:-3:bcd"),
//			expectedA: Diffs{getDiffOrDie(t, "0:+1:a")},
//			expectedB: nil,
//			offset:    -3,
//		},
//		{
//			desc:      "Overlapping, DiffB longer",
//			diffA:     getDiffOrDie(t, "0:+4:abcd"),
//			diffB:     getDiffOrDie(t, "1:-5:bcdef"),
//			expectedA: Diffs{getDiffOrDie(t, "0:+1:a")},
//			expectedB: getDiffOrDie(t, "1:-2:ef"),
//			offset:    -3,
//		},
//		{
//			desc:  "Overlapping, DiffA longer, non-matching overlap",
//			diffA: getDiffOrDie(t, "0:+4:abcd"),
//			diffB: getDiffOrDie(t, "1:-2:xy"),
//			error: "Overlapping strings did not match",
//		},
//		{
//			desc:  "Overlapping, same ending point, non-matching overlap",
//			diffA: getDiffOrDie(t, "0:+4:abcd"),
//			diffB: getDiffOrDie(t, "1:-3:xyz"),
//			error: "Overlapping strings did not match",
//		},
//		{
//			desc:  "Overlapping, DiffB longer, non-matching overlap",
//			diffA: getDiffOrDie(t, "0:+4:abcd"),
//			diffB: getDiffOrDie(t, "1:-4:wxyz"),
//			error: "Overlapping strings did not match",
//		},
//	}
//
//	runTests(t, "TestConsolidate_Diff1B", tests)
//}
//
//func TestConsolidate_Diff1C(t *testing.T) {
//	tests := []diffConsolidationTest{
//		{
//			desc:      "Non-overlapping",
//			diffA:     getDiffOrDie(t, "0:-2:ab"),
//			diffB:     getDiffOrDie(t, "3:+2:cd"),
//			expectedA: Diffs{getDiffOrDie(t, "0:-2:ab")},
//			expectedB: getDiffOrDie(t, "5:+2:cd"),
//			offset:    0,
//		},
//		{
//			desc:      "Adjacent",
//			diffA:     getDiffOrDie(t, "0:-2:ab"),
//			diffB:     getDiffOrDie(t, "2:+2:cd"),
//			expectedA: Diffs{getDiffOrDie(t, "0:-2:ab")},
//			expectedB: getDiffOrDie(t, "4:+2:cd"),
//			offset:    0,
//		},
//		{
//			desc:      "Overlapping 1",
//			diffA:     getDiffOrDie(t, "0:-2:ad"),
//			diffB:     getDiffOrDie(t, "1:+2:bc"),
//			expectedA: Diffs{getDiffOrDie(t, "0:-2:ad")},
//			expectedB: getDiffOrDie(t, "3:+2:bc"),
//			offset:    0,
//		},
//		{
//			desc:      "Overlapping 2",
//			diffA:     getDiffOrDie(t, "0:-2:ag"),
//			diffB:     getDiffOrDie(t, "1:+5:bcdef"),
//			expectedA: Diffs{getDiffOrDie(t, "0:-2:ag")},
//			expectedB: getDiffOrDie(t, "3:+5:bcdef"),
//			offset:    0,
//		},
//		{
//			desc:      "Overlapping 2",
//			diffA:     getDiffOrDie(t, "0:-8:abcdefgh"),
//			diffB:     getDiffOrDie(t, "1:+2:ij"),
//			expectedA: Diffs{getDiffOrDie(t, "0:-8:abcdefgh")},
//			expectedB: getDiffOrDie(t, "9:+2:ij"),
//			offset:    0,
//		},
//	}
//
//	runTests(t, "TestConsolidate_Diff1C", tests)
//}
//
//func TestConsolidate_Diff1D(t *testing.T) {
//	tests := []diffConsolidationTest{
//		{
//			desc:      "Non-overlapping",
//			diffA:     getDiffOrDie(t, "0:-2:ab"),
//			diffB:     getDiffOrDie(t, "3:-2:cd"),
//			expectedA: Diffs{getDiffOrDie(t, "0:-2:ab")},
//			expectedB: getDiffOrDie(t, "5:-2:cd"),
//			offset:    0,
//		},
//		{
//			desc:      "Adjacent",
//			diffA:     getDiffOrDie(t, "0:-2:ab"),
//			diffB:     getDiffOrDie(t, "2:-2:cd"),
//			expectedA: Diffs{getDiffOrDie(t, "0:-2:ab")},
//			expectedB: getDiffOrDie(t, "4:-2:cd"),
//			offset:    0,
//		},
//		{
//			desc:      "Overlapping 1",
//			diffA:     getDiffOrDie(t, "0:-2:ad"),
//			diffB:     getDiffOrDie(t, "1:-2:bc"),
//			expectedA: Diffs{getDiffOrDie(t, "0:-2:ad")},
//			expectedB: getDiffOrDie(t, "3:-2:bc"),
//			offset:    0,
//		},
//		{
//			desc:      "Overlapping 2",
//			diffA:     getDiffOrDie(t, "0:-2:ag"),
//			diffB:     getDiffOrDie(t, "1:-5:bcdef"),
//			expectedA: Diffs{getDiffOrDie(t, "0:-2:ag")},
//			expectedB: getDiffOrDie(t, "3:-5:bcdef"),
//			offset:    0,
//		},
//		{
//			desc:      "Overlapping 2",
//			diffA:     getDiffOrDie(t, "0:-8:abcdefgh"),
//			diffB:     getDiffOrDie(t, "1:-2:ij"),
//			expectedA: Diffs{getDiffOrDie(t, "0:-8:abcdefgh")},
//			expectedB: getDiffOrDie(t, "9:-2:ij"),
//			offset:    0,
//		},
//	}
//
//	runTests(t, "TestConsolidate_Diff1D", tests)
//}
//
//func TestConsolidate_Diff2A(t *testing.T) {
//	tests := []diffConsolidationTest{
//		{
//			desc:      "Overlapping 1",
//			diffA:     getDiffOrDie(t, "2:+2:ab"),
//			diffB:     getDiffOrDie(t, "2:+2:cd"),
//			expectedA: Diffs{getDiffOrDie(t, "2:+4:cdab")},
//			expectedB: nil,
//			offset:    2,
//		},
//		{
//			desc:      "Overlapping 2",
//			diffA:     getDiffOrDie(t, "1:+2:ad"),
//			diffB:     getDiffOrDie(t, "1:+2:bc"),
//			expectedA: Diffs{getDiffOrDie(t, "1:+4:bcad")},
//			expectedB: nil,
//			offset:    2,
//		},
//		{
//			desc:      "Overlapping 3",
//			diffA:     getDiffOrDie(t, "0:+2:ag"),
//			diffB:     getDiffOrDie(t, "0:+5:bcdef"),
//			expectedA: Diffs{getDiffOrDie(t, "0:+7:bcdefag")},
//			expectedB: nil,
//			offset:    5,
//		},
//		{
//			desc:      "Overlapping 4",
//			diffA:     getDiffOrDie(t, "5:+8:abcdefgh"),
//			diffB:     getDiffOrDie(t, "5:+2:ij"),
//			expectedA: Diffs{getDiffOrDie(t, "5:+10:ijabcdefgh")},
//			expectedB: nil,
//			offset:    2,
//		},
//	}
//
//	runTests(t, "TestConsolidate_Diff2A", tests)
//}
//
//func TestConsolidate_Diff2B(t *testing.T) {
//	tests := []diffConsolidationTest{
//		{
//			desc:      "DiffA longer",
//			diffA:     getDiffOrDie(t, "2:+4:abcd"),
//			diffB:     getDiffOrDie(t, "2:-2:ab"),
//			expectedA: Diffs{getDiffOrDie(t, "2:+2:cd")},
//			expectedB: nil,
//			offset:    -2,
//		},
//		{
//			desc:      "Same length",
//			diffA:     getDiffOrDie(t, "1:+2:ad"),
//			diffB:     getDiffOrDie(t, "1:-2:ad"),
//			expectedA: Diffs{},
//			expectedB: nil,
//			offset:    -2,
//		},
//		{
//			desc:      "DiffB longer",
//			diffA:     getDiffOrDie(t, "0:+2:ag"),
//			diffB:     getDiffOrDie(t, "0:-5:aghed"),
//			expectedA: Diffs{},
//			expectedB: getDiffOrDie(t, "0:-3:hed"),
//			offset:    -2,
//		},
//		{
//			desc:  "DiffA longer, non-matching overlap",
//			diffA: getDiffOrDie(t, "2:+4:abcd"),
//			diffB: getDiffOrDie(t, "2:-2:xy"),
//			error: "Overlapping strings did not match",
//		},
//		{
//			desc:  "Same length",
//			diffA: getDiffOrDie(t, "1:+2:ad"),
//			diffB: getDiffOrDie(t, "1:-2:xy"),
//			error: "Overlapping strings did not match",
//		},
//		{
//			desc:  "DiffB longer",
//			diffA: getDiffOrDie(t, "0:+2:ag"),
//			diffB: getDiffOrDie(t, "0:-5:xyhed"),
//			error: "Overlapping strings did not match",
//		},
//	}
//
//	runTests(t, "TestConsolidate_Diff2B", tests)
//}
//
//func TestConsolidate_Diff2C(t *testing.T) {
//	tests := []diffConsolidationTest{
//		{
//			desc:      "No shared string, equal lengths",
//			diffA:     getDiffOrDie(t, "2:-4:abcd"),
//			diffB:     getDiffOrDie(t, "2:+4:efgh"),
//			expectedA: Diffs{getDiffOrDie(t, "2:-4:abcd"), getDiffOrDie(t, "6:+4:efgh")},
//			expectedB: nil,
//			offset:    4,
//		},
//		{
//			desc:      "No shared string, diffB longer",
//			diffA:     getDiffOrDie(t, "2:-3:foo"),
//			diffB:     getDiffOrDie(t, "2:+6:barbaz"),
//			expectedA: Diffs{getDiffOrDie(t, "2:-3:foo"), getDiffOrDie(t, "5:+6:barbaz")},
//			expectedB: nil,
//			offset:    6,
//		},
//		{
//			desc:      "No shared string, diffA longer",
//			diffA:     getDiffOrDie(t, "2:-6:barbaz"),
//			diffB:     getDiffOrDie(t, "2:+3:foo"),
//			expectedA: Diffs{getDiffOrDie(t, "2:-6:barbaz"), getDiffOrDie(t, "8:+3:foo")},
//			expectedB: nil,
//			offset:    3,
//		},
//		{
//			desc:      "Shared prefix, equal lengths",
//			diffA:     getDiffOrDie(t, "2:-7:fooabcd"),
//			diffB:     getDiffOrDie(t, "2:+7:foowxyz"),
//			expectedA: Diffs{getDiffOrDie(t, "5:-4:abcd"), getDiffOrDie(t, "9:+4:wxyz")},
//			expectedB: nil,
//			offset:    7,
//		},
//		{
//			desc:      "Shared prefix, diffA longer",
//			diffA:     getDiffOrDie(t, "2:-7:fooabcd"),
//			diffB:     getDiffOrDie(t, "2:+5:foowx"),
//			expectedA: Diffs{getDiffOrDie(t, "5:-4:abcd"), getDiffOrDie(t, "9:+2:wx")},
//			expectedB: nil,
//			offset:    5,
//		},
//		{
//			desc:      "Shared prefix, diffB longer",
//			diffA:     getDiffOrDie(t, "2:-5:fooab"),
//			diffB:     getDiffOrDie(t, "2:+7:foowxyz"),
//			expectedA: Diffs{getDiffOrDie(t, "5:-2:ab"), getDiffOrDie(t, "7:+4:wxyz")},
//			expectedB: nil,
//			offset:    7,
//		},
//		{
//			desc:      "Shared suffix, equal lengths",
//			diffA:     getDiffOrDie(t, "2:-7:abcdfoo"),
//			diffB:     getDiffOrDie(t, "2:+7:wxyzfoo"),
//			expectedA: Diffs{getDiffOrDie(t, "2:-4:abcd"), getDiffOrDie(t, "6:+4:wxyz")},
//			expectedB: nil,
//			offset:    7,
//		},
//		{
//			desc:      "Shared suffix, diffA longer",
//			diffA:     getDiffOrDie(t, "2:-7:abcdfoo"),
//			diffB:     getDiffOrDie(t, "2:+5:wxfoo"),
//			expectedA: Diffs{getDiffOrDie(t, "2:-4:abcd"), getDiffOrDie(t, "6:+2:wx")},
//			expectedB: nil,
//			offset:    5,
//		},
//		{
//			desc:      "Shared suffix, diffB longer",
//			diffA:     getDiffOrDie(t, "2:-5:abfoo"),
//			diffB:     getDiffOrDie(t, "2:+7:wxyzfoo"),
//			expectedA: Diffs{getDiffOrDie(t, "2:-2:ab"), getDiffOrDie(t, "4:+4:wxyz")},
//			expectedB: nil,
//			offset:    7,
//		},
//		{
//			desc:      "Mid-string change, differing lengths",
//			diffA:     getDiffOrDie(t, "2:-19:+thequickbrownfox++"),
//			diffB:     getDiffOrDie(t, "2:+10:+lazyfox++"),
//			expectedA: Diffs{getDiffOrDie(t, "3:-13:thequickbrown"), getDiffOrDie(t, "16:+4:lazy")},
//			expectedB: nil,
//			offset:    10,
//		},
//	}
//
//	runTests(t, "TestConsolidate_Diff2C", tests)
//}
//
//func TestConsolidate_Diff2D(t *testing.T) {
//	tests := []diffConsolidationTest{
//		{
//			desc:      "Equal lengths 1",
//			diffA:     getDiffOrDie(t, "0:-2:ab"),
//			diffB:     getDiffOrDie(t, "0:-2:cd"),
//			expectedA: Diffs{getDiffOrDie(t, "0:-4:abcd")},
//			expectedB: nil,
//			offset:    -2,
//		},
//		{
//			desc:      "Equal lengths 2",
//			diffA:     getDiffOrDie(t, "2:-2:ab"),
//			diffB:     getDiffOrDie(t, "2:-2:cd"),
//			expectedA: Diffs{getDiffOrDie(t, "2:-4:abcd")},
//			expectedB: nil,
//			offset:    -2,
//		},
//		{
//			desc:      "DiffB longer",
//			diffA:     getDiffOrDie(t, "5:-2:ag"),
//			diffB:     getDiffOrDie(t, "5:-5:bcdef"),
//			expectedA: Diffs{getDiffOrDie(t, "5:-7:agbcdef")},
//			expectedB: nil,
//			offset:    0 - 5,
//		},
//		{
//			desc:      "DiffA longer",
//			diffA:     getDiffOrDie(t, "21:-8:abcdefgh"),
//			diffB:     getDiffOrDie(t, "21:-2:ij"),
//			expectedA: Diffs{getDiffOrDie(t, "21:-10:abcdefghij")},
//			expectedB: nil,
//			offset:    -2,
//		},
//	}
//
//	runTests(t, "TestConsolidate_Diff2D", tests)
//}
//
//func TestConsolidate_Diff3A(t *testing.T) {
//	tests := []diffConsolidationTest{
//		{
//			desc:      "Non-overlapping",
//			diffA:     getDiffOrDie(t, "5:+2:ab"),
//			diffB:     getDiffOrDie(t, "0:+2:cd"),
//			expectedA: Diffs{getDiffOrDie(t, "5:+2:ab")},
//			expectedB: getDiffOrDie(t, "0:+2:cd"),
//			offset:    0,
//		},
//		{
//			desc:      "Overlapping",
//			diffA:     getDiffOrDie(t, "1:+2:ab"),
//			diffB:     getDiffOrDie(t, "0:+2:cd"),
//			expectedA: Diffs{getDiffOrDie(t, "1:+2:ab")},
//			expectedB: getDiffOrDie(t, "0:+2:cd"),
//			offset:    0,
//		},
//	}
//
//	runTests(t, "TestConsolidate_Diff3A", tests)
//}
//
//func TestConsolidate_Diff3B(t *testing.T) {
//	tests := []diffConsolidationTest{
//		{
//			desc:      "Non-overlapping",
//			diffA:     getDiffOrDie(t, "5:+2:ab"),
//			diffB:     getDiffOrDie(t, "0:-2:cd"),
//			expectedA: Diffs{getDiffOrDie(t, "5:+2:ab")},
//			expectedB: getDiffOrDie(t, "0:-2:cd"),
//			offset:    0,
//		},
//		{
//			desc:      "Overlapping",
//			diffA:     getDiffOrDie(t, "1:+2:ab"),
//			diffB:     getDiffOrDie(t, "0:-2:ca"),
//			expectedA: Diffs{getDiffOrDie(t, "1:+1:b")},
//			expectedB: getDiffOrDie(t, "0:-1:c"),
//			offset:    -1,
//		},
//		{
//			desc:  "Overlapping, non-matching overlap",
//			diffA: getDiffOrDie(t, "1:+2:ab"),
//			diffB: getDiffOrDie(t, "0:-2:cd"),
//			error: "Overlapping strings did not match",
//		},
//		{
//			desc:      "Overlapping, EndIndexA < EndIndexB",
//			diffA:     getDiffOrDie(t, "1:+2:ab"),
//			diffB:     getDiffOrDie(t, "0:-4:cabd"),
//			expectedA: Diffs{},
//			expectedB: getDiffOrDie(t, "0:-2:cd"),
//			offset:    -2,
//		},
//		{
//			desc:      "Overlapping, EndIndexA == EndIndexB",
//			diffA:     getDiffOrDie(t, "1:+2:ab"),
//			diffB:     getDiffOrDie(t, "0:-3:cab"),
//			expectedA: Diffs{},
//			expectedB: getDiffOrDie(t, "0:-1:c"),
//			offset:    -2,
//		},
//		{
//			desc:      "Overlapping, EndIndexA > EndIndexB",
//			diffA:     getDiffOrDie(t, "1:+2:ab"),
//			diffB:     getDiffOrDie(t, "0:-2:ca"),
//			expectedA: Diffs{getDiffOrDie(t, "1:+1:b")},
//			expectedB: getDiffOrDie(t, "0:-1:c"),
//			offset:    -1,
//		},
//	}
//
//	runTests(t, "TestConsolidate_Diff3B", tests)
//}
//
//func TestConsolidate_Diff3C(t *testing.T) {
//	tests := []diffConsolidationTest{
//		{
//			desc:      "Non-overlapping",
//			diffA:     getDiffOrDie(t, "5:-2:ab"),
//			diffB:     getDiffOrDie(t, "0:+2:cd"),
//			expectedA: Diffs{getDiffOrDie(t, "5:-2:ab")},
//			expectedB: getDiffOrDie(t, "0:+2:cd"),
//			offset:    0,
//		},
//		{
//			desc:      "Overlapping",
//			diffA:     getDiffOrDie(t, "1:-2:ab"),
//			diffB:     getDiffOrDie(t, "0:+2:cd"),
//			expectedA: Diffs{getDiffOrDie(t, "1:-2:ab")},
//			expectedB: getDiffOrDie(t, "0:+2:cd"),
//			offset:    0,
//		},
//	}
//
//	runTests(t, "TestConsolidate_Diff3C", tests)
//}
//
//func TestConsolidate_Diff3D(t *testing.T) {
//	tests := []diffConsolidationTest{
//		{
//			desc:      "Non-overlapping",
//			diffA:     getDiffOrDie(t, "3:-2:cd"),
//			diffB:     getDiffOrDie(t, "0:-2:ab"),
//			expectedA: Diffs{getDiffOrDie(t, "3:-2:cd")},
//			expectedB: getDiffOrDie(t, "0:-2:ab"),
//			offset:    0,
//		},
//		{
//			desc:      "Adjacent",
//			diffA:     getDiffOrDie(t, "2:-2:cd"),
//			diffB:     getDiffOrDie(t, "0:-2:ab"),
//			expectedA: Diffs{getDiffOrDie(t, "0:-4:abcd")},
//			expectedB: nil,
//			offset:    -2,
//		},
//		{
//			desc:      "Overlapping 1",
//			diffA:     getDiffOrDie(t, "1:-2:bc"),
//			diffB:     getDiffOrDie(t, "0:-2:ad"),
//			expectedA: Diffs{getDiffOrDie(t, "0:-3:abc")},
//			expectedB: getDiffOrDie(t, "3:-1:d"),
//			offset:    -1,
//		},
//		{
//			desc:      "Overlapping 2",
//			diffA:     getDiffOrDie(t, "1:-5:bcdef"),
//			diffB:     getDiffOrDie(t, "0:-2:ag"),
//			expectedA: Diffs{getDiffOrDie(t, "0:-6:abcdef")},
//			expectedB: getDiffOrDie(t, "6:-1:g"),
//			offset:    -1,
//		},
//	}
//
//	runTests(t, "TestConsolidate_Diff3D", tests)
//}
//
type overallConsolidationTest struct {
	desc       string
	baseText   string
	diffsA     Diffs
	docLengthA int
	diffsB     Diffs
	docLengthB int
	error      error
}

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
			patch:         getPatchesOrDie(t, "v0:\n1:+1:A:\n5")[0],
			currIndex:     -1,
			wasNoOp:       false,
			expectedDiff:  "0:+0:",
			expectedIndex: -1,
		},
		{
			desc:          "First patch after noOp",
			patch:         getPatchesOrDie(t, "v0:\n1:+1:A:\n5")[0],
			currIndex:     -1,
			wasNoOp:       true,
			expectedDiff:  "1:+1:A",
			expectedIndex: 0,
		},
		{
			desc:          "Ending noOp",
			patch:         getPatchesOrDie(t, "v0:\n1:+1:A:\n5")[0],
			currIndex:     0,
			wasNoOp:       false,
			expectedDiff:  "2:+0:",
			expectedIndex: 0,
		},
		{
			desc:          "Beyond ending noOp",
			patch:         getPatchesOrDie(t, "v0:\n1:+1:A:\n5")[0],
			currIndex:     0,
			wasNoOp:       true,
			expectedDiff:  "nil",
			expectedIndex: -1,
		},
		{
			desc:          "Starting noOp with first diff at 0",
			patch:         getPatchesOrDie(t, "v0:\n0:+1:A:\n5")[0],
			currIndex:     -1,
			wasNoOp:       false,
			expectedDiff:  "0:+0:",
			expectedIndex: -1,
		},
		{
			desc:          "Ending noOp after noOp with first diff at 0",
			patch:         getPatchesOrDie(t, "v0:\n0:+1:A:\n5")[0],
			currIndex:     0,
			wasNoOp:       false,
			expectedDiff:  "1:+0:",
			expectedIndex: 0,
		},
		{
			desc:          "Ending noOp with no space after",
			patch:         getPatchesOrDie(t, "v0:\n0:+1:A:\n0")[0],
			currIndex:     0,
			wasNoOp:       false,
			expectedDiff:  "1:+0:",
			expectedIndex: 0,
		},
		{
			desc:          "Beyond ending noOp",
			patch:         getPatchesOrDie(t, "v0:\n0:+1:A:\n5")[0],
			currIndex:     0,
			wasNoOp:       true,
			expectedDiff:  "nil",
			expectedIndex: -1,
		},
		{
			desc:          "Adjacent patch insert-insert 1",
			patch:         getPatchesOrDie(t, "v0:\n0:+1:A,\n1:+1:B:\n5")[0],
			currIndex:     0,
			wasNoOp:       false,
			expectedDiff:  "0:+0:",
			expectedIndex: 0,
		},
		{
			desc:          "Adjacent patch insert-insert 2",
			patch:         getPatchesOrDie(t, "v0:\n0:+1:A,\n1:+1:B:\n5")[0],
			currIndex:     0,
			wasNoOp:       true,
			expectedDiff:  "1:+1:B",
			expectedIndex: 1,
		},
		{
			desc:          "Adjacent patch insert-remove 1",
			patch:         getPatchesOrDie(t, "v0:\n0:+1:A,\n1:-1:B:\n5")[0],
			currIndex:     0,
			wasNoOp:       false,
			expectedDiff:  "0:+0:",
			expectedIndex: 0,
		},
		{
			desc:          "Adjacent patch insert-remove 2",
			patch:         getPatchesOrDie(t, "v0:\n0:+1:A,\n1:-1:B:\n5")[0],
			currIndex:     0,
			wasNoOp:       true,
			expectedDiff:  "1:-1:B",
			expectedIndex: 1,
		},
		{
			desc:          "Adjacent patch remove-insert",
			patch:         getPatchesOrDie(t, "v0:\n0:-1:A,\n1:+1:B:\n5")[0],
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

func TestOverallConsolidation(t *testing.T) {
	tests := []overallConsolidationTest{
		{
			desc:       "Simple Add-Only test",
			baseText:   "",
			diffsA:     Diffs{getDiffOrDie(t, "0:+7:testing")},
			diffsB:     Diffs{getDiffOrDie(t, "1:+2:AB"), getDiffOrDie(t, "4:+4:CDEF"), getDiffOrDie(t, "7:+3:GHI")},
			docLengthA: 0,
			docLengthB: 7,
		},
		{
			desc:       "Simple deletion-addition test 1",
			baseText:   "testing",
			diffsA:     Diffs{getDiffOrDie(t, "2:-2:st"), getDiffOrDie(t, "5:-1:n")},
			diffsB:     Diffs{getDiffOrDie(t, "1:+2:AB"), getDiffOrDie(t, "4:+4:CDEF")},
			docLengthA: 7,
			docLengthB: 4,
		},
		{
			desc:       "Simple deletion-addition test 2",
			baseText:   "testing",
			diffsA:     Diffs{getDiffOrDie(t, "2:-2:st"), getDiffOrDie(t, "5:-1:n")},
			diffsB:     Diffs{getDiffOrDie(t, "0:-2:te"), getDiffOrDie(t, "4:+4:CDEF")},
			docLengthA: 7,
			docLengthB: 4,
		},
		{
			desc:       "Mixed deletion-addition test 1",
			baseText:   "testing",
			diffsA:     Diffs{getDiffOrDie(t, "1:+2:AB"), getDiffOrDie(t, "5:-1:n")},
			diffsB:     Diffs{getDiffOrDie(t, "0:-2:tA"), getDiffOrDie(t, "4:+4:CDEF")},
			docLengthA: 7,
			docLengthB: 8,
		},
		{
			desc:       "Mixed deletion-addition test 2",
			baseText:   "testing",
			diffsA:     Diffs{getDiffOrDie(t, "1:-3:est"), getDiffOrDie(t, "5:+2:AB")},
			diffsB:     Diffs{getDiffOrDie(t, "0:-2:ti"), getDiffOrDie(t, "3:+4:CDEF")},
			docLengthA: 7,
			docLengthB: 6,
		},
	}

	for _, test := range tests {
		patchA := NewPatch(0, test.diffsA, test.docLengthA)
		patchB := NewPatch(0, test.diffsB, test.docLengthB)

		patchedText, err := PatchText(test.baseText, []*Patch{patchA, patchB})
		assert.Nil(t, err)

		consolidatedPatch := ConsolidatePatches(patchA, patchB)
		consolidatedPatchedText, err := PatchText(test.baseText, []*Patch{consolidatedPatch})
		fmt.Printf("TestOverallConsolidation[%s]: expecting: %s, got %s\n", test.desc, patchedText, consolidatedPatchedText)
		assert.Equal(t, patchedText, consolidatedPatchedText)
	}
}
