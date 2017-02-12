package patching

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func getDiffOrDie(t *testing.T, diffStr string) *Diff {
	diff, err := NewDiffFromString(diffStr)
	if err != nil {
		t.Fatalf("Failed to build diff from string %s", diffStr)
	}

	return diff
}

func TestDiff_NewDiff(t *testing.T) {
	// Test addition
	diff := NewDiff(true, 1, "test")
	require.Equal(t, true, diff.Insertion)
	require.Equal(t, 1, diff.StartIndex)
	require.Equal(t, "test", diff.Changes)

	// Test removal
	diff = NewDiff(false, 10, "string")
	require.Equal(t, false, diff.Insertion)
	require.Equal(t, 10, diff.StartIndex)
	require.Equal(t, "string", diff.Changes)

}

func TestDiff_NewDiffFromString(t *testing.T) {
	// Test insertion from string
	diff, err := NewDiffFromString("10:+4:test")
	require.Nil(t, err)
	require.Equal(t, true, diff.Insertion)
	require.Equal(t, 10, diff.StartIndex)
	require.Equal(t, "test", diff.Changes)

	// Test removal from string
	diff, err = NewDiffFromString("3:-3:del")
	require.Nil(t, err)
	require.Equal(t, false, diff.Insertion)
	require.Equal(t, 3, diff.StartIndex)
	require.Equal(t, "del", diff.Changes)

	// Test emoji from string
	diff, err = NewDiffFromString("3:-1:ω")
	require.Nil(t, err)
	require.Equal(t, false, diff.Insertion)
	require.Equal(t, 3, diff.StartIndex)
	require.Equal(t, 1, diff.Length())
	require.Equal(t, "ω", diff.Changes)
}

func TestDiff_NewDiffFromStringInvalidFormats(t *testing.T) {
	_, err := NewDiffFromString("delete 2")
	require.NotNil(t, err, "Did not throw an error on invalid format")

	_, err = NewDiffFromString("test")
	require.NotNil(t, err, "Did not throw an error on invalid format")

	_, err = NewDiffFromString("0:@1:test")
	require.NotNil(t, err, "Did not throw an error on invalid operation type")

	_, err = NewDiffFromString("3:-1:del")
	require.NotNil(t, err, "Did not throw an error on wrong changes length")

	_, err = NewDiffFromString("0:+1:test")
	require.NotNil(t, err, "Did not throw an error on wrong changes length")

	_, err = NewDiffFromString("a:+4:test")
	require.NotNil(t, err, "Did not throw an error on invalid offset")

	_, err = NewDiffFromString("0:+err:test")
	require.NotNil(t, err, "Did not throw an erorr on invalid length.")
}

func TestDiff_ConvertToCRLF(t *testing.T) {
	diff, err := NewDiffFromString("0:+5:test%0A")
	require.Nil(t, err)
	newDiff := diff.ConvertToCRLF("\r\ntest")
	require.Equal(t, "0:+6:test%0D%0A", newDiff.String())

	diff, err = NewDiffFromString("1:+5:test%0A")
	require.Nil(t, err)
	newDiff = diff.ConvertToCRLF("\r\ntest")
	require.Equal(t, "2:+6:test%0D%0A", newDiff.String())

	diff, err = NewDiffFromString("2:+5:test%0A")
	require.Nil(t, err)
	newDiff = diff.ConvertToCRLF("\r\ntest")
	require.Equal(t, "3:+6:test%0D%0A", newDiff.String())

	diff, err = NewDiffFromString("7:+5:test%0A")
	require.Nil(t, err)
	newDiff = diff.ConvertToCRLF("\r\ntes\r\nt")
	require.Equal(t, "9:+6:test%0D%0A", newDiff.String())
}

func TestDiff_ConvertToLF(t *testing.T) {
	diff, err := NewDiffFromString("0:+6:test%0D%0A")
	require.Nil(t, err)
	newDiff := diff.ConvertToLF("\r\ntest")
	require.Equal(t, "0:+5:test%0A", newDiff.String())

	diff, err = NewDiffFromString("2:+6:test%0D%0A")
	require.Nil(t, err)
	newDiff = diff.ConvertToLF("\r\ntest")
	require.Equal(t, "1:+5:test%0A", newDiff.String())

	diff, err = NewDiffFromString("7:+6:test%0D%0A")
	require.Nil(t, err)
	newDiff = diff.ConvertToLF("\r\ntes\r\nt")
	require.Equal(t, "5:+5:test%0A", newDiff.String())
}

func TestDiff_ConvertBack(t *testing.T) {
	diff, err := NewDiffFromString("53:+3:a%0D%0A")
	require.Nil(t, err)

	LFDiff := diff.ConvertToLF("package testPkg1;\r\n" +
		"\r\n" +
		"public class TestClass1 {\r\n" +
		"\r\n" +
		"}\r\n")
	require.Equal(t, "48:+2:a%0A", LFDiff.String())

	CRLFDiff := LFDiff.ConvertToCRLF("package testPkg1;\r\n" +
		"\r\n" +
		"public class TestClass1 {\r\n" +
		"\r\n" +
		"}\r\n")
	require.Equal(t, diff.String(), CRLFDiff.String())
}

