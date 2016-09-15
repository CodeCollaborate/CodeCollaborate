package patching

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
}

func TestDiff_NewDiffFromStringInvalidFormats(t *testing.T) {

	_, err := NewDiffFromString("test")
	require.NotNil(t, err, "Did not throw an error on invalid format")

	_, err = NewDiffFromString("0:@1:test")
	require.NotNil(t, err, "Did not throw an error on invalid operation type")

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

func TestDiff_GetUndo(t *testing.T) {
	diff, err := NewDiffFromString("0:+4:str1")
	require.Nil(t, err)
	newDiff := diff.Undo()
	require.Equal(t, "0:-4:str1", newDiff.String())

	diff, err = NewDiffFromString("1:-4:str2")
	require.Nil(t, err)
	newDiff = diff.Undo()
	require.Equal(t, "1:+4:str2", newDiff.String())
}

func TestDiff_Transform1A(t *testing.T) {
	diff1 := NewDiff(true, 2, "str1")
	diff2 := NewDiff(true, 4, "str2")
	result := diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "8:+4:str2", result[0].String())

	diff1 = NewDiff(true, 0, "str1")
	diff2 = NewDiff(true, 1, "str2")
	result = diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "5:+4:str2", result[0].String())
}

func TestDiff_Transform1B(t *testing.T) {
	diff1 := NewDiff(true, 2, "str1")
	diff2 := NewDiff(false, 4, "str2")
	result := diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "8:-4:str2", result[0].String())

	diff1 = NewDiff(true, 0, "str1")
	diff2 = NewDiff(false, 1, "str2")
	result = diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "5:-4:str2", result[0].String())
}

func TestDiff_Transform1C(t *testing.T) {
	diff1 := NewDiff(false, 2, "str1")
	diff2 := NewDiff(true, 4, "str2")
	result := diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "2:+4:str2", result[0].String())

	diff1 = NewDiff(false, 2, "longerstr1")
	diff2 = NewDiff(true, 4, "str2")
	result = diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "2:+4:str2", result[0].String())

	diff1 = NewDiff(false, 2, "str1")
	diff2 = NewDiff(true, 6, "str2")
	result = diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "2:+4:str2", result[0].String())

	diff1 = NewDiff(false, 2, "str1")
	diff2 = NewDiff(true, 8, "str2")
	result = diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "4:+4:str2", result[0].String())
}

func TestDiff_Transform1D(t *testing.T) {
	diff1 := NewDiff(false, 2, "str1")
	diff2 := NewDiff(false, 8, "str2")
	result := diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "4:-4:str2", result[0].String())

	diff1 = NewDiff(false, 2, "str1")
	diff2 = NewDiff(false, 6, "str2")
	result = diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "2:-4:str2", result[0].String())

	diff1 = NewDiff(false, 2, "longerstr1")
	diff2 = NewDiff(false, 4, "str2")
	result = diff2.transform([]*Diff{diff1})
	require.Equal(t, 0, len(result))

	diff1 = NewDiff(false, 2, "str1")
	diff2 = NewDiff(false, 4, "str2")
	result = diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "2:-2:r2", result[0].String())
}

func TestDiff_Transform2A(t *testing.T) {
	diff1 := NewDiff(true, 4, "str1")
	diff2 := NewDiff(true, 4, "str2")
	result := diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "8:+4:str2", result[0].String())

	diff1 = NewDiff(true, 0, "longTestString1")
	diff2 = NewDiff(true, 0, "str2")
	result = diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "15:+4:str2", result[0].String())
}

func TestDiff_Transform2B(t *testing.T) {
	diff1 := NewDiff(true, 4, "str1")
	diff2 := NewDiff(false, 4, "str2")
	result := diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "8:-4:str2", result[0].String())

	diff1 = NewDiff(true, 0, "longTestString1")
	diff2 = NewDiff(false, 0, "str2")
	result = diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "15:-4:str2", result[0].String())
}

func TestDiff_Transform2C(t *testing.T) {
	diff1 := NewDiff(false, 4, "str1")
	diff2 := NewDiff(true, 4, "str2")
	result := diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "4:+4:str2", result[0].String())

	diff1 = NewDiff(false, 0, "longTestString1")
	diff2 = NewDiff(true, 0, "str2")
	result = diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "0:+4:str2", result[0].String())
}

func TestDiff_Transform2D(t *testing.T) {
	diff1 := NewDiff(false, 4, "str1")
	diff2 := NewDiff(false, 4, "longTestString2")
	result := diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "4:-11:TestString2", result[0].String())

	diff1 = NewDiff(false, 0, "longTestString1")
	diff2 = NewDiff(false, 0, "str2")
	result = diff2.transform([]*Diff{diff1})
	require.Equal(t, 0, len(result))

	diff1 = NewDiff(false, 4, "str1")
	diff2 = NewDiff(false, 4, "str2")
	result = diff2.transform([]*Diff{diff1})
	require.Equal(t, 0, len(result))
}

func TestDiff_Transform3A(t *testing.T) {
	diff1 := NewDiff(true, 5, "str1")
	diff2 := NewDiff(true, 4, "str2")
	result := diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "4:+4:str2", result[0].String())

	diff1 = NewDiff(true, 4, "str1")
	diff2 = NewDiff(true, 0, "longTestString2")
	result = diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "0:+15:longTestString2", result[0].String())
}

func TestDiff_Transform3B(t *testing.T) {
	diff1 := NewDiff(true, 5, "str1")
	diff2 := NewDiff(false, 4, "longStr2")
	result := diff2.transform([]*Diff{diff1})
	require.Equal(t, 2, len(result))
	require.Equal(t, "4:-1:l", result[0].String())
	require.Equal(t, "8:-7:ongStr2", result[1].String())

	diff1 = NewDiff(true, 8, "str1")
	diff2 = NewDiff(false, 0, "str2")
	result = diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "0:-4:str2", result[0].String())
}

func TestDiff_Transform3C(t *testing.T) {
	diff1 := NewDiff(false, 9, "str1")
	diff2 := NewDiff(true, 4, "str2")
	result := diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "4:+4:str2", result[0].String())

	diff1 = NewDiff(false, 5, "longTestString1")
	diff2 = NewDiff(true, 4, "str2")
	result = diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "4:+4:str2", result[0].String())
}

func TestDiff_Transform3D(t *testing.T) {
	diff1 := NewDiff(false, 6, "str1")
	diff2 := NewDiff(false, 4, "str2")
	result := diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "4:-2:st", result[0].String())

	diff1 = NewDiff(false, 8, "str1")
	diff2 = NewDiff(false, 4, "str2")
	result = diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "4:-4:str2", result[0].String())

	diff1 = NewDiff(false, 10, "str1")
	diff2 = NewDiff(false, 4, "str2")
	result = diff2.transform([]*Diff{diff1})
	require.Equal(t, 1, len(result))
	require.Equal(t, "4:-4:str2", result[0].String())
}
