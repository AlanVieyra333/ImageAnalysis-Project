// This example demonstrates a Priority queue built using the heap interface.
package tools

import (
	"container/heap"
)

// An Item is something we manage in a Priority queue.
type Item struct {
	Value    string // The Value of the item; arbitrary.
	Priority int    // The Priority of the item in the queue.
	// The Index is needed by update and is maintained by the heap.Interface methods.
	Index int // The Index of the item in the heap.
	// Coordenada
	X, Y int
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, Priority so we use greater than here.
	if pq[i].Priority < pq[j].Priority {
		return true
	} else if pq[i].Priority == pq[j].Priority {
		if pq[i].Y < pq[j].Y {
			return true
		} else if pq[i].Y == pq[j].Y {
			if pq[i].X < pq[j].X {
				return true
			}
		}
	}
	return false
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) Top() *Item {
	old := *pq
	item := old[0]
	return item
} //*/

// update modifies the Priority and Value of an Item in the queue.
func (pq *PriorityQueue) Update(item *Item, Value string, Priority int) {
	item.Value = Value
	item.Priority = Priority
	heap.Fix(pq, item.Index)
}

func NewPriorityQueue(size int) PriorityQueue {
	return make(PriorityQueue, size)
}

/*/ This example creates a PriorityQueue with some items, adds and manipulates an item,
// and then removes the items in Priority order.
func main() {
	// Some items and their priorities.
	items := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}

	// Create a Priority queue, put the items in it, and
	// establish the Priority queue (heap) invariants.
	pq := make(PriorityQueue, len(items))
	i := 0
	for Value, Priority := range items {
		pq[i] = &Item{
			Value:    Value,
			Priority: Priority,
			Index:    i,
		}
		i++
	}
	heap.Init(&pq)

	// Insert a new item and then modify its Priority.
	item := &Item{
		Value:    "orange",
		Priority: 1,
	}
	heap.Push(&pq, item)
	pq.update(item, item.Value, 5)

	// Take the items out; they arrive in decreasing Priority order.
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%.2d:%s ", item.Priority, item.Value)
	}
}
*/
