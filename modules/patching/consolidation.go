package patching

import (
	"errors"
)

// ConsolidatePatches consolidates patch others with patch A.
// Patches should be fed into this function in dependency order (A -> B -> C)
func ConsolidatePatches(patches []*Patch) (*Patch, error) {
	if len(patches) <= 0 {
		return nil, errors.New("ConsolidatePatches: No patches provided")
	}

	patchA := patches[0]
	for _, patchB := range patches[1:] {
		indexA := -1
		indexB := -1
		resultDiffs := Diffs{}
		currIndex := 0

		diffA, indexA := getNextDiff(patchA, indexA, false)
		diffB, indexB := getNextDiff(patchB, indexB, false)

		getNextDiffA := func() {
			diffA, indexA = getNextDiff(patchA, indexA, isNoOp(diffA))
		}
		getNextDiffB := func() {
			diffB, indexB = getNextDiff(patchB, indexB, isNoOp(diffB))
		}
		commit := func(diff *Diff, numChars int) {
			if numChars == -1 {
				resultDiffs = append(resultDiffs, NewDiff(diff.Insertion, currIndex, diff.Changes))
			} else {
				resultDiffs = append(resultDiffs, NewDiff(diff.Insertion, currIndex, diff.Changes[:numChars]))
			}
		}

		for diffA != nil && diffB != nil {
			// Get lengths of each diff
			lenA := 0
			lenB := 0
			if isNoOp(diffA) {
				lenA = noOpLength(diffA, patchA, indexA)
			} else {
				lenA = diffA.Length()
			}
			if isNoOp(diffB) {
				lenB = noOpLength(diffB, patchB, indexB)
			} else {
				lenB = diffB.Length()
			}

			if !diffA.Insertion && !isNoOp(diffA) {
				commit(diffA, -1)
				currIndex += lenA
				getNextDiffA()
			} else if diffB.Insertion && !isNoOp(diffB) {
				commit(diffB, -1)
				getNextDiffB()
			} else {
				// Commit changes and update currIndex as needed
				switch {
				case !isNoOp(diffA) && diffA.Insertion && !isNoOp(diffB) && !diffB.Insertion:
				// do nothing
				case !isNoOp(diffA) && diffA.Insertion && isNoOp(diffB):
					switch {
					case lenA < lenB, lenA == lenB:
						commit(diffA, -1)
					default:
						commit(diffA, lenB)
					}
				case isNoOp(diffA) && !isNoOp(diffB) && !diffB.Insertion:
					switch {
					case lenA < lenB:
						commit(diffB, lenA)
						currIndex += lenA
					case lenA == lenB:
						commit(diffB, -1)
						currIndex += lenA
					default:
						commit(diffB, -1)
						currIndex += lenB
					}
				case isNoOp(diffA) && isNoOp(diffB):
					switch {
					case lenA < lenB, lenA == lenB:
						currIndex += lenA
					default:
						currIndex += lenB
					}
				}

				// Update the diff and get new ones if needed
				switch {
				case lenA < lenB:
					if isNoOp(diffB) {
						diffB.StartIndex += lenA
					} else {
						diffB.Changes = diffB.Changes[lenA:]
					}
					getNextDiffA()
				case lenA == lenB:
					getNextDiffA()
					getNextDiffB()
				default:
					if isNoOp(diffA) {
						diffA.StartIndex += lenB
					} else {
						diffA.Changes = diffA.Changes[lenB:]
					}
					getNextDiffB()
				}
			}
		}
		patchA = NewPatch(patchA.BaseVersion, resultDiffs, patchA.DocLength)
	}
	return patchA, nil
}

func getNextDiff(patch *Patch, currIndex int, wasNoOp bool) (*Diff, int) {
	if currIndex == -1 {
		if !wasNoOp {
			return NewDiff(true, 0, ""), -1
		} else if patch.Changes.Len() <= 0 {
			return nil, -1
		}
		return patch.Changes[0].clone(), 0

	}

	// If previous one was a noOp, return either the end or the next actual diff
	if wasNoOp {
		// Return end if we have gone past it.
		if currIndex+1 >= patch.Changes.Len() {
			return nil, -1
		}
		// Return next diff otherwise.
		return patch.Changes[currIndex+1].clone(), currIndex + 1
	}
	// Else return the next value
	currDiff := patch.Changes[currIndex]
	// If we have no more slice, and our current diff does not go to the end, return a new noop diff
	if currIndex+1 >= patch.Changes.Len() {
		// If it is an insertion, return new noOp diff at start index.
		// Else, return a noop diff after the removed block
		if currDiff.Insertion {
			return NewDiff(true, currDiff.StartIndex, ""), currIndex
		}
		return NewDiff(true, currDiff.StartIndex+currDiff.Length(), ""), currIndex

	}
	//return the next diff if it is adjacent, or a noOp otherwise.
	nextDiff := patch.Changes[currIndex+1]
	if nextDiff.StartIndex == currDiff.StartIndex || !currDiff.Insertion && currDiff.StartIndex+currDiff.Length() >= nextDiff.StartIndex {
		return nextDiff.clone(), currIndex + 1
	}
	if currDiff.Insertion {
		return NewDiff(true, currDiff.StartIndex, ""), currIndex
	}
	return NewDiff(true, currDiff.StartIndex+currDiff.Length(), ""), currIndex

}

func isNoOp(diff *Diff) bool {
	return diff.Length() == 0
}

func noOpLength(diff *Diff, patch *Patch, currIndex int) int {
	if currIndex+1 >= patch.Changes.Len() {
		return patch.DocLength - diff.StartIndex + 1
	}
	return patch.Changes[currIndex+1].StartIndex - diff.StartIndex

}
