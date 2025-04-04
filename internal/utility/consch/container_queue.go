package consch

type conNode struct {
	Language string
	Code     string
	Jid      string

	next *conNode
}

func createNode(language, code, jid string) *conNode {
	return &conNode{
		Language: language,
		Code:     code,
		Jid:      jid,

		next: nil,
	}
}

type ConQueue struct {
	Front *conNode
}

func CreateConQueue() *ConQueue {
	return &ConQueue{
		Front: createNode("", "", ""),
	}
}

func (q *ConQueue) AddCode(language, code, jid string) {
	if q.Front == nil {
		q.Front = createNode(language, code, jid)
	} else {
		n := createNode(language, code, jid)
		q.Front.next = n
		q.Front = n
	}
}

func (q *ConQueue) LatestCode() {

}

func (q *ConQueue) AckCode() {

}

func (q *ConQueue) PeekCode() {

}
