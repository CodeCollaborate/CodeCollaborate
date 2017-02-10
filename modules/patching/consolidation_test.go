package patching

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			expectedDiff:  "1:+0:",
			expectedIndex: 0,
		},
		{
			desc:          "Insertion noOp",
			patch:         getPatchesOrDie(t, "v0:\n1:+1:A:\n5")[0],
			currIndex:     0,
			wasNoOp:       false,
			expectedDiff:  "1:+0:",
			expectedIndex: 0,
		},
		{
			desc:          "Removal noOp",
			patch:         getPatchesOrDie(t, "v0:\n1:-1:A:\n5")[0],
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
			expectedDiff:  "0:+0:",
			expectedIndex: 0,
		},
		{
			desc:          "Ending noOp with no space after",
			patch:         getPatchesOrDie(t, "v0:\n0:+1:A:\n0")[0],
			currIndex:     0,
			wasNoOp:       false,
			expectedDiff:  "0:+0:",
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
			patches:  getPatchesOrDie(t, "v0:\n0:+7:testing:\n0", "v1:\n1:+2:AB,\n4:+4:CDEF,\n7:+3:GHI:\n7"),
		},
		{
			desc:     "Simple deletion-addition test 1",
			baseText: "testing",
			patches:  getPatchesOrDie(t, "v0:\n2:-2:st,\n5:-1:n:\n7", "v1:\n1:+2:AB,\n4:+4:CDEF:\n4"),
		},
		{
			desc:     "Simple deletion-addition test 2",
			baseText: "testing",
			patches:  getPatchesOrDie(t, "v0:\n2:-2:st,\n5:-1:n:\n7", "v1:\n0:-2:te,\n4:+4:CDEF:\n4"),
		},
		{
			desc:     "Mixed deletion-addition test 1",
			baseText: "testing",
			patches:  getPatchesOrDie(t, "v0:\n1:+2:AB,\n5:-1:n:\n7", "v1:\n0:-2:tA,\n4:+4:CDEF:\n8"),
		},
		{
			desc:     "Mixed deletion-addition test 2",
			baseText: "testing",
			patches:  getPatchesOrDie(t, "v0:\n1:-3:est,\n5:+2:AB:\n6", "v1:\n0:-2:ti,\n3:+4:CDEF:\n9"),
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

func TestConsolidatePatchLong(t *testing.T) {

	tests := []overallConsolidationTest{
		{
			desc:     "Simple Add-Only test",
			baseText: "",
			patches: getPatchesOrDie(t,
				"v1:\n0:+1:y:\n0", "v2:\n1:+1:e:\n1", "v3:\n2:+1:p:\n2", "v4:\n3:+1:%0A:\n3", "v5:\n4:+1:%0A:\n4",
				"v6:\n3:+1:%0A:\n5", "v7:\n6:+1:O:\n6", "v8:\n4:+1:%0A:\n7", "v9:\n5:+1:%0A:\n8", "v9:\n6:+1:%0A:\n9",
				"v10:\n9:-1:O:\n10", "v11:\n6:+1:I:\n9", "v12:\n7:+1:%27:\n10", "v13:\n11:+1:O:\n11", "v14:\n12:+1:k:\n12",
				"v15:\n8:+1:m:\n13", "v16:\n9:+1:+:\n14", "v16:\n10:+1:n:\n15", "v17:\n11:+1:o:\n16", "v17:\n12:+1:t:\n17",
				"v17:\n13:+1:+:\n18", "v17:\n14:+1:s:\n19", "v17:\n15:+1:u:\n20", "v18:\n16:+1:p:\n21", "v19:\n17:+1:e:\n22",
				"v20:\n18:+1:r:\n23", "v21:\n19:+1:+:\n24", "v21:\n20:+1:g:\n25", "v21:\n21:+1:o:\n26", "v21:\n22:+1:o:\n27",
				"v21:\n23:+1:d:\n28", "v21:\n24:+1:+:\n29", "v22:\n25:+1:a:\n30", "v22:\n26:+1:t:\n31", "v23:\n27:+1:+:\n32",
				"v24:\n33:+1:.:\n33", "v25:\n28:+1:t:\n34", "v26:\n29:+1:y:\n35", "v27:\n30:+1:p:\n36", "v28:\n31:+1:i:\n37",
				"v29:\n32:+1:n:\n38", "v30:\n33:+1:g:\n39", "v31:\n34:+1:+:\n40", "v32:\n35:+1:h:\n41", "v33:\n36:+1:o:\n42",
				"v34:\n37:+1:n:\n43", "v35:\n38:+1:e:\n44", "v36:\n39:+1:s:\n45", "v37:\n46:+1:+:\n46", "v37:\n47:+1:S:\n47",
				"v37:\n48:+1:o:\n48", "v37:\n49:+1:+:\n49", "v38:\n40:+1:t:\n50", "v38:\n41:+1:l:\n51", "v38:\n42:+1:y:\n52",
				"v39:\n53:+1:w:\n53", "v39:\n54:+1:e:\n54", "v39:\n54:-1:e:\n55", "v39:\n53:-1:w:\n54", "v39:\n53:+1:w:\n53",
				"v39:\n54:+1:h:\n54", "v39:\n55:+1:e:\n55", "v39:\n56:+1:n:\n56", "v39:\n57:+1:+:\n57", "v39:\n58:+1:w:\n58",
				"v39:\n59:+1:e:\n59", "v39:\n60:+1:+:\n60", "v39:\n61:+1:a:\n61", "v39:\n62:+1:r:\n62", "v39:\n63:+1:e:\n63",
				"v39:\n64:+1:+:\n64", "v39:\n65:+1:t:\n65", "v39:\n66:+1:y:\n66", "v39:\n67:+1:p:\n67", "v39:\n68:+1:i:\n68",
				"v39:\n69:+1:n:\n69", "v39:\n70:+1:g:\n70", "v39:\n71:+1:+:\n71", "v39:\n72:+1:t:\n72", "v39:\n73:+1:o:\n73",
				"v39:\n74:+1:g:\n74", "v39:\n75:+1:e:\n75", "v40:\n76:+1:t:\n76", "v40:\n77:+1:h:\n77", "v41:\n78:+1:e:\n78",
				"v42:\n79:+1:r:\n79", "v43:\n80:+1:%2C:\n80", "v43:\n81:+1:+:\n81", "v43:\n82:+1:i:\n82", "v44:\n83:+1:t:\n83",
				"v45:\n84:+1:+:\n84", "v46:\n85:+1:u:\n85", "v47:\n86:+1:p:\n86", "v48:\n87:+1:d:\n87", "v49:\n43:+1:.:\n88",
				"v50:\n44:+1:.:\n89", "v51:\n45:+1:.:\n90", "v52:\n91:+1:a:\n91", "v53:\n92:+1:t:\n92", "v53:\n93:+1:e:\n93",
				"v53:\n94:+1:s:\n94", "v53:\n95:+1:%2C:\n95", "v53:\n96:+1:+:\n96", "v53:\n97:+1:a:\n97", "v53:\n98:+1:n:\n98",
				"v53:\n99:+1:d:\n99", "v53:\n100:+1:+:\n100", "v53:\n101:+1:d:\n101", "v53:\n102:+1:o:\n102", "v53:\n103:+1:e:\n103",
				"v53:\n104:+1:s:\n104", "v53:\n105:+1:+:\n105", "v53:\n106:+1:i:\n106", "v53:\n107:+1:t:\n107", "v53:\n108:+1:+:\n108",
				"v53:\n109:+1:i:\n109", "v53:\n110:+1:n:\n110", "v53:\n111:+1:+:\n111", "v53:\n112:+1:t:\n112", "v53:\n113:+1:h:\n113",
				"v53:\n114:+1:e:\n114", "v53:\n115:+1:+:\n115", "v53:\n116:+1:c:\n116", "v53:\n117:+1:o:\n117", "v54:\n118:+1:r:\n118",
				"v54:\n119:+1:r:\n119", "v54:\n120:+1:e:\n120", "v55:\n121:+1:c:\n121", "v56:\n122:+1:+:\n122", "v56:\n123:+1:t:\n123",
				"v57:\n46:+1:%0A:\n124", "v58:\n47:+1:%0A:\n125", "v58:\n48:+1:s:\n126", "v58:\n49:+1:o:\n127", "v58:\n50:+1:m:\n128",
				"v58:\n51:+1:e:\n129", "v58:\n52:+1:t:\n130", "v58:\n53:+1:i:\n131", "v58:\n54:+1:m:\n132", "v58:\n55:+1:e:\n133",
				"v58:\n56:+1:s:\n134", "v58:\n57:+1:+:\n135", "v59:\n135:-1:t:\n136", "v60:\n134:-1:+:\n135", "v60:\n134:+1:t:\n134",
				"v60:\n135:+1:+:\n135", "v60:\n136:+1:l:\n136", "v61:\n137:+1:o:\n137", "v61:\n138:+1:c:\n138", "v61:\n139:+1:a:\n139",
				"v61:\n140:+1:t:\n140", "v61:\n141:+1:i:\n141", "v61:\n142:+1:o:\n142", "v61:\n143:+1:n:\n143", "v61:\n144:+1:.:\n144",
				"v62:\n58:+1:i:\n145", "v63:\n146:+1:%0A:\n146", "v64:\n147:+1:%0A:\n147", "v65:\n148:+1:S:\n148", "v66:\n149:+1:o:\n149",
				"v67:\n59:+1:t:\n150", "v67:\n60:+1:+:\n151", "v68:\n152:+1:+:\n152", "v69:\n61:+1:w:\n153", "v69:\n62:+1:i:\n154",
				"v69:\n63:+1:l:\n155", "v69:\n64:+1:l:\n156", "v69:\n65:+1:+:\n157", "v69:\n66:+1:k:\n158", "v69:\n67:+1:i:\n159",
				"v69:\n68:+1:n:\n160", "v69:\n69:+1:d:\n161", "v69:\n70:+1:a:\n162", "v70:\n71:+1:+:\n163", "v70:\n72:+1:w:\n164",
				"v71:\n72:-1:w:\n165", "v72:\n72:+1:t:\n164", "v72:\n73:+1:a:\n165", "v72:\n74:+1:k:\n166", "v72:\n75:+1:e:\n167",
				"v72:\n76:+1:+:\n168", "v72:\n77:+1:a:\n169", "v72:\n78:+1:+:\n170", "v72:\n79:+1:l:\n171", "v72:\n80:+1:i:\n172",
				"v72:\n81:+1:t:\n173", "v72:\n82:+1:t:\n174", "v72:\n83:+1:l:\n175", "v72:\n84:+1:e:\n176", "v72:\n85:+1:+:\n177",
				"v72:\n86:+1:b:\n178", "v72:\n87:+1:i:\n179", "v72:\n88:+1:t:\n180", "v73:\n89:+1:+:\n181", "v73:\n90:+1:t:\n182",
				"v73:\n91:+1:o:\n183", "v73:\n92:+1:+:\n184", "v73:\n93:+1:c:\n185", "v73:\n94:+1:a:\n186", "v74:\n95:+1:t:\n187",
				"v75:\n96:+1:c:\n188", "v76:\n97:+1:h:\n189", "v77:\n189:-1:+:\n190", "v78:\n188:-1:o:\n189", "v79:\n187:-1:S:\n188",
				"v80:\n187:+1:I:\n187", "v81:\n188:+1:+:\n188", "v82:\n189:+1:s:\n189", "v83:\n190:+1:u:\n190", "v84:\n191:+1:s:\n191",
				"v85:\n192:+1:p:\n192", "v85:\n193:+1:e:\n193", "v85:\n194:+1:c:\n194", "v86:\n195:+1:t:\n195", "v87:\n196:+1:+:\n196",
				"v88:\n197:+1:t:\n197", "v89:\n198:+1:h:\n198", "v89:\n199:+1:a:\n199", "v89:\n200:+1:t:\n200", "v89:\n201:+1:+:\n201",
				"v89:\n202:+1:i:\n202", "v89:\n203:+1:s:\n203", "v89:\n204:+1:+:\n204", "v89:\n205:+1:p:\n205", "v90:\n206:+1:i:\n206",
				"v90:\n207:+1:n:\n207", "v90:\n208:+1:g:\n208", "v90:\n209:+1:.:\n209", "v91:\n98:+1:+:\n210", "v91:\n99:+1:u:\n211",
				"v91:\n100:+1:p:\n212", "v91:\n101:+1:+:\n213", "v92:\n101:-1:+:\n214", "v93:\n213:+1:%0A:\n213", "v94:\n214:+1:%0A:\n214",
				"v94:\n215:+1:%0A:\n215",
				"v95:\n0:-216:yep%0A%0A%0AI%27m+not+super+good+at+typing+honestly...%0A%0Asometimes+it+will+kinda+take+a+little+bit+to+catch+up%0A%0A%0AOk.+So+when+we+are+typing+together%2C+it+updates%2C+and+does+it+in+the+correct+location.%0A%0AI+suspect+that+is+ping.%0A%0A%0A:\n216",
				"v96:\n0:+1:t:\n0", "v97:\n1:+1:e:\n1", "v98:\n2:+1:s:\n2", "v99:\n3:+1:t:\n3", "v100:\n4:+1:i:\n4",
				"v101:\n5:+1:n:\n5", "v102:\n6:+1:g:\n6", "v103:\n6:-1:g:\n7", "v104:\n5:-1:n:\n6", "v105:\n4:-1:i:\n5",
				"v106:\n3:-1:t:\n4", "v107:\n2:-1:s:\n3", "v108:\n1:-1:e:\n2", "v109:\n0:-1:t:\n1", "v110:\n0:+1:l:\n0",
				"v111:\n1:+1:t:\n1", "v112:\n2:+1:j:\n2", "v113:\n3:+1:s:\n3", "v114:\n4:+1:l:\n4", "v115:\n5:+1:e:\n5",
				"v116:\n6:+1:k:\n6", "v117:\n7:+1:l:\n7", "v118:\n8:+1:a:\n8", "v119:\n9:+1:s:\n9", "v120:\n10:+1:d:\n10",
				"v121:\n11:+1:f:\n11", "v122:\n0:-12:ltjsleklasdf,\n0:+1:H:\n12", "v123:\n1:+1:e:\n1", "v124:\n2:+1:l:\n2",
				"v125:\n3:+1:l:\n3", "v126:\n4:+1:o:\n4", "v127:\n5:+1:+:\n5", "v128:\n6:+1:m:\n6", "v129:\n7:+1:y:\n7",
				"v130:\n8:+1:+:\n8", "v131:\n9:+1:n:\n9", "v132:\n10:+1:a:\n10", "v133:\n11:+1:m:\n11", "v134:\n12:+1:e:\n12",
				"v135:\n13:+1:+:\n13", "v136:\n14:+1:i:\n14", "v137:\n15:+1:s:\n15", "v138:\n16:+1:+:\n16", "v139:\n17:+1:b:\n17",
				"v140:\n18:+1:e:\n18", "v141:\n19:+1:n:\n19", "v142:\n20:+1:.:\n20", "v143:\n21:+1:+:\n21", "v144:\n22:+1:%0A:\n22",
				"v145:\n23:+1:%0A:\n23", "v146:\n24:+1:I:\n24", "v147:\n25:+1:t:\n25", "v148:\n26:+1:+:\n26", "v149:\n27:+1:w:\n27",
				"v150:\n27:-1:w:\n28", "v151:\n27:+1:h:\n27", "v152:\n28:+1:a:\n28", "v153:\n29:+1:n:\n29", "v154:\n30:+1:d:\n30",
				"v155:\n31:+1:l:\n31", "v156:\n32:+1:e:\n32", "v157:\n33:+1:s:\n33", "v158:\n34:+1:c:\n34", "v159:\n35:+1:o:\n35",
				"v160:\n36:+1:n:\n36", "v161:\n37:+1:c:\n37", "v162:\n38:+1:u:\n38", "v163:\n39:+1:r:\n39", "v164:\n40:+1:r:\n40",
				"v165:\n41:+1:e:\n41", "v166:\n42:+1:n:\n42", "v167:\n43:+1:t:\n43", "v168:\n44:+1:+:\n44", "v169:\n45:+1:e:\n45",
				"v170:\n46:+1:d:\n46", "v171:\n47:+1:i:\n47", "v172:\n48:+1:t:\n48", "v173:\n49:+1:s:\n49", "v174:\n50:+1:+:\n50",
				"v175:\n51:+1:o:\n51", "v176:\n52:+1:n:\n52", "v177:\n53:+1:+:\n53", "v178:\n0:+1:d:\n54", "v179:\n1:+1:i:\n55",
				"v180:\n2:+1:f:\n56", "v181:\n3:+1:f:\n57", "v182:\n4:+1:e:\n58", "v183:\n5:+1:r:\n59", "v184:\n6:+1:e:\n60",
				"v185:\n7:+1:n:\n61", "v186:\n8:+1:t:\n62", "v187:\n9:+1:+:\n63", "v188:\n10:+1:e:\n64", "v189:\n11:+1:d:\n65",
				"v190:\n12:+1:i:\n66", "v191:\n13:+1:t:\n67", "v192:\n14:+1:o:\n68", "v193:\n15:+1:r:\n69", "v194:\n16:+1:s:\n70",
				"v195:\n17:+1:+:\n71", "v196:\n72:+1:p:\n72", "v197:\n73:+1:e:\n73", "v198:\n74:+1:r:\n74", "v199:\n75:+1:f:\n75",
				"v200:\n76:+1:e:\n76", "v201:\n77:+1:c:\n77", "v202:\n78:+1:t:\n78", "v203:\n79:+1:+:\n79", "v204:\n80:+1:l:\n80",
				"v205:\n81:+1:y:\n81", "v206:\n82:+1:+:\n82", "v207:\n82:-1:+:\n83", "v208:\n81:-1:y:\n82", "v209:\n80:-1:l:\n81",
				"v210:\n79:-1:+:\n80", "v211:\n79:+1:l:\n79", "v212:\n80:+1:y:\n80", "v213:\n81:+1:+:\n81", "v214:\n82:+1:w:\n82",
				"v215:\n83:+1:e:\n83", "v216:\n84:+1:l:\n84", "v217:\n85:+1:l:\n85", "v218:\n86:+1:.:\n86",
				"v219:\n0:-87:different+editors+Hello+my+name+is+ben.+%0A%0AIt+handlesconcurrent+edits+on+perfectly+well.,\n0:+1:t:\n87",
				"v220:\n1:+1:e:\n1", "v221:\n2:+1:s:\n2", "v222:\n3:+1:t:\n3", "v223:\n4:+1:i:\n4", "v224:\n5:+1:n:\n5", "v225:\n6:+1:g:\n6",
				"v226:\n7:+1:+:\n7", "v227:\n8:+1:h:\n8", "v228:\n9:+1:e:\n9", "v229:\n10:+1:l:\n10", "v230:\n11:+1:l:\n11", ""+"v231:\n12:+1:o:\n12",
				"v232:\n13:+1:e:\n13", "v233:\n14:+1:s:\n14", "v234:\n15:+1:a:\n15", "v235:\n16:+1:f:\n16",
			),
		},
	}

	for _, test := range tests {
		patchedText, err := PatchText(test.baseText, test.patches)
		require.Nil(t, err)

		consolidatedPatch, err := ConsolidatePatches(test.patches)
		require.Nil(t, err)

		consolidatedPatchedText, err := PatchText(test.baseText, []*Patch{consolidatedPatch})
		require.Equal(t, patchedText, consolidatedPatchedText, "TestConsolidatePatchLong[%s]: Expected %s but got %s", test.desc, patchedText, consolidatedPatchedText)
	}
}
