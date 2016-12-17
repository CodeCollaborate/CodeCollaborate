package patching

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// Diffs represents an array of Diff objects, mainly used for sorting
type Diffs []*Diff

func (slice Diffs) Len() int {
	return len(slice)
}

func (slice Diffs) Less(i, j int) bool {
	return slice[i].StartIndex < slice[j].StartIndex
}

func (slice Diffs) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// Diff represents a single change in the document.
type Diff struct {
	Insertion  bool
	StartIndex int
	Changes    string
}

// NewDiff creates a new diff object from the given parameters
func NewDiff(insertion bool, startIndex int, changes string) *Diff {
	return &Diff{
		Insertion:  insertion,
		StartIndex: startIndex,
		Changes:    changes,
	}
}

// NewDiffFromString parses a diff from its string representation.
func NewDiffFromString(str string) (*Diff, error) {
	regex, err := regexp.Compile("\\d+:(\\+|-)\\d+:.+")
	if err != nil {
		return nil, err
	}
	if !regex.MatchString(str) {
		return nil, errors.New("Illegal patch format; should be %d:+%d:%s or %d:-%d:%s")
	}

	parts := strings.Split(str, ":")
	diff := Diff{}

	// Parse startIndex
	diff.StartIndex, err = strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("Invalid offset: %s", parts[0])
	}

	// Switch on diff type
	switch parts[1][0] {
	case '+':
		diff.Insertion = true
	case '-':
		diff.Insertion = false
	default:
		return nil, fmt.Errorf("Invalid operation: %s", parts[1][0:1])
	}

	// Get length of diff
	length, err := strconv.Atoi(parts[1][1:])
	if err != nil {
		return nil, fmt.Errorf("Invalid length: %s", parts[1][1:])
	}

	// Un-URL encode the body
	unescapedChanges, err := url.QueryUnescape(parts[2])
	if err != nil {
		return nil, err
	}

	// Validate changes
	if length != len(unescapedChanges) {
		return nil, fmt.Errorf("Length does not match length of change: %d != %s", length, parts[2])
	}

	diff.Changes = unescapedChanges
	return &diff, nil
}

// Length returns the number of characters changed.
func (diff *Diff) Length() int {
	return len(diff.Changes)
}

// String gives the encoded format of a diff as a String.
func (diff *Diff) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(strconv.Itoa(diff.StartIndex))
	buffer.WriteString(":")
	if diff.Insertion {
		buffer.WriteString("+")
	} else {
		buffer.WriteString("-")
	}
	buffer.WriteString(strconv.Itoa(len(diff.Changes)))
	buffer.WriteString(":")
	buffer.WriteString(url.QueryEscape(diff.Changes))

	return buffer.String()
}

// ConvertToCRLF converts this diff from using LF to CRLF line separators given the base text to patch.
// TODO(wongb): Review this for performance
func (diff *Diff) ConvertToCRLF(base string) *Diff {
	newStartIndex := diff.StartIndex
	newChanges := strings.Replace(diff.Changes, "\n", "\r\n", -1)

	for i := 0; i < newStartIndex && i < len(base)-1; i++ {
		if base[i] == '\r' && base[i+1] == '\n' {
			newStartIndex++
		}
	}

	return NewDiff(diff.Insertion, newStartIndex, newChanges)
}

// ConvertToLF converts this diff from using CRLF to LF line separators given the base (CRLF) text that it was generated from.
// TODO(wongb): Review this for performance
func (diff *Diff) ConvertToLF(base string) *Diff {
	newStartIndex := diff.StartIndex
	newChanges := strings.Replace(diff.Changes, "\r\n", "\n", -1)

	for i := 0; i < diff.StartIndex-1 && i < len(base)-1; i++ {
		if base[i] == '\r' && base[i+1] == '\n' {
			newStartIndex--
		}
	}

	return NewDiff(diff.Insertion, newStartIndex, newChanges)
}

// Undo reverses this diff, producing a diff to undo the changes done by applying the diff.
func (diff *Diff) Undo() *Diff {
	return NewDiff(!diff.Insertion, diff.StartIndex, diff.Changes)
}

// OffsetDiff shifts the start index of this diff by the provided offset
func (diff *Diff) OffsetDiff(offset int) *Diff {
	return NewDiff(diff.Insertion, diff.StartIndex+offset, diff.Changes)
}

func (diff *Diff) subChanges(start, end int) *Diff {
	return NewDiff(diff.Insertion, diff.StartIndex, diff.Changes[start:end])
}

func (diff *Diff) subChangesStartingFrom(start int) *Diff {
	return NewDiff(diff.Insertion, diff.StartIndex, diff.Changes[start:])
}

func (diff *Diff) subChangesEndingAt(end int) *Diff {
	return NewDiff(diff.Insertion, diff.StartIndex, diff.Changes[:end])
}

