package tools

type Queue []*Item

func (q Queue) Len() int { return len(q) }

func (q *Queue) Push(x interface{}) {
	item := x.(*Item)
	*q = append(*q, item)
}

func (q *Queue) Pop() *Item {
	item := (*q)[0]
	*q = (*q)[1:]
	return item
}

func (q *Queue) Top() *Item {
	item := (*q)[0]
	return item
}

func (q *Queue) Last() *Item {
	item := (*q)[len(*q)-1]
	return item
}

func (q *Queue) RemoveLast() {
	*q = (*q)[:len(*q)-1]
}

func NewQueue(size int) Queue {
	return make(Queue, size)
}
