package consch

type ConQueue struct {
	Language string
	Code     string
	Jid      string

	Next *ConQueue
}

func (q *ConQueue) AddCode() {

}

func (q *ConQueue) LatestCode() {

}

func (q *ConQueue) AckCode() {

}

func (q *ConQueue) PeekCode() {

}
