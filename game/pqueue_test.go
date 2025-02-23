package game

import "testing"

func TestPriorityQueueOperations(t *testing.T) {
	var pq pqueue

	// Test empty queue
	if len(pq) != 0 {
		t.Error("New queue should be empty")
	}

	// Test push
	pos1 := Pos{X: 1, Y: 1}
	pq = pq.push(pos1, 5)
	if len(pq) != 1 {
		t.Error("Queue should have one element after push")
	}

	// Test push maintaining heap property
	pos2 := Pos{X: 2, Y: 2}
	pq = pq.push(pos2, 3) // Higher priority (lower number)
	if pq[0].priority != 3 {
		t.Error("Heap property not maintained after push")
	}

	// Test pop
	pq, popped := pq.pop()
	if popped != pos2 {
		t.Errorf("Expected pos {2,2}, got pos {%d,%d}", popped.X, popped.Y)
	}
	if len(pq) != 1 {
		t.Error("Queue should have one element after pop")
	}

	// Test pop last element
	pq, popped = pq.pop()
	if popped != pos1 {
		t.Errorf("Expected pos {1,1}, got pos {%d,%d}", popped.X, popped.Y)
	}
	if len(pq) != 0 {
		t.Error("Queue should be empty after popping all elements")
	}

	// Test pop empty queue
	pq, popped = pq.pop()
	if popped != (Pos{}) {
		t.Error("Popping empty queue should return zero value Pos")
	}
}
