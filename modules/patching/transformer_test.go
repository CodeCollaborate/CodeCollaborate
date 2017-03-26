package patching

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func getPatchOrDie(t *testing.T, patchStr string) *Patch {
	patch, err := NewPatchFromString(patchStr)
	require.Nil(t, err)

	return patch
}

/**
 * Created by wongb on 3/25/17.
 */
func TestPrecedence(t *testing.T) {

	baseText := "abdeletion"
	patchA, err := NewPatchFromString("v0:\n2:-8:deletion,\n2:+6:insert:\n10")
	require.Nil(t, err)
	patchB, err := NewPatchFromString("v0:\n10:+1:a:\n10")
	require.Nil(t, err)

	// Test with A having precedence
	result, err := TransformPatches(patchA, patchB)
	require.Nil(t, err)
	// Validate the the document that each path produces is the same
	resStringA, err := PatchText(baseText, []*Patch{patchA, result.PatchYPrime})
	require.Nil(t, err)
	resStringB, err := PatchText(baseText, []*Patch{patchB, result.PatchXPrime})
	require.Nil(t, err)
	require.Equal(t, resStringA, resStringB, "Resultant documents were different")

	// Test with B having precedence
	result, err = TransformPatches(patchB, patchA)
	require.Nil(t, err)
	// Validate the the document that each path produces is the same
	resStringA, err = PatchText(baseText, []*Patch{patchB, result.PatchYPrime})
	require.Nil(t, err)
	resStringB, err = PatchText(baseText, []*Patch{patchA, result.PatchXPrime})
	require.Nil(t, err)
	require.Equal(t, resStringA, resStringB, "Resultant documents were different")

}

func TestOverlappingDeletes(t *testing.T) {
	baseText := "\n\tHello, my name is Ben. This is a test of whether this works properly.\n" +
		"\t\n" +
		"\tWow eclipse is dumb. It changed my \"Properly\" word to a entire public main method\n" +
		"\ttesting\n" +
		"\tSystem.out.println(\"Hellow this is a test\");\n" +
		"\t\n" +
		"\tif (true == false) {\n" +
		"\t\tdo all the things\n" +
		"\t\t2 = 1\n" +
		"\t\tThis is definitely coherent.  DEFINITELY\n" +
		"\t}"
	patchA, err := NewPatchFromString("v557:\n0:-309:%0A%09Hello%2C+my+name+is+Ben.+This+is+a+test+of+whether+this+works+properly.%0A%09%0A%09Wow+eclipse+is+dumb.+It+changed+my+%22Properly%22+word+to+a+entire+public+main+method%0A%09testing%0A%09System.out.println%28%22Hellow+this+is+a+test%22%29%3B%0A%09%0A%09if+%28true+%3D%3D+false%29+%7B%0A%09%09do+all+the+things%0A%09%092+%3D+1%0A%09%09This+is+definitely+coherent.++DEFINITELY%0A%09%7D,\n0:+1:m:\n309")
	require.Nil(t, err)
	patchB, err := NewPatchFromString("v557:\n308:-1:%7D:\n309")
	require.Nil(t, err)

	// Test with A having precedence
	result, err := TransformPatches(patchA, patchB)
	require.Nil(t, err)
	// Validate the the document that each path produces is the same
	resStringA, err := PatchText(baseText, []*Patch{patchA, result.PatchYPrime})
	require.Nil(t, err)
	resStringB, err := PatchText(baseText, []*Patch{patchB, result.PatchXPrime})
	require.Nil(t, err)
	require.Equal(t, resStringA, resStringB, "Resultant documents were different")

	// Test with B having precedence
	result, err = TransformPatches(patchB, patchA)
	require.Nil(t, err)
	// Validate the the document that each path produces is the same
	resStringA, err = PatchText(baseText, []*Patch{patchB, result.PatchYPrime})
	require.Nil(t, err)
	resStringB, err = PatchText(baseText, []*Patch{patchA, result.PatchXPrime})
	require.Nil(t, err)
	require.Equal(t, resStringA, resStringB, "Resultant documents were different")
}

type transformationTest struct {
	desc                string
	patchA              *Patch
	patchB              *Patch
	baseText            string
	expectedPatchAPrime *Patch
	expectedPatchBPrime *Patch
	canReverse          bool
}

