package patching

type point struct {
	x int
	y int
}

type snake struct {
	start *point
	mid   *point
	end   *point
}

func myersDiff(str1, str2 string) (Diffs, error) {
	N := len(str1)
	M := len(str2)

	snakes := map[point]*snake{}

	maxXByK := map[int]int{}
	maxXByK[1] = 0
	maxXByK[-1] = -1

	found := false
	for d := 0; d <= N+M && !found; d++ {
		for k := -d; k <= d; k += 2 {
			down := k == -d || (k != d && maxXByK[k-1] < maxXByK[k+1])

			prevK := k - 1
			if down {
				prevK = k + 1
			}

			xStart := maxXByK[prevK]
			yStart := xStart - prevK

			xMid := xStart + 1
			yMid := yStart
			if down {
				xMid = xStart
				yMid = yStart + 1
			}

			diagonalsTaken := 0
			xEnd := xMid
			yEnd := yMid

			// While we are still in the grid, and the next character matches, take a diagonal
			for (xEnd < len(str1) && yEnd < len(str2)) && (str1[xEnd] == str2[yEnd]) {
				xEnd++
				yEnd++
				diagonalsTaken++
			}

			maxXByK[k] = xEnd

			snake := snake{
				start: &point{
					x: xStart,
					y: yStart,
				},
				mid: &point{
					x: xMid,
					y: yMid,
				},
				end: &point{
					x: xEnd,
					y: yEnd,
				},
			}

			// Group deletions together if possible.
			if val := snakes[*snake.end]; val == nil || val.start.x <= snake.start.x {
				snakes[*snake.end] = &snake
			}

			if xEnd >= len(str1) && yEnd >= len(str2) {
				found = true
			}
		}
	}

	flippedResults := Diffs{}

	p := &point{x: len(str1), y: len(str2)}
	for p.y != -1 {
		s := snakes[*p]

		if s.start.y != -1 {
			if s.mid.x == s.start.x {
				flippedResults = append(flippedResults, NewDiff(true, s.mid.x, string(str2[s.mid.y-1])))
			} else {
				flippedResults = append(flippedResults, NewDiff(false, s.mid.x-1, string(str1[s.mid.x-1])))
			}
		}

		p = s.start
	}

	results := make(Diffs, len(flippedResults))
	for i := 0; i < len(flippedResults); i++ {
		results[len(flippedResults)-1-i] = flippedResults[i]
	}

	return results, nil
}
