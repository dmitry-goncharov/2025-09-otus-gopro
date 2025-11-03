package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type element struct {
	key *Key
	val interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	v, ok := c.items[key]
	if ok {
		el := v.Value.(element)
		el.val = value
		v.Value = el
		c.queue.MoveToFront(v)
	} else {
		el := element{
			key: &key,
			val: value,
		}
		v = c.queue.PushFront(el)
		c.items[key] = v
		if c.queue.Len() > c.capacity {
			tail := c.queue.Back()
			c.queue.Remove(tail)
			delete(c.items, *tail.Value.(element).key)
		}
	}
	return ok
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	v, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(v)
		return v.Value.(element).val, ok
	}
	return nil, ok
}

func (c *lruCache) Clear() {
	for k, v := range c.items {
		delete(c.items, k)
		c.queue.Remove(v)
	}
}
