package fibheap

import (
	"testing"
)

func assert(t *testing.T, actual any, expected int) {
	if expected != actual.(int) {
		t.Errorf("❌ expected: %d actual: %d\n", expected, actual.(int))
	}
}

func TestHeapInsert(t *testing.T) {
	count := 32
	h := &Heap[int, any]{}
	for i := 0; i < count; i++ {
		h.Insert(i, nil)
	}
	if h.Size() != count {
		t.Fail()
	}
}

func TestHeapExtract(t *testing.T) {
	count := 32
	h := &Heap[int, any]{}
	for i := 0; i < count; i++ {
		h.Insert(i, nil)
	}
	for i := 0; i < count-1; i++ {
		assert(t, h.ExtractMin().key, i)
	}
}

func TestHeapDecreasing(t *testing.T) {
	h := &Heap[int, any]{}
	elements := make([]*Element[int, any], 512)
	for i := 0; i < 512; i++ {
		elements[i] = h.Insert(i, i)
	}

	for i := 0; i < 10; i++ {
		h.ExtractMin()
	}

	for i := 99; i >= 50; i-- {
		h.Decreasing(elements[i], -1000)
		x := h.ExtractMin()
		if h.min.key == -1000 {
			break
		}
		assert(t, x.Value, i)
	}
}

func TestHeapRemove(t *testing.T) {
	h := &Heap[int, any]{}
	elements := make([]*Element[int, any], 100)
	for i := 0; i < 100; i++ {
		elements[i] = h.Insert(i, i)
	}
	for i := 0; i < 50; i++ {
		h.Remove(elements[i], -1)
	}
	if h.Min() != elements[50] {
		t.Fail()
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Should panic()")
		}
	}()
	h.Remove(elements[99], 100)
}

func TestUnion(t *testing.T) {
	h := &Heap[int, any]{}
	g := &Heap[int, any]{}
	for i := 0; i < 10; i++ {
		h.Insert(i, nil)
	}
	for i := 10; i < 20; i++ {
		g.Insert(i, nil)
	}

	k := h.Union(g)
	for i := 0; i < 20; i++ {
		x := k.ExtractMin()
		assert(t, x.Key(), i)
	}

	if h.min != nil || h.elements != 0 {
		t.Fatal("h should be clear after Union")
	}
	if g.min != nil || g.elements != 0 {
		t.Fatal("g should be clear after Union")
	}
}
