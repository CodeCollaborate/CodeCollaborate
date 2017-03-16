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
	// Consolidate first two patches, if there are any
	for _, patchB := range patches[1:] {
		indexA := -1
		indexB := -1
		resultDiffs := Diffs{}
		currIndex := 0

		diffA, indexA := getNextDiff(patchA, indexA, false)
		diffB, indexB := getNextDiff(patchB, indexB, false)

		// Convenience update functions
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

		// If either diff is nil, we have hit the end
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

			if !diffA.Insertion && !isNoOp(diffA) { // Cases A1 through A3
				commit(diffA, -1)
				currIndex += lenA
				getNextDiffA()
			} else if diffB.Insertion && !isNoOp(diffB) { // Cases G1 through G3
				commit(diffB, -1)
				getNextDiffB()
			} else {
				// Commit changes and update currIndex as needed
				switch {
				// Cases B1 through B3
				case !isNoOp(diffA) && diffA.Insertion && !isNoOp(diffB) && !diffB.Insertion:
				// Commit type 0
				default:
					// Do nothing
					break

				// Cases C1 through C3
				case !isNoOp(diffA) && diffA.Insertion && isNoOp(diffB):
					switch {
					// Commit type 1 (Cases C1, C2)
					case lenA < lenB, lenA == lenB:
						commit(diffA, -1)

					// Commit type 2 (Case C3)
					default:
						commit(diffA, lenB)
					}

				// Cases H1 through H3
				case isNoOp(diffA) && !isNoOp(diffB) && !diffB.Insertion:
					switch {
					// Commit Type 3 (Case H1)
					case lenA < lenB:
						commit(diffB, lenA)
						currIndex += lenA
					// Commit Type 4 (Case H2)
					case lenA == lenB:
						commit(diffB, -1)
						currIndex += lenA
					// Commit Type 5 (Case H3)
					default:
						commit(diffB, -1)
						currIndex += lenB
					}
				// Case I1 through I3
				case isNoOp(diffA) && isNoOp(diffB):
					switch {
					// Commit type 6 (Cases I1, I2)
					case lenA < lenB, lenA == lenB:
						currIndex += lenA
					// Commit Type 7 (Cases I3)
					default:
						currIndex += lenB
					}
				}

				// Update diff, get a new one if needed.
				switch {
				// Iteration Type 1 (All cases in 1 column)
				case lenA < lenB:
					if isNoOp(diffB) {
						diffB.StartIndex += lenA
					} else {
						diffB.Changes = diffB.Changes[lenA:]
					}
					getNextDiffA()

				// Iteration Type 2 (All cases in 2 column)
				case lenA == lenB:
					getNextDiffA()
					getNextDiffB()

				// Iteration Type 3 (All cases in 3 column)
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
	// If we are just starting, currIndex will be -1
	if currIndex == -1 {
		if !wasNoOp { // If we are just starting, wasNoOp will be false - return the starting no-op
			return NewDiff(true, 0, ""), -1
		}

		// If the length of the patches is 0, return nil diff
		if patch.Changes.Len() <= 0 {
			return nil, -1
		}

		// Else return the first diff
		return patch.Changes[0].clone(), 0
	}

	// If previous one was a noOp, return either the end or the next actual diff
	if wasNoOp {
		// Return nil diff if we have gone past the end
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
		if currDiff.Insertion {
			return NewDiff(true, currDiff.StartIndex, ""), currIndex
		}

		// Else, return a noop diff after the removed block
		return NewDiff(true, currDiff.StartIndex+currDiff.Length(), ""), currIndex

	}

	//return the next diff if it is adjacent, or a noOp otherwise.
	nextDiff := patch.Changes[currIndex+1]

	// If next diff is adjacent, return it directly
	if nextDiff.StartIndex == currDiff.StartIndex || !currDiff.Insertion && currDiff.StartIndex+currDiff.Length() >= nextDiff.StartIndex {
		return nextDiff.clone(), currIndex + 1
	}

	// Otherwise, return a new no-op diff
	// If the current diff is an insertion, start the no-op diff at the current location
	if currDiff.Insertion {
		return NewDiff(true, currDiff.StartIndex, ""), currIndex
	}

	// Else, start it at the end of the removal
	return NewDiff(true, currDiff.StartIndex+currDiff.Length(), ""), currIndex

}

func isNoOp(diff *Diff) bool {
	return diff.Length() == 0
}

func noOpLength(diff *Diff, patch *Patch, currIndex int) int {
	// If this is the last diff, return the remaining untouched length of the document
	if currIndex+1 >= patch.Changes.Len() {
		return patch.DocLength - diff.StartIndex + 1
	}

	// Else, return the length from the no-op patch's startIndex to the next patch's startIndex
	return patch.Changes[currIndex+1].StartIndex - diff.StartIndex

}
