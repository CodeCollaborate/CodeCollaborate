package patching

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
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

// Simplify merges deletion slice within this patch.
// This does not merge insertions, because insertions within the same patch are not actually adjacent
// due to the character that is in between.
func (slice Diffs) Simplify() Diffs {
	if len(slice) == 0 {
		return slice
	}

	result := Diffs{}
	result = append(result, slice[0].clone())

	for i, j := 1, 0; i < len(slice); i++ {
		if slice[i] == nil {
			break
		}
		curr := slice[i]
		prev := result[j]

		if !curr.Insertion && !prev.Insertion && prev.StartIndex+prev.Length() == curr.StartIndex {
			prev.Changes = prev.Changes + curr.Changes
		} else if curr.Insertion && prev.Insertion && prev.StartIndex == curr.StartIndex {
			prev.Changes = prev.Changes + curr.Changes
		} else {
			j++
			result = append(result, curr.clone())
		}
	}

	return result
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
	if length != utf8.RuneCountInString(unescapedChanges) {
		return nil, fmt.Errorf("Length does not match length of change: %d != %s", length, parts[2])
	}

	diff.Changes = unescapedChanges
	return &diff, nil
}

// Length returns the number of characters changed.
func (diff *Diff) Length() int {
	return utf8.RuneCountInString(diff.Changes)
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
	buffer.WriteString(strconv.Itoa(utf8.RuneCountInString(diff.Changes)))
	buffer.WriteString(":")
	buffer.WriteString(url.QueryEscape(diff.Changes))

	return buffer.String()
}

// ConvertToCRLF converts this diff from using LF to CRLF line separators given the base text to patch.
// TODO(wongb): Review this for performance
func (diff *Diff) ConvertToCRLF(base string) *Diff {
	newStartIndex := diff.StartIndex
	newChanges := strings.Replace(diff.Changes, "\n", "\r\n", -1)

	for i := 0; i < newStartIndex && i < utf8.RuneCountInString(base)-1; i++ {
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

	for i := 0; i < diff.StartIndex-1 && i < utf8.RuneCountInString(base)-1; i++ {
		if base[i] == '\r' && base[i+1] == '\n' {
			newStartIndex--
		}
	}

	return NewDiff(diff.Insertion, newStartIndex, newChanges)
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

func (diff *Diff) clone() *Diff {
	return NewDiff(diff.Insertion, diff.StartIndex, diff.Changes)
}