func runTests(t *testing.T, tests []transformationTest) {
	for _, test := range tests {

		result, err := TransformPatches(test.patchA, test.patchB)
		require.Nil(t, err)

		require.Equal(t, test.expectedPatchAPrime.String(), result.PatchXPrime.String(), "TestConsolidator[%s]: Patch A' was incorrect expected [%s], got [%s]",
			test.desc, strings.Replace(test.expectedPatchAPrime.String(), "\n", "\\n", -1), strings.Replace(result.PatchXPrime.String(), "\n", "\\n", -1))
		require.Equal(t, test.expectedPatchBPrime.String(), result.PatchYPrime.String(), "TestConsolidator[%s]: Patch B' was incorrect expected [%s], got [%s]",
			test.desc, strings.Replace(test.expectedPatchBPrime.String(), "\n", "\\n", -1), strings.Replace(result.PatchYPrime.String(), "\n", "\\n", -1))

		// Validate the the document that each path produces is the same
		resStringA, err := PatchText(test.baseText, []*Patch{test.patchA, result.PatchYPrime})
		require.Nil(t, err)
		resStringB, err := PatchText(test.baseText, []*Patch{test.patchB, result.PatchXPrime})
		require.Nil(t, err)
		require.Equal(t, resStringA, resStringB, "TestConsolidator[%s]: Document was different based on patch application order", test.desc)

		// If is reversible (does not require precedence, try running it in reverse
		if test.canReverse {
			result, err := TransformPatches(test.patchB, test.patchA)
			require.Nil(t, err)

			require.Equal(t, test.expectedPatchBPrime.String(), result.PatchXPrime.String(), "TestConsolidator[%s-Reverse]: Patch B' was incorrect expected [%s], got [%s]",
				test.desc, strings.Replace(test.expectedPatchBPrime.String(), "\n", "\\n", -1), strings.Replace(result.PatchXPrime.String(), "\n", "\\n", -1))
			require.Equal(t, test.expectedPatchAPrime.String(), result.PatchYPrime.String(), "TestConsolidator[%s-Reverse]: Patch A' was incorrect expected [%s], got [%s]",
				test.desc, strings.Replace(test.expectedPatchAPrime.String(), "\n", "\\n", -1), strings.Replace(result.PatchYPrime.String(), "\n", "\\n", -1))

			// Validate the the document that each path produces is the same
			resStringA, err := PatchText(test.baseText, []*Patch{test.patchB, result.PatchYPrime})
			require.Nil(t, err)
			resStringB, err := PatchText(test.baseText, []*Patch{test.patchA, result.PatchXPrime})
			require.Nil(t, err)
			require.Equal(t, resStringA, resStringB, "TestConsolidator[%s-Reverse]: Document was different based on patch application order", test.desc)
		}
	}
}

func TestTransform1A(t *testing.T) {
	tests := []transformationTest{
		{
			"Non-Overlapping strings",
			getPatchOrDie(t, "v1:\n0:+4:str1:\n8"),
			getPatchOrDie(t, "v1:\n6:+4:str2:\n8"),
			"baseText",
			getPatchOrDie(t, "v2:\n0:+4:str1:\n12"),
			getPatchOrDie(t, "v2:\n10:+4:str2:\n12"),
			true,
		}, {
			"Overlapping strings",
			getPatchOrDie(t, "v1:\n2:+4:str1:\n8"),
			getPatchOrDie(t, "v1:\n4:+4:str2:\n8"),
			"baseText",
			getPatchOrDie(t, "v2:\n2:+4:str1:\n12"),
			getPatchOrDie(t, "v2:\n8:+4:str2:\n12"),
			true,
		},
	}
	runTests(t, tests)
}

func TestTransform1B(t *testing.T) {
	tests := []transformationTest{
		{
			"Non-Overlapping strings",
			getPatchOrDie(t, "v1:\n0:+2:s1:\n8"),
			getPatchOrDie(t, "v1:\n4:-4:Text:\n8"),
			"baseText",
			getPatchOrDie(t, "v2:\n0:+2:s1:\n4"),
			getPatchOrDie(t, "v2:\n6:-4:Text:\n10"),
			true,
		},
		{
			"Overlapping strings",
			getPatchOrDie(t, "v1:\n2:+4:str1:\n8"),
			getPatchOrDie(t, "v1:\n4:-4:Text:\n8"),
			"baseText",
			getPatchOrDie(t, "v2:\n2:+4:str1:\n4"),
			getPatchOrDie(t, "v2:\n8:-4:Text:\n12"),
			true,
		},
	}

	runTests(t, tests)
}

