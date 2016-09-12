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

type Diff struct {
	Insertion  bool
	StartIndex int
	Changes    string
}

func NewDiff(insertion bool, startIndex int, changes string) *Diff {
	return &Diff{
		Insertion:  insertion,
		StartIndex: startIndex,
		Changes:    changes,
	}
}

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

	diff.StartIndex, err = strconv.Atoi(parts[0])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Invalid offset: %s", parts[0]))
	}

	switch parts[1][0] {
	case '+':
		diff.Insertion = true
	case '-':
		diff.Insertion = false
	default:
		return nil, errors.New(fmt.Sprintf("Invalid operation: %s", parts[1][0]))
	}

	length, err := strconv.Atoi(parts[1][1:])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Invalid length: %s", parts[1][1:]))
	}

	unescapedChanges, err := url.QueryUnescape(parts[2])
	if err != nil {
		return nil, err
	}

	if length != len(unescapedChanges) {
		return nil, errors.New(fmt.Sprintf("Length does not match length of change: %d != %s", length, parts[2]))
	}

	diff.Changes = unescapedChanges

	return &diff, nil
}

func (diff *Diff) Length() int {
	return len(diff.Changes)
}

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

func (diff *Diff) ConvertToLF(base string) *Diff {
	newStartIndex := diff.StartIndex
	newChanges := strings.Replace(diff.Changes, "\r\n", "\n", -1)

	for i := 0; i < diff.StartIndex && i < len(base)-1; i++ {
		if base[i] == '\n' {
			newStartIndex--
		}
	}

	return NewDiff(diff.Insertion, newStartIndex, newChanges)
}

func (diff *Diff) GetUndo() *Diff {
	return NewDiff(!diff.Insertion, diff.StartIndex, diff.Changes)
}

func (diff *Diff) OffsetDiff(offset int) *Diff {
	return NewDiff(diff.Insertion, diff.StartIndex+offset, diff.Changes)
}

func (diff *Diff) GetSubChanges(start, end int) *Diff {
	return NewDiff(diff.Insertion, diff.StartIndex, diff.Changes[start:end])
}

func (diff *Diff) GetSubChangesStartingFrom(start int) *Diff {
	return NewDiff(diff.Insertion, diff.StartIndex, diff.Changes[start:])
}

func (diff *Diff) GetSubChangesEndingAt(end int) *Diff {
	return NewDiff(diff.Insertion, diff.StartIndex, diff.Changes[:end])
}

func (diff *Diff) transform(others []*Diff) []*Diff {

	intermediateDiffs := []*Diff{}
	intermediateDiffs = append(intermediateDiffs, diff)

	for _, other := range others {
		newIntermediateDiffs := []*Diff{}
		for _, current := range intermediateDiffs {
			switch {

			// CASE 1: IndexA < IndexB
			case other.StartIndex < current.StartIndex:
				switch {
				// CASE 1a, 1b: Ins - Ins, Ins - Rmv
				case other.Insertion && current.Insertion, other.Insertion && !current.Insertion:
					newDiff := current.OffsetDiff(other.Length())
					newIntermediateDiffs = append(newIntermediateDiffs, newDiff)

				// CASE 1c: Rmv - Ins
				case !other.Insertion && current.Insertion:
					if (other.StartIndex + other.Length()) > current.StartIndex {
						newDiff := diff.OffsetDiff(-(current.StartIndex - other.StartIndex))
						newIntermediateDiffs = append(newIntermediateDiffs, newDiff)
					} else {
						newDiff := diff.OffsetDiff(-other.Length())
						newIntermediateDiffs = append(newIntermediateDiffs, newDiff)
					}

				// CASE 1d: Rmv - Rmv
				case !other.Insertion && !current.Insertion:
					if (other.StartIndex + other.Length()) <= current.StartIndex {
						newDiff := diff.OffsetDiff(-other.Length())
						newIntermediateDiffs = append(newIntermediateDiffs, newDiff)
					} else if (other.StartIndex + other.Length()) >= (current.StartIndex + current.Length()) {
						// do nothing
					} else {
						overlap := other.StartIndex + other.Length() - current.StartIndex
						newDiff := diff.OffsetDiff(-other.Length() + overlap)
						newDiff = newDiff.GetSubChangesStartingFrom(overlap)
						newIntermediateDiffs = append(newIntermediateDiffs, newDiff)
					}

				// Fail; should never have gotten to here.
				default:
					panic(fmt.Sprintf("Got to invalid state 1e while transforming [%s] on predessors [%+v]", current.String(), others))
				}

			// CASE 2: IndexA = IndexB
			case other.StartIndex == current.StartIndex:

				switch {

				// CASES 2a, 2b: Ins - Ins, Ins - Rmv
				case other.Insertion && current.Insertion, other.Insertion && !current.Insertion:
					newDiff := current.OffsetDiff(other.Length())
					newIntermediateDiffs = append(newIntermediateDiffs, newDiff)

				// CASE 2c: Rmv - Ins
				// Do nothing
				case !other.Insertion && current.Insertion:
					newIntermediateDiffs = append(newIntermediateDiffs, current)

				// CASE 2d: Rmv - Ins
				case !other.Insertion && !current.Insertion:
					if current.Length() > other.Length() {
						newDiff := current.GetSubChangesStartingFrom(other.Length())
						newIntermediateDiffs = append(newIntermediateDiffs, newDiff)
					}
				// Else do nothing; already done by previous patch.

				// Fail; should never have gotten to here.
				default:
					panic(fmt.Sprintf("Got to invalid state 2e while transforming [%s] on predessors [%+v]", current.String(), others))
				}
			case other.StartIndex > current.StartIndex:
				switch {

				// CASE 3a, 3c: Ins - Ins, Rmv - Ins
				case other.Insertion && current.Insertion, !other.Insertion && current.Insertion:
					newIntermediateDiffs = append(newIntermediateDiffs, current)

				// CASE 3b: Ins - Rmv
				case other.Insertion && !current.Insertion:
					if (current.StartIndex + current.Length()) > other.StartIndex {
						length1 := other.StartIndex - current.StartIndex

						diff1 := current.GetSubChangesEndingAt(length1)
						diff2 := current.GetSubChangesStartingFrom(length1).OffsetDiff(other.Length())

						newIntermediateDiffs = append(newIntermediateDiffs, diff1, diff2)
					} else {
						newIntermediateDiffs = append(newIntermediateDiffs, current)
					}

				// CASE 3d: Rmv - Rmv
				case !other.Insertion && !current.Insertion:
					if (current.StartIndex + current.Length()) > other.StartIndex {
						nonOverlap := other.StartIndex - current.StartIndex

						newDiff := current.GetSubChangesEndingAt(current.Length() - nonOverlap)

						newIntermediateDiffs = append(newIntermediateDiffs, newDiff)
					} else {
						newIntermediateDiffs = append(newIntermediateDiffs, current)
					}

				// Fail; should never have gotten to here.
				default:
					panic(fmt.Sprintf("Got to invalid state 3e while transforming [%s] on predessors [%+v]", current.String(), others))
				}
			}
		}
		intermediateDiffs = newIntermediateDiffs
	}
	return intermediateDiffs
}