func TestDiff_Undo(t *testing.T) {
	diff, err := NewDiffFromString("0:+4:str1")
	require.Nil(t, err)
	newDiff := diff.Undo()
	require.Equal(t, "0:-4:str1", newDiff.String())
	originalDiff := newDiff.Undo()
	require.Equal(t, diff.String(), originalDiff.String())

	diff, err = NewDiffFromString("1:-4:str2")
	require.Nil(t, err)
	newDiff = diff.Undo()
	require.Equal(t, "1:+4:str2", newDiff.String())
	originalDiff = newDiff.Undo()
	require.Equal(t, diff.String(), originalDiff.String())
}

func TestDiff_Transform1A(t *testing.T) {
	diff1, err := NewDiffFromString("2:+4:str1")
	require.Nil(t, err)
	diff2, err := NewDiffFromString("4:+4:str2")
	require.Nil(t, err)
	result := diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "8:+4:str2", result[0].String())

	diff1, err = NewDiffFromString("0:+4:str1")
	require.Nil(t, err)
	diff2, err = NewDiffFromString("1:+4:str2")
	require.Nil(t, err)
	result = diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "5:+4:str2", result[0].String())
}

func TestDiff_Transform1B(t *testing.T) {
	diff1, err := NewDiffFromString("2:+4:str1")
	require.Nil(t, err)
	diff2, err := NewDiffFromString("4:-4:str2")
	require.Nil(t, err)
	result := diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "8:-4:str2", result[0].String())

	diff1, err = NewDiffFromString("0:+4:str1")
	require.Nil(t, err)
	diff2, err = NewDiffFromString("1:-4:str2")
	require.Nil(t, err)
	result = diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "5:-4:str2", result[0].String())
}

func TestDiff_Transform1C(t *testing.T) {
	// Test case 1: if (IndexA + LenA) > IndexB, shift B down by amoun of A that comes before IndexB
	diff1, err := NewDiffFromString("2:-4:str1")
	require.Nil(t, err)
	diff2, err := NewDiffFromString("4:+4:str2")
	require.Nil(t, err)
	result := diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "2:+4:str2", result[0].String())

	diff1, err = NewDiffFromString("2:-10:longerstr1")
	require.Nil(t, err)
	diff2, err = NewDiffFromString("4:+4:str2")
	require.Nil(t, err)
	result = diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "2:+4:str2", result[0].String())

	// Test else case
	diff1, err = NewDiffFromString("2:-4:str1")
	require.Nil(t, err)
	diff2, err = NewDiffFromString("6:+4:str2")
	require.Nil(t, err)
	result = diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "2:+4:str2", result[0].String())

	diff1, err = NewDiffFromString("2:-4:str1")
	require.Nil(t, err)
	diff2, err = NewDiffFromString("8:+4:str2")
	require.Nil(t, err)
	result = diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "4:+4:str2", result[0].String())
}

func TestDiff_Transform1D(t *testing.T) {
	// Test case 1: if IndexA + LenA < IndexB (No overlap), shift B down by LenA
	diff1, err := NewDiffFromString("2:-4:str1")
	require.Nil(t, err)
	diff2, err := NewDiffFromString("8:-4:str2")
	require.Nil(t, err)
	result := diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "4:-4:str2", result[0].String())

	diff1, err = NewDiffFromString("2:-4:str1")
	require.Nil(t, err)
	diff2, err = NewDiffFromString("6:-4:str2")
	require.Nil(t, err)
	result = diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "2:-4:str2", result[0].String())

	// Test case 2: if IndexA + LenA >= IndexB + LenB, ignore B
	diff1, err = NewDiffFromString("2:-10:longerstr1")
	require.Nil(t, err)
	diff2, err = NewDiffFromString("4:-4:str2")
	require.Nil(t, err)
	result = diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 0, len(result))

	// Test else cases: if overlapping, shorten B by overlap, shift down by LenA - overlap
	diff1, err = NewDiffFromString("2:-4:str1")
	require.Nil(t, err)
	diff2, err = NewDiffFromString("4:-4:str2")
	require.Nil(t, err)
	result = diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "2:-2:r2", result[0].String())
}

func TestDiff_Transform2A(t *testing.T) {
	diff1, err := NewDiffFromString("4:+4:str1")
	require.Nil(t, err)
	diff2, err := NewDiffFromString("4:+4:str2")
	require.Nil(t, err)
	result := diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "8:+4:str2", result[0].String())

	diff1, err = NewDiffFromString("0:+15:longTestString1")
	require.Nil(t, err)
	diff2, err = NewDiffFromString("0:+4:str2")
	require.Nil(t, err)
	result = diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "15:+4:str2", result[0].String())

	diff1, err = NewDiffFromString("4:+4:str1")
	require.Nil(t, err)
	diff2, err = NewDiffFromString("4:+4:str2")
	require.Nil(t, err)
	result = diff2.transform(Diffs{diff1}, false)
	require.Equal(t, 1, len(result))
	require.Equal(t, "4:+4:str2", result[0].String())

	diff1, err = NewDiffFromString("0:+15:longTestString1")
	require.Nil(t, err)
	diff2, err = NewDiffFromString("0:+4:str2")
	require.Nil(t, err)
	result = diff2.transform(Diffs{diff1}, false)
	require.Equal(t, 1, len(result))
	require.Equal(t, "0:+4:str2", result[0].String())
}

