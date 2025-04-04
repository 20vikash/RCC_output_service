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
	return &ConQueue{}
}

func (q *ConQueue) AddCode(language, code, jid string) {
	if q.Front == nil {
		q.Front = createNode(language, code, jid)
	} else {
		n := createNode(language, code, jid)
		t := q.Front
		for t.next != nil {
			t = t.next
		}

		t.next = n
	}
}

func (q *ConQueue) LatestCode() *conNode {
	if q.Front != nil {
		return q.Front
	}

	return nil
}

func (q *ConQueue) AckCode() *conNode {
	if q.Front == nil {
		return nil
	}

	t := q.Front
	q.Front = q.Front.next

	return t
}
