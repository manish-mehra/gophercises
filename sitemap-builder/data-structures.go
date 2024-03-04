package main

type Queue struct {
	items []string
}

// Enqueue adds a string to the end of the queue
func (q *Queue) Enqueue(item string) {
	q.items = append(q.items, item)
}

// Dequeue removes and returns the first string from the queue
func (q *Queue) Dequeue() (string, bool) {
	if len(q.items) == 0 {
		return "", false // Queue is empty
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// IsEmpty checks if the queue is empty
func (q *Queue) IsEmpty() bool {
	return len(q.items) == 0
}