func TestDiff_Transform2B(t *testing.T) {
	diff1, err := NewDiffFromString("4:+4:str1")
	require.Nil(t, err)
	diff2, err := NewDiffFromString("4:-4:str2")
	require.Nil(t, err)
	result := diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "8:-4:str2", result[0].String())

	diff1, err = NewDiffFromString("0:+15:longTestString1")
	require.Nil(t, err)
	diff2, err = NewDiffFromString("0:-4:str2")
	require.Nil(t, err)
	result = diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "15:-4:str2", result[0].String())
}

func TestDiff_Transform2C(t *testing.T) {
	diff1, err := NewDiffFromString("4:-4:str1")
	require.Nil(t, err)
	diff2, err := NewDiffFromString("4:+4:str2")
	require.Nil(t, err)
	result := diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "4:+4:str2", result[0].String())

	diff1, err = NewDiffFromString("0:-15:longTestString1")
	require.Nil(t, err)
	diff2, err = NewDiffFromString("0:+4:str2")
	require.Nil(t, err)
	result = diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "0:+4:str2", result[0].String())
}

func TestDiff_Transform2D(t *testing.T) {
	// Test case 1: If LenB > LenA, remove LenA characters from B
	diff1, err := NewDiffFromString("4:-4:str1")
	require.Nil(t, err)
	diff2, err := NewDiffFromString("4:-15:longTestString2")
	require.Nil(t, err)
	result := diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "4:-11:TestString2", result[0].String())

	// Test else case - if LenB <= LenA
	diff1, err = NewDiffFromString("0:-15:longTestString1")
	require.Nil(t, err)
	diff2, err = NewDiffFromString("0:-4:str2")
	require.Nil(t, err)
	result = diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 0, len(result))

	diff1, err = NewDiffFromString("4:-4:str1")
	require.Nil(t, err)
	diff2, err = NewDiffFromString("4:-4:str2")
	require.Nil(t, err)
	result = diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 0, len(result))
}

func TestDiff_Transform3A(t *testing.T) {
	diff1, err := NewDiffFromString("5:+4:str1")
	require.Nil(t, err)
	diff2, err := NewDiffFromString("4:+4:str2")
	require.Nil(t, err)
	result := diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "4:+4:str2", result[0].String())

	diff1, err = NewDiffFromString("4:+4:str1")
	require.Nil(t, err)
	diff2, err = NewDiffFromString("0:+15:longTestString2")
	require.Nil(t, err)
	result = diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "0:+15:longTestString2", result[0].String())
}

func TestDiff_Transform3B(t *testing.T) {
	// Test case 1: If IndexB + LenB > IndexA, split B into two diffs
	diff1, err := NewDiffFromString("5:+4:str1")
	require.Nil(t, err)
	diff2, err := NewDiffFromString("4:-8:longStr2")
	require.Nil(t, err)
	result := diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 2, len(result))
	require.Equal(t, "4:-1:l", result[0].String())
	require.Equal(t, "8:-7:ongStr2", result[1].String())

	// Test else case: no change
	diff1, err = NewDiffFromString("8:+4:str1")
	require.Nil(t, err)
	diff2, err = NewDiffFromString("0:-4:str2")
	require.Nil(t, err)
	result = diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "0:-4:str2", result[0].String())
}

func TestDiff_Transform3C(t *testing.T) {
	diff1, err := NewDiffFromString("9:-4:str1")
	require.Nil(t, err)
	diff2, err := NewDiffFromString("4:+4:str2")
	require.Nil(t, err)
	result := diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "4:+4:str2", result[0].String())

	diff1, err = NewDiffFromString("5:-15:longTestString1")
	require.Nil(t, err)
	diff2, err = NewDiffFromString("4:+4:str2")
	require.Nil(t, err)
	result = diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "4:+4:str2", result[0].String())
}

func TestDiff_Transform3D(t *testing.T) {
	// Test case 1: If IndexB + LenB > IndexA, shorten B by overlap (from end)
	diff1, err := NewDiffFromString("6:-4:str1")
	require.Nil(t, err)
	diff2, err := NewDiffFromString("4:-4:str2")
	require.Nil(t, err)
	result := diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "4:-2:st", result[0].String())

	// Test else case: No change if no overlap
	diff1, err = NewDiffFromString("8:-4:str1")
	require.Nil(t, err)
	diff2, err = NewDiffFromString("4:-4:str2")
	require.Nil(t, err)
	result = diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "4:-4:str2", result[0].String())

	diff1, err = NewDiffFromString("10:-4:str1")
	require.Nil(t, err)
	diff2, err = NewDiffFromString("4:-4:str2")
	require.Nil(t, err)
	result = diff2.transform(Diffs{diff1}, true)
	require.Equal(t, 1, len(result))
	require.Equal(t, "4:-4:str2", result[0].String())
}
