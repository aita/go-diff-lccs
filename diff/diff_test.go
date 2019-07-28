package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTraverseBalanced(t *testing.T) {
	xs := TraverseBalanced(
		[]string{"A", "B", "C", "B", "D", "A", "B"},
		[]string{"B", "D", "C", "A", "B", "A"},
	)
	assert.Equal(t, []Change{
		{ActionDelete, 0, "A", 0, "B"},
		{ActionEqual, 1, "B", 0, "B"},
		{ActionDelete, 2, "C", 1, "D"},
		{ActionDelete, 3, "B", 1, "D"},
		{ActionEqual, 4, "D", 1, "D"},
		{ActionAdd, 5, "A", 2, "C"},
		{ActionEqual, 5, "A", 3, "A"},
		{ActionEqual, 6, "B", 4, "B"},
		{ActionAdd, 7, "", 5, "A"},
	}, xs)
}

func TestLCS(t *testing.T) {
	xs := lcs(
		[]string{"A", "B", "C", "B", "D", "A", "B"},
		[]string{"B", "D", "C", "A", "B", "A"},
	)
	assert.Equal(t, []int{-1, 0, -1, -1, 1, 3, 4}, xs)
}
