package patching

import "fmt"

// ConsolidatePatches consolidates patch others with patch A.
// Patches should be fed into this function in dependency order (A -> B -> C)
func ConsolidatePatches(patchA *Patch, others ...*Patch) *Patch {
	for _, patchB := range others {
		indexA := -1
		indexB := -1
		resultDiffs := Diffs{}

		diffA, indexA := getNextDiff(patchA, indexA, false)
		diffB, indexB := getNextDiff(patchB, indexB, false)

		getNextDiffA := func() {
			diffA, indexA = getNextDiff(patchA, indexA, isNoOp(diffA))
		}
		getNextDiffB := func() {
			diffB, indexB = getNextDiff(patchB, indexB, isNoOp(diffB))
		}

		currIndex := 0
		for diffA != nil && diffB != nil {
			fmt.Printf("DiffA: [%s], DiffB: [%s]\n", diffA.String(), diffB.String())

			if !diffA.Insertion && !isNoOp(diffA) {
				resultDiffs = append(resultDiffs, NewDiff(diffA.Insertion, currIndex, diffA.Changes))
				currIndex += diffA.Length()
				getNextDiffA()
				continue
			}
			if diffB.Insertion && !isNoOp(diffB) {
				resultDiffs = append(resultDiffs, NewDiff(diffB.Insertion, currIndex, diffB.Changes))
				//currIndex += diffB.Length()
				getNextDiffB()
				continue
			}

			if isNoOp(diffA) && isNoOp(diffB) {
				lenA := noOpLength(diffA, patchA, indexA)
				lenB := noOpLength(diffB, patchB, indexB)
				if lenA > lenB {
					currIndex += lenB
					diffA.StartIndex += lenB
					getNextDiffB()
				} else if lenA == lenB {
					currIndex += lenB
					getNextDiffA()
					getNextDiffB()
				} else {
					currIndex += lenA
					diffB.StartIndex += lenA
					getNextDiffA()
				}
			} else if diffA.Insertion && !isNoOp(diffA) && !diffB.Insertion {
				if diffA.Length() > diffB.Length() {
					diffA.Changes = diffA.Changes[diffB.Length():]
					getNextDiffB()
				} else if diffA.Length() == diffB.Length() {
					getNextDiffA()
					getNextDiffB()
				} else {
					diffB.Changes = diffB.Changes[diffA.Length():]
					getNextDiffA()
				}
			} else if diffA.Insertion && isNoOp(diffB) {
				lenB := noOpLength(diffB, patchB, indexB)
				if diffA.Length() > lenB {
					resultDiffs = append(resultDiffs, NewDiff(diffA.Insertion, currIndex, diffA.Changes[:lenB]))
					//currIndex += lenB
					diffA.Changes = diffA.Changes[lenB:]
					getNextDiffB()
				} else if diffA.Length() == lenB {
					resultDiffs = append(resultDiffs, NewDiff(diffA.Insertion, currIndex, diffA.Changes))
					//currIndex += diffA.Length()
					getNextDiffA()
					getNextDiffB()
				} else {
					resultDiffs = append(resultDiffs, NewDiff(diffA.Insertion, currIndex, diffA.Changes))
					//currIndex += diffA.Length()
					diffB.StartIndex += diffA.Length()
					getNextDiffA()
				}
			} else if isNoOp(diffA) && !diffB.Insertion {
				lenA := noOpLength(diffA, patchA, indexA)
				if lenA > diffB.Length() {
					resultDiffs = append(resultDiffs, NewDiff(diffB.Insertion, currIndex, diffB.Changes))
					currIndex += diffB.Length()
					diffA.StartIndex += diffB.Length()
					getNextDiffB()
				} else if lenA == diffB.Length() {
					resultDiffs = append(resultDiffs, NewDiff(diffB.Insertion, currIndex, diffB.Changes))
					currIndex += diffB.Length()
					getNextDiffA()
					getNextDiffB()
				} else {
					resultDiffs = append(resultDiffs, NewDiff(diffB.Insertion, currIndex, diffB.Changes[:lenA]))
					currIndex += lenA
					diffB.Changes = diffB.Changes[lenA:]
					getNextDiffA()
				}
			} else {
				panic("WHAT IS GOING ON")
			}
		}

		fmt.Println(resultDiffs)
		//return

		patchA = NewPatch(patchA.BaseVersion, resultDiffs, patchA.DocLength)
	}
	return patchA
}

func getNextDiff(patch *Patch, currIndex int, wasNoOp bool) (*Diff, int) {
	if currIndex == -1 {
		if !wasNoOp {
			return NewDiff(true, 0, ""), -1
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
