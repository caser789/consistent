package consistent

import (
	"errors"
	"fmt"
	"hash/crc32"
	"sort"
)

type uints []uint32

func (x uints) Len() int           { return len(x) }
func (x uints) Less(i, j int) bool { return x[i] < x[j] }
func (x uints) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

var ErrEmptyCircle = errors.New("empty circle")

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

// Remove removes an element from the hash.
func (c *Consistent) Remove(elt string) {
	for i := 0; i < c.NumberOfReplicas; i++ {
		delete(c.circle, c.hashKey(c.eltKey(elt, i)))
	}
	c.updateSortedHashes()
	c.count--
}

// Get returns an element close to where name hashes to in the circle.
func (c *Consistent) Get(name string) (string, error) {
	if len(c.circle) == 0 {
		return "", ErrEmptyCircle
	}
	key := c.hashKey(name)
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}
	i := sort.Search(len(c.sortedHashes), f)
	if i >= len(c.sortedHashes) {
		i = 0
	}
	return c.circle[c.sortedHashes[i]], nil
}

// Get two returns the two closest distinct elements to the name input in the circle.
func (c *Consistent) GetTwo(name string) (string, string, error) {
	if len(c.circle) == 0 {
		return "", "", ErrEmptyCircle
	}
	key := c.hashKey(name)
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}
	i := sort.Search(len(c.sortedHashes), f)
	if i >= len(c.sortedHashes) {
		i = 0
	}
	a := c.circle[c.sortedHashes[i]]

	if c.count == 1 {
		return a, "", nil
	}

	start := i
	i++
	if i >= len(c.sortedHashes) {
		i = 0
	}
	var b string
	for i != start {
		b = c.circle[c.sortedHashes[i]]
		if b != a {
			break
		}
		i++
		if i >= len(c.sortedHashes) {
			i = 0
		}
	}
	return a, b, nil
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
