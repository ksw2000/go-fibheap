// Package fibheap implements a Fibonacci heap.
// A Fibonacci heap is a data structure for priority queue operations,
// consisting of a collection of heap-ordered trees. The amortized time
// complexity of Fibonacci heap operations are as follows:
// fetching the minimum is Θ(1),
// extracting the minimum is O(log n),
// inserting is Θ(1),
// decreasing a key is Θ(1),
// and merging two heaps is Θ(1).
//
// In our implementation, we do not additionally track the key value of each
// element. Therefore, users should be aware that they should not insert
// elements with the same key into the Fibonacci heap.
package fibheap

import (
	"golang.org/x/exp/constraints"
)

type Element[K constraints.Ordered, V any] struct {
	p        *Element[K, V]
	r        *Element[K, V]
	l        *Element[K, V]
	children *Element[K, V]
	// store mark in the LSB
	degree uint32
	key    K
	// The value stored with this element.
	Value V
}

func (e *Element[K, V]) getDegree() int {
	return int(e.degree >> 1)
}

func (e *Element[K, V]) increaseDegree() {
	e.degree += 2
}

func (e *Element[K, V]) decreaseDegree() {
	e.degree -= 2
}

func (e *Element[K, V]) getMark() bool {
	return e.degree&1 == 1
}

func (e *Element[K, V]) clearMark() {
	e.degree = e.degree &^ 1
}

func (e *Element[K, V]) setMark() {
	e.degree = e.degree | 1
}

// Key returns the key of the element e
func (e *Element[K, V]) Key() K {
	return e.key
}

// append appends new element m to the left of element n
func (n *Element[K, V]) append(m *Element[K, V]) *Element[K, V] {
	if m == nil {
		return n
	}
	if n == nil {
		m.l, m.r = m, m
		return m
	}
	r := n.r
	n.r = m
	m.l = n
	m.r = r
	r.l = m
	return n
}

// Heap represents the fibonacci heap.
type Heap[K constraints.Ordered, V any] struct {
	elements int
	min      *Element[K, V]
}

// Size returns the number of elements in the heap h
func (h *Heap[K, V]) Size() int {
	return h.elements
}

// Insert inserts the key-value pair (key, value) to the heap h and returns the
// inserted element with amortized running time Θ(1)
func (h *Heap[K, V]) Insert(key K, value V) *Element[K, V] {
	n := &Element[K, V]{key: key, Value: value}
	h.elements++
	h.min = h.min.append(n)
	if n.key < h.min.key {
		h.min = n
	}
	return n
}

// Min fetches the minimum key from the heap h with running time Θ(1)
func (h *Heap[K, V]) Min() *Element[K, V] {
	return h.min
}

// ExtractMin() fetches and removes the minimum key from the heap h with
// amortized running time O(log n)
func (h *Heap[K, V]) ExtractMin() *Element[K, V] {
	if h == nil || h.min == nil {
		return nil
	}

	if h.min.children != nil {
		h.min.children.p = nil
		for c := h.min.children.r; c != h.min.children; c = c.r {
			c.p = nil
		}
		l := h.min.children.l
		r := h.min.r
		h.min.r = h.min.children
		h.min.children.l = h.min
		l.r = r
		r.l = l
	}

	z := h.min
	if h.min.r == h.min.l && h.min.r == h.min {
		h.min = nil
	} else {
		h.min.l.r = h.min.r
		h.min.r.l = h.min.l
		h.min = h.min.r
		h.consolidate()
	}

	return z
}

// d returns math.Floor(math.Log2(n))
func d(a int) int {
	i := 0
	for a > 1 {
		a = a >> 1
		i++
	}
	return i
}

func (h *Heap[K, V]) consolidate() {
	a := make([]*Element[K, V], d(h.elements)+1)
	end := h.min.l
	for w := h.min; ; {
		next := w.r
		x := w
		d := x.getDegree()
		for a[d] != nil {
			y := a[d]
			if y.key < x.key {
				x, y = y, x
			}
			h.link(y, x)
			a[d] = nil
			d++
		}
		a[d] = x
		if w == end {
			break
		}
		w = next
	}
	h.min = nil
	for _, node := range a {
		if node == nil {
			continue
		}
		node.l.r = node.r
		node.r.l = node.l
		node.l = node
		node.r = node

		if h.min == nil {
			h.min = node
			continue
		}
		h.min = h.min.append(node)
		if node.key < h.min.key {
			h.min = node
		}
	}
}

// link removes y from the root list, and makes y a children of x.
func (h *Heap[K, V]) link(y, x *Element[K, V]) {
	// remove y form the root list
	y.l.r = y.r
	y.r.l = y.l

	x.children = x.children.append(y)

	x.increaseDegree()
	y.p = x
	y.clearMark()
}

// Decreasing decreases the key of element with the minimum key with amortized
// running time Θ(1). If the new key k is larger or equal than the key of x,
// Decreasing does nothing.
func (h *Heap[K, V]) Decreasing(x *Element[K, V], key K) {
	if key >= x.key {
		return
	}
	x.key = key
	p := x.p
	if p != nil && x.key < p.key {
		h.cut(x, p)
		h.cascadingCut(p)
	}
	if x.key < h.min.key {
		h.min = x
	}
}

// Remove removes the element x by given a key minimumKey which is smaller than
// any key in the heap h.
func (h *Heap[K, V]) Remove(x *Element[K, V], minimumKey K) {
	h.Decreasing(x, minimumKey)
	if n := h.Min(); n != x {
		panic("fibheap: Remove will remove unexpected element")
	}
	h.ExtractMin()
}

// cut cuts the link between x and its parent p and makes x a root.
func (h *Heap[K, V]) cut(x, p *Element[K, V]) {
	p.decreaseDegree()

	if x == x.r {
		p.children = nil
	} else {
		x.l.r = x.r
		x.r.l = x.l

		if p.children == x {
			p.children = x.r
		}
	}

	// add x to the list of h
	x.l = x
	x.r = x
	x.p = nil
	x.clearMark()
	h.min = h.min.append(x)
}

// cascadingCut handles the ancestral consequences of cutting an element.
func (h *Heap[K, V]) cascadingCut(y *Element[K, V]) {
	z := y.p
	if z != nil {
		if !y.getMark() {
			y.setMark()
		} else {
			h.cut(y, z)
			h.cascadingCut(z)
		}
	}
}

// Union unions the two fibonacci heaps h and g, and returns the new fibonacci
// heap with amortized running time Θ(1). The heap h and g will be reset after
// unioning.
func (h *Heap[K, V]) Union(g *Heap[K, V]) *Heap[K, V] {
	if h == nil || g == nil {
		panic("fibheap: Union expects non-nil heap h and g")
	}

	m := &Heap[K, V]{
		elements: g.elements + h.elements,
	}
	if h.min != nil && g.min != nil {
		l := g.min.l
		r := h.min.r
		h.min.r = g.min
		g.min.l = h.min
		l.r = r
		r.l = l

		if h.min.key < g.min.key {
			m.min = h.min
		} else {
			m.min = g.min
		}
	} else if h.min != nil {
		m.min = h.min
	} else {
		m.min = g.min
	}

	// clear heap h and heap g
	h.min = nil
	h.elements = 0
	g.min = nil
	g.elements = 0

	return m
}