func TestTransform1C(t *testing.T) {
	tests := []transformationTest{
		{
			"Non-Overlapping strings",
			getPatchOrDie(t, "v1:\n0:-2:ba:\n8"),
			getPatchOrDie(t, "v1:\n4:+4:abcd:\n8"),
			"baseText",
			getPatchOrDie(t, "v2:\n0:-2:ba:\n12"),
			getPatchOrDie(t, "v2:\n2:+4:abcd:\n6"),
			true,
		},
		{
			"Overlapping strings",
			getPatchOrDie(t, "v1:\n2:-4:seTe:\n8"),
			getPatchOrDie(t, "v1:\n4:+4:abcd:\n8"),
			"baseText",
			getPatchOrDie(t, "v2:\n2:-2:se,\n8:-2:Te:\n12"),
			getPatchOrDie(t, "v2:\n2:+4:abcd:\n4"),
			true,
		},
	}

	runTests(t, tests)
}

func TestTransform1D(t *testing.T) {
	tests := []transformationTest{
		{
			"Non-Overlapping strings",
			getPatchOrDie(t, "v1:\n2:-4:str1:\n16"),
			getPatchOrDie(t, "v1:\n8:-4:str2:\n16"),
			"bastr1sestr2Text",
			getPatchOrDie(t, "v2:\n2:-4:str1:\n12"),
			getPatchOrDie(t, "v2:\n4:-4:str2:\n12"),
			true,
		},
		{
			"Non-Overlapping strings, adjacent",
			getPatchOrDie(t, "v1:\n2:-4:str1:\n16"),
			getPatchOrDie(t, "v1:\n6:-4:str2:\n16"),
			"bastr1str2seText",
			getPatchOrDie(t, "v2:\n2:-4:str1:\n12"),
			getPatchOrDie(t, "v2:\n2:-4:str2:\n12"),
			true,
		},
		{
			"Overlapping strings",
			getPatchOrDie(t, "v1:\n2:-4:seTe:\n8"),
			getPatchOrDie(t, "v1:\n4:-4:Text:\n8"),
			"baseText",
			getPatchOrDie(t, "v2:\n2:-2:se:\n4"),
			getPatchOrDie(t, "v2:\n2:-2:xt:\n4"),
			true,
		},
		{
			"Overlapping strings, B subset of A",
			getPatchOrDie(t, "v1:\n2:-6:seText:\n8"),
			getPatchOrDie(t, "v1:\n4:-2:Te:\n8"),
			"baseText",
			getPatchOrDie(t, "v2:\n2:-4:sext:\n6"),
			getPatchOrDie(t, "v2:\n:\n2"),
			true,
		},
	}
	runTests(t, tests)
}

func TestTransform2A(t *testing.T) {
	tests := []transformationTest{
		{
			"Same length strings",
			getPatchOrDie(t, "v1:\n4:+4:str1:\n8"),
			getPatchOrDie(t, "v1:\n4:+4:str2:\n8"),
			"testText",
			getPatchOrDie(t, "v2:\n4:+4:str1:\n12"),
			getPatchOrDie(t, "v2:\n8:+4:str2:\n12"),
			false,
		},
		{
			"A longer",
			getPatchOrDie(t, "v1:\n4:+8:longstr1:\n8"),
			getPatchOrDie(t, "v1:\n4:+4:str2:\n8"),
			"testText",
			getPatchOrDie(t, "v2:\n4:+8:longstr1:\n12"),
			getPatchOrDie(t, "v2:\n12:+4:str2:\n16"),
			false,
		},
		{
			"B longer",
			getPatchOrDie(t, "v1:\n4:+4:str1:\n8"),
			getPatchOrDie(t, "v1:\n4:+8:longstr2:\n8"),
			"testText",
			getPatchOrDie(t, "v2:\n4:+4:str1:\n16"),
			getPatchOrDie(t, "v2:\n8:+8:longstr2:\n12"),
			false,
		},
	}
	runTests(t, tests)
}

