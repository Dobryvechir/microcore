/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvjson

type OrderedItem struct {
	k interface{}
	v interface{}
}

type OrderedMap struct {
	items []*OrderedItem
	quick map[interface{}]*OrderedItem
}

func CreateOrderedMap(n int) *OrderedMap {
	return &OrderedMap{
		items: make([]*OrderedItem, 0, n),
		quick: make(map[interface{}]*OrderedItem),
	}
}

func findOrderedItem(m *OrderedMap, k interface{}) *OrderedItem {
	switch k.(type) {
	case string, int, int64:
		return m.quick[k]
	default:
		n := len(m.items)
		for i := 0; i < n; i++ {
			if m.items[i].k == k {
				return m.items[i]
			}
		}
	}
	return nil
}

func (m *OrderedMap) Put(k interface{}, v interface{}) bool {
	item := findOrderedItem(m, k)
	if item != nil {
		item.v = v
		return true
	} else {
		item = &OrderedItem{k, v}
		m.items = append(m.items, item)
		switch k.(type) {
		case string, int, int64:
			m.quick[k] = item
		}
	}
	return false
}

func (m *OrderedMap) Get(k interface{}) (interface{}, bool) {
	item := findOrderedItem(m, k)
	if item == nil {
		return nil, false
	}
	return item.v, true
}

func (m *OrderedMap) GetByInt64Key(k int64) (interface{}, bool) {
	item := findOrderedItem(m, k)
	if item == nil {
		item = findOrderedItem(m, int(k))
		if item == nil {
			return nil, false
		}
	}
	return item.v, true
}

func (m *OrderedMap) GetInt64ByInt64Key(k int64) (int64, bool) {
	item := findOrderedItem(m, k)
	if item == nil {
		item = findOrderedItem(m, int(k))
		if item == nil {
			return 0, false
		}
	}
	r := item.v
	switch r.(type) {
	case int64:
		return r.(int64), true
	case int:
		return int64(r.(int)), true
	}
	return -1, false
}

func (m *OrderedMap) Size() int {
	return len(m.items)
}

func (m *OrderedMap) GetAtSafe(pos int) (interface{}, interface{}, bool) {
	item := m.items[pos]
	if item == nil {
		return nil, nil, false
	}
	return item.k, item.v, true
}

func (m *OrderedMap) GetAt(pos int) (interface{}, interface{}) {
	item := m.items[pos]
	if item == nil {
		return nil, nil
	}
	return item.k, item.v
}

func (m *OrderedMap) Remove(k interface{}) bool {
	switch k.(type) {
	case string, int, int64:
		item := m.quick[k]
		if item == nil {
			return false
		}
		delete(m.quick, k)
	}
	n := len(m.items)
	for i := 0; i < n; i++ {
		if m.items[i].k == k {
			m.items = append(m.items[:i], m.items[i+1:]...)
			return true
		}
	}
	return false
}

func (m *OrderedMap) RemoveAt(pos int) bool {
	if pos >= len(m.items) {
		return false
	}
	item := m.items[pos]
	m.items = append(m.items[:pos], m.items[pos+1:]...)
	k := item.k
	switch k.(type) {
	case string, int, int64:
		item := m.quick[k]
		if item == nil {
			return false
		}
		delete(m.quick, k)
	}
	return true
}

func (m *OrderedMap) IsSimpleObject() bool {
	for _, it := range m.items {
		if _, ok := it.k.(string); !ok {
			return false
		}
	}
	return true
}

func (m *OrderedMap) IsSimpleArray(base int) bool {
	n := len(m.items) + base
	n64 := int64(n)
	base64 := int64(base)
	for _, it := range m.items {
		k := it.k
		switch k.(type) {
		case int:
			ki := k.(int)
			if ki < base || ki >= n {
				return false
			}
		case int64:
			ki := k.(int64)
			if ki < base64 || ki >= n64 {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func CreateOrderedMapByMap(m map[interface{}]interface{}) *OrderedMap {
	r := CreateOrderedMap(len(m))
	for k, v := range m {
		r.Put(k, v)
	}
	return r
}

func CreateOrderedMapByArray(m []interface{}, base int) *OrderedMap {
	r := CreateOrderedMap(len(m))
	for k, v := range m {
		r.Put(k+base, v)
	}
	return r
}
