package patching

import "errors"

// TransformResult is a struct aggregating the results of the TransformPatches function
type TransformResult struct {
	PatchXPrime *Patch
	PatchYPrime *Patch
}

// ErrorBaseDocumentLengthsDifferent is the error thrown when base document lengths for two patches are different
var ErrorBaseDocumentLengthsDifferent = errors.New("Base document lengths for patchX and patchY were different")

// ErrorIllegalStateNoOpXYLen is the error thrown when we get into an invalid state (-1)
var ErrorIllegalStateNoOpXYLen = errors.New("Got to invalid state based on noOpXLen and noOpYLen")

// TransformPatches takes two patches, and produces their component opposites (A -> A'), (B -> B')
func TransformPatches(patchX *Patch, patchY *Patch) (*TransformResult, error) {
	if patchX.DocLength != patchY.DocLength {
		return nil, ErrorBaseDocumentLengthsDifferent
	}

	patchXPrime := Diffs{}
	patchYPrime := Diffs{}

	indexX := -1
	indexY := -1

	startIndexX := 0
	startIndexY := 0

	noOpXLen := 0
	var diffX *Diff
	noOpYLen := 0
	var diffY *Diff

	getNextDiffX := func() {
		indexX++

		if indexX < patchX.Changes.Len() {
			diffX = patchX.Changes[indexX].clone()

			noOpXLen = diffX.StartIndex
			if indexX > 0 {
				prev := patchX.Changes[indexX-1]
				if prev.Insertion || prev.StartIndex == diffX.StartIndex {
					noOpXLen = diffX.StartIndex - prev.StartIndex
				} else {
					noOpXLen = diffX.StartIndex - (prev.StartIndex + prev.Length())
				}
			}
		} else {
			// last no-op
			if patchX.Changes.Len() != 0 {
				prev := patchX.Changes[patchX.Changes.Len()-1]
				noOpXLen = patchX.DocLength - prev.StartIndex
				if !prev.Insertion {
					noOpXLen -= prev.Length()
				}
			} else {
				noOpXLen = patchX.DocLength
			}
			diffX = nil
		}
	}

	getNextDiffY := func() {
		indexY++

		if indexY < patchY.Changes.Len() {
			diffY = patchY.Changes[indexY].clone()

			noOpYLen = diffY.StartIndex
			if indexY > 0 {
				prev := patchY.Changes[indexY-1]
				if prev.Insertion || prev.StartIndex == diffY.StartIndex {
					noOpYLen = diffY.StartIndex - prev.StartIndex
				} else {
					noOpYLen = diffY.StartIndex - (prev.StartIndex + prev.Length())
				}
			}
		} else {
			// last no-op
			if patchY.Changes.Len() != 0 {
				prev := patchY.Changes[patchY.Changes.Len()-1]
				noOpYLen = patchY.DocLength - prev.StartIndex
				if !prev.Insertion {
					noOpYLen -= prev.Length()
				}
			} else {
				noOpYLen = patchY.DocLength
			}
			diffY = nil
		}
	}

	getNextDiffX()
	getNextDiffY()

	for {
		// Min of noOpXLen, noOpYLen
		noOpLength := noOpXLen
		if noOpYLen < noOpXLen {
			noOpLength = noOpYLen
		}
		// Max of noOpLength, 0
		if noOpLength < 0 {
			noOpLength = 0
		}
		startIndexX += noOpLength
		startIndexY += noOpLength
		noOpXLen -= noOpLength
		noOpYLen -= noOpLength
		if diffX == nil && diffY == nil {
			break
		}

		diffXLen := noOpXLen
		if diffX != nil {
			diffXLen = diffX.Length()
		}
		diffYLen := noOpYLen
		if diffY != nil {
			diffYLen = diffY.Length()
		}

		if diffX != nil && diffX.Insertion && noOpXLen == 0 {
			patchXPrime = append(patchXPrime, NewDiff(true, startIndexX, diffX.Changes))

			startIndexY += diffXLen

			getNextDiffX()
			continue
		} else if diffY != nil && diffY.Insertion && noOpYLen == 0 {
			patchYPrime = append(patchYPrime, NewDiff(true, startIndexY, diffY.Changes))

			startIndexX += diffYLen

			getNextDiffY()
			continue
		}

		if noOpXLen == 0 && noOpYLen == 0 {
			if diffXLen < diffYLen {
				// Nothing to commit

				// Already been deleted; remove first lenX characters from diffY
				diffY.Changes = diffY.Changes[diffXLen:]

				// Done with this diffX
				getNextDiffX()
			} else if diffXLen == diffYLen {
				// Both deleted the same text; nothing to commit

				getNextDiffX()
				getNextDiffY()
			} else {
				// Nothing to commit

				// Already been deleted; remove first lenX characters from diffY
				diffX.Changes = diffX.Changes[diffYLen:]

				// Done with this diffX
				getNextDiffY()
			}
		} else if noOpXLen == 0 && noOpYLen > 0 {
			commitLength := noOpYLen
			if diffXLen < commitLength {
				commitLength = diffXLen
			}

			patchXPrime = append(patchXPrime, NewDiff(false, startIndexX, diffX.Changes[0:commitLength]))

			startIndexX += commitLength
			diffX.Changes = diffX.Changes[commitLength:]

			noOpYLen -= commitLength

			// If we have exhausted the entire diffY, proceed to next diff
			if diffX.Length() == 0 {
				getNextDiffX()
			}
		} else if noOpXLen > 0 && noOpYLen == 0 {
			commitLength := noOpXLen
			if diffYLen < commitLength {
				commitLength = diffYLen
			}

			patchYPrime = append(patchYPrime, NewDiff(false, startIndexY, diffY.Changes[0:commitLength]))

			startIndexY += commitLength
			diffY.Changes = diffY.Changes[commitLength:]

			noOpXLen -= commitLength

			// If we have exhausted the entire diffY, proceed to next diff
			if diffY.Length() == 0 {
				getNextDiffY()
			}
		} else {
			return nil, ErrorIllegalStateNoOpXYLen
		}
	}

	// Less efficient, but simpler
	newDocXLen := patchX.DocLength
	for _, diff := range patchY.Changes {
		if diff.Insertion {
			newDocXLen += diff.Length()
		} else {
			newDocXLen -= diff.Length()
		}
	}

	newDocYLen := patchY.DocLength
	for _, diff := range patchX.Changes {
		if diff.Insertion {
			newDocYLen += diff.Length()
		} else {
			newDocYLen -= diff.Length()
		}
	}

	return &TransformResult{
		NewPatch(patchY.BaseVersion+1, patchXPrime, newDocXLen),
		NewPatch(patchX.BaseVersion+1, patchYPrime, newDocYLen),
	}, nil
}
