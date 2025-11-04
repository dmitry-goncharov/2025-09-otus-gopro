package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len  int
	head *ListItem
	tail *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	l.len++
	elem := &ListItem{
		Value: v,
	}
	if l.head != nil {
		l.head.Prev = elem
		elem.Next = l.head
	}
	l.head = elem
	if l.tail == nil {
		l.tail = elem
	}
	return l.head
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.len++
	elem := &ListItem{
		Value: v,
	}
	if l.tail != nil {
		l.tail.Next = elem
		elem.Prev = l.tail
	}
	l.tail = elem
	if l.head == nil {
		l.head = elem
	}
	return l.tail
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}
	l.len--
	if l.len == 0 {
		l.head, l.tail = nil, nil
		return
	}
	if i == l.head {
		l.head = i.Next
		i.Next.Prev = nil
		return
	}
	if i == l.tail {
		l.tail = i.Prev
		i.Prev.Next = nil
		return
	}
	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil {
		return
	}
	if i == l.head {
		return
	}
	i.Prev.Next = i.Next
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if i == l.tail {
		l.tail = i.Prev
	}
	l.head.Prev = i
	i.Next = l.head
	l.head = i
	l.head.Prev = nil
}
