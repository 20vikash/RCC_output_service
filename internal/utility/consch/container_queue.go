package consch

type ConQueue struct {
	Language string
	Code     string
	Jid      string

	Next *ConQueue
}

func CreateConQueue(language, code, jid string) *ConQueue {
	return &ConQueue{
		Language: language,
		Code:     code,
		Jid:      jid,
		Next:     nil,
	}
}

func (q *ConQueue) AddCode() {

}

func (q *ConQueue) LatestCode() {

}

func (q *ConQueue) AckCode() {

}

func (q *ConQueue) PeekCode() {

}
