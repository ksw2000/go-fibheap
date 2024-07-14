package fibheap_test

import (
	"fmt"

	"github.com/ksw2000/go-fibheap"
)

func ExampleHeap() {
	h := &fibheap.Heap[int, string]{}
	h.Insert(3, "three")
	h.Insert(2, "two")
	h.Insert(1, "one")

	min := h.ExtractMin()
	fmt.Println(min.Key(), min.Value)

	min = h.Min()
	fmt.Println(min.Key(), min.Value)
	// Output: 1 one
	//2 two
}

func ExampleHeap_Decreasing() {
	h := &fibheap.Heap[int, string]{}
	list := []*fibheap.Element[int, string]{}
	list = append(list, h.Insert(5, "one"))
	list = append(list, h.Insert(6, "two"))
	list = append(list, h.Insert(7, "three"))

	h.Decreasing(list[0], 1)
	h.Decreasing(list[1], 2)
	h.Decreasing(list[2], 3)

	min := h.ExtractMin()
	fmt.Println(min.Key(), min.Value)

	min = h.ExtractMin()
	fmt.Println(min.Key(), min.Value)

	min = h.ExtractMin()
	fmt.Println(min.Key(), min.Value)

	// Output: 1 one
	//2 two
	//3 three
}

func ExampleHeap_Remove() {
	h := &fibheap.Heap[int, any]{}
	list := []*fibheap.Element[int, any]{}
	list = append(list, h.Insert(5, nil))
	list = append(list, h.Insert(6, nil))
	list = append(list, h.Insert(7, nil))

	h.Remove(list[0], 0)
	h.Remove(list[1], 0)

	fmt.Println("size:", h.Size())
	fmt.Println("min:", h.Min().Key())

	// Output: size: 3
	//min: 7
}