func TestTransform2B(t *testing.T) {
	tests := []transformationTest{
		{
			"Same length strings",
			getPatchOrDie(t, "v1:\n2:+4:str1:\n8"),
			getPatchOrDie(t, "v1:\n2:-4:stTe:\n8"),
			"testText",
			getPatchOrDie(t, "v2:\n2:+4:str1:\n4"),
			getPatchOrDie(t, "v2:\n6:-4:stTe:\n12"),
			true,
		},
		{
			"A longer",
			getPatchOrDie(t, "v1:\n2:+8:longstr1:\n8"),
			getPatchOrDie(t, "v1:\n2:-4:stTe:\n8"),
			"testText",
			getPatchOrDie(t, "v2:\n2:+8:longstr1:\n4"),
			getPatchOrDie(t, "v2:\n10:-4:stTe:\n16"),
			true,
		},
		{
			"B longer",
			getPatchOrDie(t, "v1:\n4:+4:str1:\n14"),
			getPatchOrDie(t, "v1:\n4:-6:longer:\n14"),
			"testlongerText",
			getPatchOrDie(t, "v2:\n4:+4:str1:\n8"),
			getPatchOrDie(t, "v2:\n8:-6:longer:\n18"),
			true,
		},
	}
	runTests(t, tests)
}

func TestTransform2C(t *testing.T) {
	tests := []transformationTest{
		{
			"Same length strings",
			getPatchOrDie(t, "v1:\n2:-4:stTe:\n8"),
			getPatchOrDie(t, "v1:\n2:+4:str2:\n8"),
			"testText",
			getPatchOrDie(t, "v2:\n6:-4:stTe:\n12"),
			getPatchOrDie(t, "v2:\n2:+4:str2:\n4"),
			true,
		},
		{
			"A longer",
			getPatchOrDie(t, "v1:\n4:-6:longer:\n14"),
			getPatchOrDie(t, "v1:\n4:+4:str2:\n14"),
			"testlongerText",
			getPatchOrDie(t, "v2:\n8:-6:longer:\n18"),
			getPatchOrDie(t, "v2:\n4:+4:str2:\n8"),
			true,
		},
		{
			"B longer",
			getPatchOrDie(t, "v1:\n2:-4:stTe:\n8"),
			getPatchOrDie(t, "v1:\n2:+8:longstr2:\n8"),
			"testText",
			getPatchOrDie(t, "v2:\n10:-4:stTe:\n16"),
			getPatchOrDie(t, "v2:\n2:+8:longstr2:\n4"),
			true,
		},
	}
	runTests(t, tests)
}

func TestTransform2D(t *testing.T) {
	tests := []transformationTest{
		{
			"Same length strings",
			getPatchOrDie(t, "v1:\n2:-4:stTe:\n8"),
			getPatchOrDie(t, "v1:\n2:-4:stTe:\n8"),
			"testText",
			getPatchOrDie(t, "v2:\n:\n4"),
			getPatchOrDie(t, "v2:\n:\n4"),
			true,
		},
		{
			"A longer",
			getPatchOrDie(t, "v1:\n4:-6:longer:\n14"),
			getPatchOrDie(t, "v1:\n4:-4:long:\n14"),
			"testlongerText",
			getPatchOrDie(t, "v2:\n4:-2:er:\n10"),
			getPatchOrDie(t, "v2:\n:\n8"),
			true,
		},
		{
			"B longer",
			getPatchOrDie(t, "v1:\n2:-4:stlo:\n14"),
			getPatchOrDie(t, "v1:\n2:-6:stlong:\n14"),
			"testlongerText",
			getPatchOrDie(t, "v2:\n:\n8"),
			getPatchOrDie(t, "v2:\n2:-2:ng:\n10"),
			true,
		},
	}
	runTests(t, tests)
}

