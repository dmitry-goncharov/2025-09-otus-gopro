package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("list 1", func(t *testing.T) {
		l := NewList()

		l.PushFront(1)

		require.Equal(t, 1, l.Len())
		require.Equal(t, 1, l.Front().Value)
		require.Equal(t, 1, l.Back().Value)
	})

	t.Run("PushFront 5", func(t *testing.T) {
		l := NewList()

		l.PushFront(5)
		l.PushFront(4)
		l.PushFront(3)
		l.PushFront(2)
		l.PushFront(1)

		require.Equal(t, 5, l.Len())
		require.Equal(t, 1, l.Front().Value)
		require.Equal(t, 5, l.Back().Value)
	})

	t.Run("PushBack 5", func(t *testing.T) {
		l := NewList()

		l.PushBack(1)
		l.PushBack(2)
		l.PushBack(3)
		l.PushBack(4)
		l.PushBack(5)

		require.Equal(t, 5, l.Len())
		require.Equal(t, 1, l.Front().Value)
		require.Equal(t, 5, l.Back().Value)
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}

func TestListExtra(t *testing.T) {
	t.Run("Remove empty", func(t *testing.T) {
		l := NewList()

		l.Remove(l.Front())

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("Remove one", func(t *testing.T) {
		l := NewList()

		l.PushFront(1)

		l.Remove(l.Front())

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("Remove head", func(t *testing.T) {
		l := NewList()

		l.PushFront(5)
		l.PushFront(4)
		l.PushFront(3)
		l.PushFront(2)
		l.PushFront(1)

		l.Remove(l.Front())

		require.Equal(t, 4, l.Len())
		require.Equal(t, 2, l.Front().Value)
		require.Equal(t, 5, l.Back().Value)
	})

	t.Run("Remove tail", func(t *testing.T) {
		l := NewList()

		l.PushFront(5)
		l.PushFront(4)
		l.PushFront(3)
		l.PushFront(2)
		l.PushFront(1)

		l.Remove(l.Back())

		require.Equal(t, 4, l.Len())
		require.Equal(t, 1, l.Front().Value)
		require.Equal(t, 4, l.Back().Value)
	})

	t.Run("Remove middle", func(t *testing.T) {
		l := NewList()

		l.PushFront(5)
		l.PushFront(4)
		l.PushFront(3)
		l.PushFront(2)
		l.PushFront(1)

		l.Remove(l.Front().Next.Next)

		require.Equal(t, 4, l.Len())
		require.Equal(t, 1, l.Front().Value)
		require.Equal(t, 5, l.Back().Value)
	})

	t.Run("MoveToFront empty", func(t *testing.T) {
		l := NewList()

		l.MoveToFront(l.Front())

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("MoveToFront one", func(t *testing.T) {
		l := NewList()

		l.PushFront(1)

		l.MoveToFront(l.Front())

		require.Equal(t, 1, l.Len())
		require.Equal(t, 1, l.Front().Value)
		require.Equal(t, 1, l.Back().Value)
	})

	t.Run("MoveToFront tail", func(t *testing.T) {
		l := NewList()

		l.PushFront(5)
		l.PushFront(4)
		l.PushFront(3)
		l.PushFront(2)
		l.PushFront(1)

		l.MoveToFront(l.Back())

		require.Equal(t, 5, l.Len())
		require.Equal(t, 5, l.Front().Value)
		require.Equal(t, 4, l.Back().Value)
	})

	t.Run("MoveToFront middle", func(t *testing.T) {
		l := NewList()

		l.PushFront(5)
		l.PushFront(4)
		l.PushFront(3)
		l.PushFront(2)
		l.PushFront(1)

		l.MoveToFront(l.Front().Next.Next)

		require.Equal(t, 5, l.Len())
		require.Equal(t, 3, l.Front().Value)
		require.Equal(t, 5, l.Back().Value)
	})
}
