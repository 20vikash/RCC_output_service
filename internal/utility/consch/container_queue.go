package consch

type ConQueue struct {
	Language string
	Code     string
	Jid      string

	Next *ConQueue
}

var head *ConQueue

func CreateConQueue() *ConQueue {
	return head
}

func (q *ConQueue) AddCode() {

}

func (q *ConQueue) LatestCode() {

}

func (q *ConQueue) AckCode() {

}

func (q *ConQueue) PeekCode() {

}
