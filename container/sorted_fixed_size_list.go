/*
 	Copyright 2012 Mark Weaver

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/
package container

type Comparer interface {
	Less(other Comparer) bool
}

type node struct {
	prev, next *node
	data       Comparer
}

type SortedFixedSizeList struct {
	maxSize    int
	head, tail *node
	freeList   chan *node
}

func NewSortedFixedSizeList(size int) *SortedFixedSizeList {
	l := SortedFixedSizeList{}
	l.maxSize = size
	l.freeList = make(chan *node, size)
	for i := 0; i < size; i++ {
		l.freeList <- new(node)
	}
	return &l
}

/*
	Add an item only if it is greater than other items in the list, or if the
	list hasn't reached its max size yet.
*/
func (l *SortedFixedSizeList) MaybeAdd(item Comparer) bool {
	var toAdd *node
	select {
	case toAdd = <-l.freeList:
		// This list still isn't its max size
		toAdd.data = item

		if l.head == nil && l.tail == nil {
			// The list is empty so this is a special case.
			l.head = toAdd
			l.tail = toAdd
		} else {
			l.addToList(toAdd)
		}
		return true // Because we added the item
	default:
		// The list has reached it max size
		if l.tail.data.Less(item) {
			// We should add the item to the list. Remove the current tail to
			// reuse the node.
			toAdd = l.tail
			prev := l.tail.prev

			toAdd.next, toAdd.prev, toAdd.data = nil, nil, item

			prev.next = nil
			l.tail = prev

			l.addToList(toAdd)
			return true // Because we added the item
		}
		return false // Because we didn't add the item
	}
}

func (l *SortedFixedSizeList) addToList(toAdd *node) {
	if l.head.data.Less(toAdd.data) {
		// Adding to head.
		toAdd.next = l.head
		l.head.prev = toAdd
		l.head = toAdd
		toAdd.prev = nil
		return
	} else if toAdd.data.Less(l.tail.data) {
		// Adding to tail.
		toAdd.prev = l.tail
		l.tail.next = toAdd
		l.tail = toAdd
		toAdd.next = nil
		return
	}

	// Normal case, adding to middle of list
	current := l.tail
	for current.data.Less(toAdd.data) {
		current = current.prev
		if current == nil {
			// This shouldn't happen!
			// TODO(mww): Return an indication that an error occured.
			return
		}
	}
	toAdd.next = current.next
	toAdd.next.prev = toAdd
	current.next = toAdd
	toAdd.prev = current
}

func (l *SortedFixedSizeList) Iterator() *Iterator {
	i := Iterator{l, nil}
	i.current = l.head
	return &i
}

type Iterator struct {
	list    *SortedFixedSizeList
	current *node
}

func (i *Iterator) HasNext() bool {
	return i.current != nil
}

func (i *Iterator) Next() Comparer {
	v := i.current.data
	i.current = i.current.next
	return v
}

func (i *Iterator) Value() Comparer {
	return i.current.data
}
