package consch

import (
	"sync"
)

type conNode struct {
	Language string
	Code     string
	Jid      string
	next     *conNode
}

type ConQueue struct {
	Front *conNode
	Rear  *conNode
	Size  int
	mu    sync.Mutex
}

func CreateConQueue() *ConQueue {
	return &ConQueue{}
}

func (q *ConQueue) AddCode(language, code, jid string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	n := &conNode{Language: language, Code: code, Jid: jid}
	if q.Front == nil {
		q.Front, q.Rear = n, n
	} else {
		q.Rear.next = n
		q.Rear = n
	}
	q.Size++
}

func (q *ConQueue) AckCode() *conNode {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.Front == nil {
		return nil
	}

	t := q.Front
	q.Front = q.Front.next
	if q.Front == nil {
		q.Rear = nil
	}
	q.Size--
	return t
}

func (q *ConQueue) LatestCode() *conNode {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.Front
}
