package lru

/*
   绿色的是字典(map)，存储键和值的映射关系。这样根据某个键(key)查找对应的值(value)的复杂是O(1)，在字典中插入一条记录的复杂度也是O(1)。
   红色的是双向链表(double linked list)实现的队列。将所有的值放到双向链表中，这样，当访问到某个值时，将其移动到队尾的复杂度是O(1)，在队尾新增一条记录以及删除一条记录的复杂度均为O(1)。
*/
import "container/list"

type LRUCache struct {
	maxLen    int64 //最大值
	nLen      int64
	ll        *list.List               //链表，表示其LRU值。队尾的LRU值高
	cache     map[string]*list.Element //map，存储节点，方便取值
	onEvicted func(value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	GetKey() string
	Update(val Value)
}

func New(maxLen int64, onEvited func(Value)) *LRUCache {
	return &LRUCache{
		maxLen:    maxLen,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		onEvicted: onEvited,
	}
}

// Get 查找主要有 2 个步骤，第一步是从字典中找到对应的双向链表的节点，第二步，将该节点移动到队尾。
func (c *LRUCache) Get(key string) (val Value, ok bool) {
	if ele, ok := c.cache[key]; ok { //如果键对应的链表节点存在，返回查找到的值。
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return val, false
}

// RemoveOldset 这里的删除，实际上是缓存淘汰。即移除最近最少访问的节点（队首）
func (c *LRUCache) RemoveOldset() {
	ele := c.ll.Back() //c.ll.Back() 取到队首节点，从链表中删除。

	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key) //从字典中 c.cache 删除该节点的映射关系。
		c.nLen--
		if c.onEvicted != nil {
			c.onEvicted(kv.value) //如果回调函数 OnEvicted 不为 nil，则调用回调函数。
		}
	}
}

//Remove 从Cache里删除
func (c *LRUCache) Remove(key string) {
	ele, ok := c.cache[key]
	if ok {
		c.ll.Remove(ele)
		delete(c.cache, key)
		c.nLen--
		if c.onEvicted != nil {
			v := ele.Value.(*entry)
			c.onEvicted(v.value)
		}
	}
}

func (c *LRUCache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele) //如果键存在，则更新对应节点的值，并将该节点移到队尾。
		kv := ele.Value.(*entry)
		kv.value.Update(value)
		//不存在则是新增场景，首先队尾添加新节点 &entry{key, value}, 并字典中添加 key 和节点的映射关系。
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nLen++
	}
	//更新 c.nBytes，如果超过了设定的最大值 c.maxBytes，则移除最少访问的节点。
	for c.maxLen != 0 && c.maxLen < c.nLen {
		c.RemoveOldset()
	}
}

func (c *LRUCache) Len() int {
	if c.Empty() {
		return 0
	}
	return c.ll.Len()
}
func (c *LRUCache) Empty() bool {
	return c.ll == nil || c.nLen == 0
}

func (c *LRUCache) IsFull() bool {
	if c == nil {
		return false
	}
	return int64(c.Len()) > c.maxLen
}

func (c *LRUCache) Contain(key string) bool {
	if c.Empty() {
		return false
	}
	if _, ok := c.cache[key]; ok {
		return true
	}
	return false
}