func TestTransform3A(t *testing.T) {
	tests := []transformationTest{
		{
			"Overlapping strings",
			getPatchOrDie(t, "v1:\n5:+4:str1:\n8"),
			getPatchOrDie(t, "v1:\n4:+4:str2:\n8"),
			"testText",
			getPatchOrDie(t, "v2:\n9:+4:str1:\n12"),
			getPatchOrDie(t, "v2:\n4:+4:str2:\n12"),
			true,
		},
		{
			"Non-overlapping strings",
			getPatchOrDie(t, "v1:\n4:+4:str1:\n8"),
			getPatchOrDie(t, "v1:\n0:+15:longTestString2:\n8"),
			"testText",
			getPatchOrDie(t, "v2:\n19:+4:str1:\n23"),
			getPatchOrDie(t, "v2:\n0:+15:longTestString2:\n12"),
			true,
		},
	}
	runTests(t, tests)
}

func TestTransform3B(t *testing.T) {
	tests := []transformationTest{
		{
			"Overlapping strings",
			getPatchOrDie(t, "v1:\n5:+4:str1:\n8"),
			getPatchOrDie(t, "v1:\n4:-4:Text:\n8"),
			"testText",
			getPatchOrDie(t, "v2:\n4:+4:str1:\n4"),
			getPatchOrDie(t, "v2:\n4:-1:T,\n9:-3:ext:\n12"),
			true,
		},
		{
			"Non-overlapping strings",
			getPatchOrDie(t, "v1:\n5:+4:str1:\n8"),
			getPatchOrDie(t, "v1:\n2:-2:st:\n8"),
			"testText",
			getPatchOrDie(t, "v2:\n3:+4:str1:\n6"),
			getPatchOrDie(t, "v2:\n2:-2:st:\n12"),
			true,
		},
	}
	runTests(t, tests)
}

func TestTransform3C(t *testing.T) {
	tests := []transformationTest{
		{
			"Overlapping strings",
			getPatchOrDie(t, "v1:\n5:-3:ext:\n8"),
			getPatchOrDie(t, "v1:\n4:+4:str2:\n8"),
			"testText",
			getPatchOrDie(t, "v2:\n9:-3:ext:\n12"),
			getPatchOrDie(t, "v2:\n4:+4:str2:\n5"),
			true,
		},
		{
			"Non-overlapping strings",
			getPatchOrDie(t, "v1:\n5:-3:ext:\n8"),
			getPatchOrDie(t, "v1:\n1:+2:s2:\n8"),
			"testText",
			getPatchOrDie(t, "v2:\n7:-3:ext:\n10"),
			getPatchOrDie(t, "v2:\n1:+2:s2:\n5"),
			true,
		},
	}
	runTests(t, tests)
}

func TestTransform3D(t *testing.T) {
	tests := []transformationTest{
		{
			"Overlapping strings, A extends past B",
			getPatchOrDie(t, "v1:\n4:-8:LongerTe:\n14"),
			getPatchOrDie(t, "v1:\n2:-4:stLo:\n14"),
			"testLongerText",
			getPatchOrDie(t, "v2:\n2:-6:ngerTe:\n10"),
			getPatchOrDie(t, "v2:\n2:-2:st:\n6"),
			true,
		},
		{
			"Overlapping strings, A ends at same index as B",
			getPatchOrDie(t, "v1:\n4:-4:Long:\n14"),
			getPatchOrDie(t, "v1:\n2:-6:stLong:\n14"),
			"testLongerText",
			getPatchOrDie(t, "v2:\n:\n8"),
			getPatchOrDie(t, "v2:\n2:-2:st:\n10"),
			true,
		},
		{
			"Overlapping strings, A ends before B",
			getPatchOrDie(t, "v1:\n4:-3:Lon:\n14"),
			getPatchOrDie(t, "v1:\n2:-6:stLong:\n14"),
			"testLongerText",
			getPatchOrDie(t, "v2:\n:\n8"),
			getPatchOrDie(t, "v2:\n2:-3:stg:\n11"),
			true,
		},
		{
			"Non-overlapping strings",
			getPatchOrDie(t, "v1:\n5:-3:ext:\n8"),
			getPatchOrDie(t, "v1:\n1:-2:es:\n8"),
			"testText",
			getPatchOrDie(t, "v2:\n3:-3:ext:\n6"),
			getPatchOrDie(t, "v2:\n1:-2:es:\n5"),
			true,
		},
	}
	runTests(t, tests)
}
