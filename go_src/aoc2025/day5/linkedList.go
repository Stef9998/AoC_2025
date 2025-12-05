package main

type D5ll struct {
	b    int
	e    int
	next *D5ll
}

func (cur *D5ll) mergeNext() bool {
	next := cur.next
	if cur.e < next.b {
		return false
	}
	if next.e > cur.e {
		cur.e = next.e
	}
	cur.next = next.next
	return true
}

func (cur *D5ll) insertAfter(ins D5ll) {
	ins.next = cur.next
	cur.next = &ins
}

func NewLL(rang ran) D5ll {
	return D5ll{
		b:    rang.b,
		e:    rang.e,
		next: nil,
	}
}
