package diff

type Action int

const (
	ActionDelete Action = iota
	ActionAdd
	ActionEqual
	ActionChange
)

type Change struct {
	Action      Action
	OldPosition int
	OldElement  string
	NewPosition int
	NewElement  string
}

func TraverseBalanced(a, b []string) (result []Change) {
	matches := lcs(a, b)
	aLen, bLen := len(a), len(b)
	var ai, bj, mb int
	ma := -1

	getItem := func(x []string, i int) string {
		if i < len(x) {
			return x[i]
		} else {
			return ""
		}
	}

	for {
		for {
			ma++
			if !(ma < len(matches) && matches[ma] < 0) {
				break
			}
		}

		if ma >= len(matches) {
			break
		}
		mb = matches[ma]

		for ai < ma || bj < mb {
			ax := getItem(a, ai)
			bx := getItem(b, bj)

			x := ai < ma
			y := bj < mb
			if x && y {
				result = append(result, Change{
					Action:      ActionChange,
					OldPosition: ai,
					OldElement:  ax,
					NewPosition: bj,
					NewElement:  bx,
				})
				ai++
				bj++
			} else if x && !y {
				result = append(result, Change{
					Action:      ActionDelete,
					OldPosition: ai,
					OldElement:  ax,
					NewPosition: bj,
					NewElement:  bx,
				})
				ai++
			} else if !x && y {
				result = append(result, Change{
					Action:      ActionAdd,
					OldPosition: ai,
					OldElement:  ax,
					NewPosition: bj,
					NewElement:  bx,
				})
				bj++
			}
		}

		ax := getItem(a, ai)
		bx := getItem(b, bj)
		result = append(result, Change{
			Action:      ActionEqual,
			OldPosition: ai,
			OldElement:  ax,
			NewPosition: bj,
			NewElement:  bx,
		})
		ai++
		bj++
	}

	for ai < aLen || bj < bLen {
		ax := getItem(a, ai)
		bx := getItem(b, bj)

		x := ai < aLen
		y := bj < bLen
		if x && y {
			result = append(result, Change{
				Action:      ActionChange,
				OldPosition: ai,
				OldElement:  ax,
				NewPosition: bj,
				NewElement:  bx,
			})
			ai++
			bj++
		} else if x && !y {
			result = append(result, Change{
				Action:      ActionDelete,
				OldPosition: ai,
				OldElement:  ax,
				NewPosition: bj,
				NewElement:  bx,
			})
			ai++
		} else if !x && y {
			result = append(result, Change{
				Action:      ActionAdd,
				OldPosition: ai,
				OldElement:  ax,
				NewPosition: bj,
				NewElement:  bx,
			})
			bj++
		}
	}
	return
}

func lcs(a, b []string) []int {
	aStart, bStart := 0, 0
	aEnd, bEnd := len(a)-1, len(b)-1

	length := len(a)
	if len(a) < len(b) {
		length = len(b)
	}
	result := make([]int, length)
	for i := range result {
		result[i] = -1
	}

	for aStart <= aEnd && bStart <= bEnd && a[aStart] == b[bStart] {
		result[aStart] = bStart
		aStart++
		bStart++
	}
	bStart = aStart

	for aStart <= aEnd && bStart <= bEnd && a[aEnd] == b[bEnd] {
		result[aEnd] = bEnd
		aEnd -= 1
		bEnd -= 1
	}

	bMatches := map[string][]int{}
	for i, line := range b[bStart : bEnd+1] {
		bMatches[line] = append(bMatches[line], i)
	}
	thresh := []int{}

	type link struct {
		last           *link
		aIndex, bIndex int
	}
	links := make([]*link, length)

	for i := aStart; i <= aEnd; i++ {
		ai := a[i]
		bm := bMatches[ai]
		k := -1
		for l := len(bm) - 1; l >= 0; l-- {
			j := bm[l]
			if k > 0 && thresh[k] > j && thresh[k-1] < j {
				thresh[k] = j
			} else {
				k, thresh = replaceNextLarger(thresh, j, k)
			}
			if k != -1 {
				var last *link
				if k > 0 {
					last = links[k-1]
				}
				links[k] = &link{
					last:   last,
					aIndex: i,
					bIndex: j,
				}
			}
		}
	}

	if len(thresh) > 0 {
		link := links[len(thresh)-1]
		for link != nil {
			result[link.aIndex] = link.bIndex
			link = link.last
		}
	}
	return result
}

func replaceNextLarger(enum []int, value, lastIndex int) (int, []int) {
	if len(enum) == 0 || value > enum[len(enum)-1] {
		enum = append(enum, value)
		return len(enum) - 1, enum
	}

	// Binary search for the insertion point
	if lastIndex < 0 {
		lastIndex = len(enum)
	}
	firstIndex := 0
	for firstIndex <= lastIndex {
		i := (firstIndex + lastIndex) >> 1
		found := enum[i]
		if value == found {
			return -1, enum
		}

		if value > found {
			firstIndex = i + 1
		} else {
			lastIndex = i - 1
		}
	}

	// The insertion point is in first_index; overwrite the next larger value.
	enum[firstIndex] = value
	return firstIndex, enum
}
