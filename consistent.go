package consistent

import (
	"fmt"
	"hash/crc32"
	"sort"
)

type uints []uint32

func (x uints) Len() int           { return len(x) }
func (x uints) Less(i, j int) bool { return x[i] < x[j] }
func (x uints) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type Consistent struct {
	circle           map[uint32]string
	sortedHashes     uints
	NumberOfReplicas int
	count            int64
}

func New() *Consistent {
	c := new(Consistent)
	c.NumberOfReplicas = 20
	c.circle = make(map[uint32]string)
	return c
}

func (c *Consistent) Add(elt string) {
	for i := 0; i < c.NumberOfReplicas; i++ {
		c.circle[c.hashKey(c.eltKey(elt, i))] = elt
	}
	c.updateSortedHashes()
	c.count++
}

func (c *Consistent) hashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func (c *Consistent) updateSortedHashes() {
	hashes := uints(nil)
	for k := range c.circle {
		hashes = append(hashes, k)
	}
	sort.Sort(hashes)
	c.sortedHashes = hashes
}

// eltKey generates a string key for an element with an index.
func (c *Consistent) eltKey(elt string, idx int) string {
	return fmt.Sprintf("%s|%d", elt, idx)
}
