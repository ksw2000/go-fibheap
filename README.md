# High-Efficiency Fibonacci Heap

Package fibheap implements a Fibonacci heap. A Fibonacci heap is a data structure for priority queue operations, consisting of a collection of heap-ordered trees.

In our implementation, we do not additionally track the key value of each element. Therefore, users should be aware that they should not insert elements with the same key into the Fibonacci heap.

We compared our package with [Workiva/go-datastructures](https://github.com/Workiva/go-datastructures).

```
goos: linux
goarch: amd64
pkg: github.com/ksw2000/go-fibheap
cpu: Intel(R) Core(TM) i7-4790 CPU @ 3.60GHz
                       │  baseline.txt  │               ours.txt               │
                       │     sec/op     │    sec/op     vs base                │
HeapExtractMin100-8       122.21µ ± 12%   59.62µ ± 16%  -51.21% (p=0.000 n=10)
HeapExtractMin1000-8      6533.6µ ± 11%   804.0µ ± 11%  -87.69% (p=0.000 n=10)
HeapExtractMin10000-8    172.212m ± 10%   8.747m ± 11%  -94.92% (p=0.000 n=10)
HeapExtractMin100000-8   12685.2m ±  2%   135.1m ±  5%  -98.93% (p=0.000 n=10)
geomean                    36.34m         2.744m        -92.45%

                       │  baseline.txt   │               ours.txt               │
                       │      B/op       │     B/op      vs base                │
HeapExtractMin100-8         87.59Ki ± 0%   12.44Ki ± 0%  -85.80% (p=0.000 n=10)
HeapExtractMin1000-8       8408.9Ki ± 0%   140.5Ki ± 0%  -98.33% (p=0.000 n=10)
HeapExtractMin10000-8     817.852Mi ± 0%   1.678Mi ± 0%  -99.79% (p=0.000 n=10)
HeapExtractMin100000-8   77056.52Mi ± 0%   19.84Mi ± 0%  -99.97% (p=0.000 n=10)
geomean                     81.57Mi        497.0Ki       -99.40%

                       │ baseline.txt │              ours.txt               │
                       │  allocs/op   │  allocs/op   vs base                │
HeapExtractMin100-8        298.0 ± 0%    199.0 ± 0%  -33.22% (p=0.000 n=10)
HeapExtractMin1000-8      2.998k ± 0%   1.999k ± 0%  -33.32% (p=0.000 n=10)
HeapExtractMin10000-8     30.01k ± 0%   20.00k ± 0%  -33.36% (p=0.000 n=10)
HeapExtractMin100000-8    300.5k ± 0%   200.0k ± 0%  -33.45% (p=0.000 n=10)
geomean                   9.474k        6.316k       -33.34%
```

The part of code for testing benchmark.

```go
// ours
func benchmarkHeapExtractMin(n int) {
	h := &Heap[float64, struct{}]{}
	for i := 0; i < n; i++ {
		h.Insert(float64(i), struct{}{})
	}
	for i := 0; i < n; i++ {
		h.ExtractMin()
	}
}

// github.com/Workiva/go-datastructures/fibheap
func benchmarkHeapExtractMin(n int) {
	h := heap.NewFloatFibHeap()
	for i := 0; i < n; i++ {
		h.Enqueue(float64(i))
	}
	for i := 0; i < n; i++ {
		h.DequeueMin()
	}
}

func BenchmarkHeapExtractMin100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkHeapExtractMin(100)
	}
}
```

## Example

```go
package main

import (
	"fmt"

	"github.com/ksw2000/go-fibheap"
)

func main() {
	h := &fibheap.Heap[int, string]{}
	nodes := make([]*fibheap.Element[int, string], 10)
	for i := range nodes {
		nodes[i] = h.Insert(i, fmt.Sprint(i))
	}
	// remove the nodes by given a key smaller
	// than all the key in the heap
	h.Remove(nodes[0], -100)
	h.Remove(nodes[2], -100)
	h.Remove(nodes[4], -100)
	h.Remove(nodes[6], -100)
	h.Remove(nodes[8], -100)

	// remove the element with minimum key
	m := h.ExtractMin()
	fmt.Printf("key: %d, val: %s\n", m.Key(), m.Value)
	// key: 1, val: 1

	// or you can fetch the minimum key without removing it
	m = h.Min()
	fmt.Printf("key: %d, val: %s\n", m.Key(), m.Value)
	// key: 3, val: 3

	// you can also decrease the key of elements in the heap
	h.Decreasing(nodes[5], 0)
	m = h.Min()
	fmt.Printf("key: %d, val: %s\n", m.Key(), m.Value)
	// key: 0, val: 5
}
```

