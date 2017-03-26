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