func (diff *Diff) transform(others Diffs) Diffs {
	intermediateDiffs := Diffs{}
	intermediateDiffs = append(intermediateDiffs, diff)

	for _, other := range others {
		newIntermediateDiffs := Diffs{}
		for _, current := range intermediateDiffs {
			switch {
			// CASE 1: IndexA < IndexB
			case other.StartIndex < current.StartIndex:
				switch {
				// CASE 1a, 1b: Ins - Ins, Ins - Rmv
				case other.Insertion && current.Insertion, other.Insertion && !current.Insertion:
					newIntermediateDiffs = doTransform(transformType2, newIntermediateDiffs, current, other)
				// CASE 1c: Rmv - Ins
				case !other.Insertion && current.Insertion:
					newIntermediateDiffs = doTransform(transformType3, newIntermediateDiffs, current, other)
				// CASE 1d: Rmv - Rmv
				case !other.Insertion && !current.Insertion:
					newIntermediateDiffs = doTransform(transformType4, newIntermediateDiffs, current, other)
				// Fail; should never have gotten to here.
				default:
					panic(fmt.Sprintf("Got to invalid state 1e while transforming [%s] on predessor [%+v], from list [%+v]", current.String(), other, others))
				}
			// CASE 2: IndexA = IndexB
			case other.StartIndex == current.StartIndex:
				switch {
				// CASES 2a, 2b: Ins - Ins, Ins - Rmv
				case other.Insertion && current.Insertion, other.Insertion && !current.Insertion:
					newIntermediateDiffs = doTransform(transformType2, newIntermediateDiffs, current, other)
				// CASE 2c: Rmv - Ins
				// Do nothing
				case !other.Insertion && current.Insertion:
					newIntermediateDiffs = doTransform(transformType1, newIntermediateDiffs, current, other)
				// CASE 2d: Rmv - Ins
				case !other.Insertion && !current.Insertion:
					newIntermediateDiffs = doTransform(transformType5, newIntermediateDiffs, current, other)
				// Fail; should never have gotten to here.
				default:
					panic(fmt.Sprintf("Got to invalid state 2e while transforming [%s] on predessor [%+v], from list [%+v]", current.String(), other, others))
				}
			// CASE 3: IndexA > IndexB
			case other.StartIndex > current.StartIndex:
				switch {
				// CASE 3a, 3c: Ins - Ins, Rmv - Ins
				case other.Insertion && current.Insertion, !other.Insertion && current.Insertion:
					newIntermediateDiffs = doTransform(transformType1, newIntermediateDiffs, current, other)
				// CASE 3b: Ins - Rmv
				case other.Insertion && !current.Insertion:
					newIntermediateDiffs = doTransform(transformType6, newIntermediateDiffs, current, other)
				// CASE 3d: Rmv - Rmv
				case !other.Insertion && !current.Insertion:
					newIntermediateDiffs = doTransform(transformType7, newIntermediateDiffs, current, other)
				// Fail; should never have gotten to here.
				default:
					panic(fmt.Sprintf("Got to invalid state 3e while transforming [%s] on predessor [%+v], from list [%+v]", current.String(), other, others))
				}
			}
		}
		intermediateDiffs = newIntermediateDiffs
	}
	return intermediateDiffs
}

func doTransform(transformFunc func(current, other *Diff) Diffs, currResults Diffs, current, other *Diff) Diffs {
	return append(currResults, transformFunc(current, other)...)
}

func transformType1(current, other *Diff) Diffs {
	return Diffs{current}
}

func transformType2(current, other *Diff) Diffs {
	return Diffs{current.OffsetDiff(other.Length())}
}

func transformType3(current, other *Diff) Diffs {
	if (other.StartIndex + other.Length()) > current.StartIndex {
		return Diffs{current.OffsetDiff(-(current.StartIndex - other.StartIndex))}
	}
	return Diffs{current.OffsetDiff(-other.Length())}
}

func transformType4(current, other *Diff) Diffs {
	if (other.StartIndex + other.Length()) <= current.StartIndex {
		return Diffs{current.OffsetDiff(-other.Length())}
	} else if (other.StartIndex + other.Length()) >= (current.StartIndex + current.Length()) {
		return Diffs{} // Ignore change
	}
	overlap := other.StartIndex + other.Length() - current.StartIndex
	newDiff := current.OffsetDiff(-other.Length() + overlap)
	newDiff = newDiff.subChangesStartingFrom(overlap)
	return Diffs{newDiff}
}

func transformType5(current, other *Diff) Diffs {
	if current.Length() > other.Length() {
		return Diffs{current.subChangesStartingFrom(other.Length())}
	} // Else do nothing; already done by previous patch.
	return Diffs{}
}

func transformType6(current, other *Diff) Diffs {
	if (current.StartIndex + current.Length()) > other.StartIndex {
		length1 := other.StartIndex - current.StartIndex

		diff1 := current.subChangesEndingAt(length1)
		diff2 := current.subChangesStartingFrom(length1).OffsetDiff(other.Length())

		return Diffs{diff1, diff2}
	}
	return Diffs{current}
}

func transformType7(current, other *Diff) Diffs {
	if (current.StartIndex + current.Length()) > other.StartIndex {
		nonOverlap := other.StartIndex - current.StartIndex

		return Diffs{current.subChangesEndingAt(current.Length() - nonOverlap)}
	}
	return Diffs{current}
}
