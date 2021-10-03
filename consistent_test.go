package consistent

import (
	"sort"
	"testing"
)

func checkNum(num, expected int, t *testing.T) {
	if num != expected {
		t.Errorf("expected %d, got %d", expected, num)
	}
}

func TestNew(t *testing.T) {
	x := New()
	if x == nil {
		t.Errorf("expected obj")
	}
	checkNum(x.NumberOfReplicas, 20, t)
}

func TestAdd(t *testing.T) {
	x := New()
	x.Add("abcdefg")
	checkNum(len(x.circle), 20, t)
	checkNum(len(x.sortedHashes), 20, t)
	if sort.IsSorted(x.sortedHashes) == false {
		t.Errorf("expected sorted hashes to be sorted")
	}
	x.Add("qwer")
	checkNum(len(x.circle), 40, t)
	checkNum(len(x.sortedHashes), 40, t)
	if sort.IsSorted(x.sortedHashes) == false {
		t.Errorf("expected sorted hashes to be sorted")
	}
}

func TestRemove(t *testing.T) {
	x := New()
	x.Add("abcdefg")
	x.Remove("abcdefg")
	checkNum(len(x.circle), 0, t)
	checkNum(len(x.sortedHashes), 0, t)
}

func TestRemoveNonExisting(t *testing.T) {
	x := New()
	x.Add("abcdefg")
	x.Remove("abcdefghijk")
	checkNum(len(x.circle), 20, t)
}
