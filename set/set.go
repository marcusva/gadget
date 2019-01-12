// Package set features a simple set implementation.
package set

import (
	"fmt"
	"strings"
)

var (
	// MinCapacity defines the default minimum capacity for a MapSet.
	MinCapacity uint = 10
)

// Set is a simple interface for set implementations
type Set interface {
	Contains(...interface{}) bool
	Add(...interface{})
	Remove(...interface{})
	Items() []interface{}
	Iterate(func(interface{}) bool)
	Len() int
	Union(...Set) Set
	Subset(Set) bool
	Superset(Set) bool
	Disjoint(Set) bool
	Intersection(...Set) Set
	Difference(...Set) Set
}

// MapSet is a Set implementation based on the map built-in type
type MapSet map[interface{}]bool

// NewMapSet creates a new MapSet Set.
func NewMapSet(values ...interface{}) *MapSet {
	m := MapSet(make(map[interface{}]bool, 10))
	for _, v := range values {
		m[v] = true
	}
	return &m
}

// Contains checks, if the MapSet contains all of the passed values.
func (m *MapSet) Contains(vals ...interface{}) bool {
	for _, v := range vals {
		if _, ok := (*m)[v]; !ok {
			return false
		}
	}
	return true
}

// Add adds the passed values to the MapSet.
func (m *MapSet) Add(vals ...interface{}) {
	for _, v := range vals {
		(*m)[v] = true
	}
}

// Remove removes the passed values from the MapSet. If one or more values are
// not contained in the MapSet, this will be a no-op for those values.
func (m *MapSet) Remove(vals ...interface{}) {
	for _, v := range vals {
		delete(*m, v)
	}
}

// Union creates a new Set containing the values from the current and passed
// Set. Values contained in both sets, will only occur once in the new Set.
func (m *MapSet) Union(sets ...Set) Set {
	newSet := NewMapSet()
	for k := range *m {
		(*newSet)[k] = true
	}
	for _, set := range sets {
		for _, k := range set.Items() {
			(*newSet)[k] = true
		}
	}
	return newSet
}

// Items returns all values contained in the Set
func (m *MapSet) Items() []interface{} {
	vals := make([]interface{}, 0, len(*m))
	for k := range *m {
		vals = append(vals, k)
	}
	return vals
}

// Iterate iterates over all values of the map, calling the passed callback
// function cb with the value. If the callback returns false, the iteration will
// stop immediately.
func (m *MapSet) Iterate(cb func(interface{}) bool) {
	for k := range *m {
		if !cb(k) {
			return
		}
	}
}

// Len returns the amount of values conained in the Set.
func (m *MapSet) Len() int {
	return len(*m)
}

// Subset checks, if all items of the Mapset are contained in the passed Set.
func (m *MapSet) Subset(s Set) bool {
	for k := range *m {
		if !s.Contains(k) {
			return false
		}
	}
	return true
}

// Superset checks, if all items of the passed Set are contained in the MapSet.
func (m *MapSet) Superset(s Set) bool {
	return s.Subset(m)
}

// Disjoint checks, if the MapSet has no items in common with the other Set.
func (m *MapSet) Disjoint(s Set) bool {
	for k := range *m {
		if s.Contains(k) {
			return false
		}
	}
	return true
}

// Intersection returns a Set, which contains only those items that are common
// to the MapSet and all other Sets provided.
func (m *MapSet) Intersection(sets ...Set) Set {
	newSet := NewMapSet()
	for k := range *m {
		found := true
		for _, s := range sets {
			if !s.Contains(k) {
				found = false
				break
			}
		}
		if found {
			(*newSet)[k] = true
		}
	}
	return newSet
}

// Difference returns a Set, which contains only those items that are unique to
// the MapSet and not available in any of the other Sets provided.
func (m *MapSet) Difference(sets ...Set) Set {
	newSet := NewMapSet()
	for k := range *m {
		unique := true
		for _, s := range sets {
			if s.Contains(k) {
				unique = false
				break
			}
		}
		if unique {
			(*newSet)[k] = true
		}
	}
	return newSet
}

// String returns a string representation of the MapSet.
func (m *MapSet) String() string {
	items := make([]string, 0, len(*m))
	for k := range *m {
		items = append(items, fmt.Sprintf("%#v", k))
	}
	return fmt.Sprintf("MapSet{%s}", strings.Join(items, " "))
}
