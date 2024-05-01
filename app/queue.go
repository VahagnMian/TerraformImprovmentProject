package main

import (
	"fmt"
)

// Queue represents a queue that holds a slice of strings.
type ApplyQueue struct {
	elements []string
}

// Enqueue adds an element to the end of the queue.
func (q *ApplyQueue) Enqueue(element string) {
	q.elements = append(q.elements, element)
}

// Dequeue removes an element from the front of the queue and returns it.
// If the queue is empty, it returns an empty string and an error.
func (q *ApplyQueue) Dequeue() (string, error) {
	if len(q.elements) == 0 {
		return "", fmt.Errorf("queue is empty")
	}
	element := q.elements[0]
	q.elements = q.elements[1:]
	return element, nil
}

// IsEmpty returns true if the queue is empty.
func (q *ApplyQueue) IsEmpty() bool {
	return len(q.elements) == 0
}

func RunQueue(element string ) {
	q := ApplyQueue{}

	q.Enqueue(element)

	// Dequeue elements
	for !q.IsEmpty() {
		element, err := q.Dequeue()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(element)
	}
}
