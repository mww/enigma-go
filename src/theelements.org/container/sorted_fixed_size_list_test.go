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

import (
	"fmt"
	"testing"
)

type Int struct {
	value int32
}

func (i *Int) Less(other Comparer) bool {
	o, ok := other.(*Int)
	if !ok {
		return false
	}
	return i.value < o.value
}

func (i *Int) String() string {
	return fmt.Sprintf("%d", i.value)
}

func assertValues(t *testing.T, expected []int32, actual *Iterator) {
	for i := 0; actual.HasNext(); i++ {
		val := actual.Next()
		v, ok := val.(*Int)
		if !ok {
			t.Errorf("Was expecting an Int, got something else: %s", val)
		}
		if expected[i] != v.value {
			t.Errorf("Expected %d, got %d", expected[i], v.value)
		}
	}
}

func TestAddItemToList(t *testing.T) {
	l := NewSortedFixedSizeList(3)
	l.MaybeAdd(&Int{1})
	l.MaybeAdd(&Int{2})
	l.MaybeAdd(&Int{3})
	l.MaybeAdd(&Int{4})
	l.MaybeAdd(&Int{5})

	assertValues(t, []int32{5, 4, 3}, l.Iterator())
}

func TestAddItemsToListOutOfOrder(t *testing.T) {
	l := NewSortedFixedSizeList(3)
	l.MaybeAdd(&Int{3})
	l.MaybeAdd(&Int{13})
	l.MaybeAdd(&Int{1})
	l.MaybeAdd(&Int{8})
	l.MaybeAdd(&Int{25})
	l.MaybeAdd(&Int{7})

	assertValues(t, []int32{25, 13, 8}, l.Iterator())
}
