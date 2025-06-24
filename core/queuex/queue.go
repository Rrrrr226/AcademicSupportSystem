package queuex

import "container/list"

type Queue struct {
	list *list.List
}

func New() *Queue {
	return &Queue{list.New()}
}

// Push inserts element to the queue
func (q *Queue) Push(value interface{}) {
	q.list.PushBack(value)
}

// Front returns first element of the queue
func (q *Queue) Front() interface{} {
	it := q.list.Front()
	if it != nil {
		return it.Value
	}
	return nil
}

// Back returns last element of the queue
func (q *Queue) Back() interface{} {
	it := q.list.Back()
	if it != nil {
		return it.Value
	}
	return nil
}

// Pop returns and deletes first element of the queue
func (q *Queue) Pop() interface{} {
	it := q.list.Front()
	if it != nil {
		q.list.Remove(it)
		return it.Value
	}
	return nil
}

// Size returns size of the queue
func (q *Queue) Size() int {
	return q.list.Len()
}

// Empty returns whether queue is empty
func (q *Queue) Empty() bool {
	return q.list.Len() == 0
}

// Clear clears the queue
func (q *Queue) Clear() {
	for !q.Empty() {
		q.Pop()
	}
}
