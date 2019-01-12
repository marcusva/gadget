package set_test

import (
	"fmt"
	"github.com/marcusva/gadget/set"
	"github.com/marcusva/gadget/testing/assert"
	"strings"
	"testing"
)

func contains(items []interface{}, val interface{}) bool {
	for _, v := range items {
		if v == val {
			return true
		}
	}
	return false
}

func TestNewMapSet(t *testing.T) {
	m := set.NewMapSet()
	assert.NotNil(t, m)

	filled := set.NewMapSet("test", 1, true)
	assert.NotNil(t, m)
	assert.Equal(t, filled.Len(), 3)
}

func TestMapSetContains(t *testing.T) {
	m := set.NewMapSet("test", 1, true)
	assert.Equal(t, m.Contains("test"), true, "Contains failed for 'test'")
	assert.Equal(t, m.Contains(1), true, "Contains failed for 1")
	assert.Equal(t, m.Contains(true), true, "Contains failed for true")
	assert.Equal(t, m.Contains(false), false, "Contains failed for false")
	assert.Equal(t, m.Contains(27), false, "Contains failed for 27")
	assert.Equal(t, m.Contains(m.Items()...), true, "Contains failed for all items of the set")
	assert.Equal(t, m.Contains(3, true), false, "Contains failed multiple args with invalid values")
}

func TestMapSetAdd(t *testing.T) {
	data := []string{"test", "1", "true"}
	m := set.NewMapSet()
	for _, v := range data {
		m.Add(v)
		assert.Equal(t, m.Contains(v), true, "Add failed for ", v)
	}
	for _, v := range data {
		m.Add(v)
		assert.Equal(t, m.Len(), 3, m, "has the wrong length")
	}
}

func TestMapSetRemove(t *testing.T) {
	data := []string{"test", "1", "true"}
	m := set.NewMapSet("test", "1", "true")
	for _, v := range data {
		assert.Equal(t, m.Contains(v), true, "invalid initial MapSet for ", v)
		m.Remove(v)
		assert.Equal(t, m.Contains(v), false, "Remove failed for ", v)
	}
}
func TestMapSetUnion(t *testing.T) {
	m1 := set.NewMapSet(1, 2, 3)
	m2 := set.NewMapSet("1", "2", "3")
	m3 := set.NewMapSet(true, false)

	union1 := m1.Union(m2)
	for _, s := range []set.Set{m1, m2} {
		for _, v := range s.Items() {
			assert.Equal(t, union1.Contains(v), true, fmt.Sprintf("Union 1 failed for %v:%T", v, v))
		}
	}

	union2 := m3.Union(m1, m2)
	for _, s := range []set.Set{m1, m2, m3} {
		for _, v := range s.Items() {
			assert.Equal(t, union2.Contains(v), true, fmt.Sprintf("Union 2 failed for %v:%T", v, v))
		}
	}
}
func TestMapSetItems(t *testing.T) {
	m := set.NewMapSet(1, "test", true)
	items := m.Items()
	assert.Equal(t, len(items), m.Len(), "Items failed: ", m, items)
	assert.Equal(t, contains(items, 1), true, "Items failed for 1")
	assert.Equal(t, contains(items, "test"), true, "Items failed for 'test'")
	assert.Equal(t, contains(items, true), true, "Items failed for true")
}
func TestMapSetIterate(t *testing.T) {
	ref := make(map[interface{}]bool)
	filler := func(i interface{}) bool {
		ref[i] = true
		return true
	}

	m := set.NewMapSet(1, "test", true)
	m.Iterate(filler)
	for _, k := range m.Items() {
		if _, ok := ref[k]; !ok {
			t.Errorf("Iterate failed for %v", k)
		}
	}

	ref = make(map[interface{}]bool)
	nofill := func(i interface{}) bool {
		if i == true {
			return false
		}
		ref[i] = true
		return true
	}
	m.Iterate(nofill)
	assert.Equal(t, len(ref) == m.Len(), false, "Iterate failed for a false retval")
}

func TestMapSetLen(t *testing.T) {
	m := set.NewMapSet()
	assert.Equal(t, m.Len(), 0, "Len failed for NewMapSet()")

	m = set.NewMapSet(1, "test", true)
	assert.Equal(t, m.Len(), 3, "Len failed for NewMapSet(args)")

	m.Add("123", 77)
	assert.Equal(t, m.Len(), 5, "Len failed for Add()")

	m.Remove(1, "test", true, "123", 77)
	assert.Equal(t, m.Len(), 0, "Len failed for Remove()")
}

func TestMapSetSubset(t *testing.T) {
	m1 := set.NewMapSet(1, 2, 3)
	m2 := set.NewMapSet(9, false, 2, "true", 1, 3)
	m3 := set.NewMapSet(4, 1, 2)
	empty := set.NewMapSet()

	if !empty.Subset(m1) {
		t.Error("Subset failed: the empty set is always a valid subset")
	}
	if !m1.Subset(m2) {
		t.Errorf("Subset failed: %v is a subset of %v", m1, m2)
	}
	if m1.Subset(m3) {
		t.Errorf("Subset failed: %v is not a subset of %v", m1, m3)
	}
	if !empty.Subset(empty) {
		t.Error("Subset failed: the empty set has itself as subset")
	}
}

func TestMapSetSuperset(t *testing.T) {
	m1 := set.NewMapSet(1, 2, 3)
	m2 := set.NewMapSet(9, false, 2, "true", 1, 3)
	m3 := set.NewMapSet(4, 1, 2)
	empty := set.NewMapSet()

	if !m1.Superset(empty) {
		t.Error("Superset failed: the empty set is always a valid subset")
	}
	if !m2.Superset(m1) {
		t.Errorf("Superset failed: %v is a subset of %v", m1, m2)
	}
	if m3.Superset(m1) {
		t.Errorf("Superset failed: %v is not a subset of %v", m1, m3)
	}
	if !empty.Superset(empty) {
		t.Error("Superset failed: the empty set has itself as subset")
	}
}

func TestMapSetDisjoint(t *testing.T) {
	m1 := set.NewMapSet(1, 2, 3)
	m2 := set.NewMapSet(9, false, 2, "true", 1, 3)
	m3 := set.NewMapSet(4, "banana", 7)
	empty := set.NewMapSet()

	if m1.Disjoint(m2) {
		t.Errorf("Disjoint failed %v : %v share elements", m1, m2)
	}
	if m2.Disjoint(m1) {
		t.Errorf("Disjoint failed %v : %v share elements", m2, m1)
	}
	if !m3.Disjoint(m1) {
		t.Errorf("Disjoint failed %v : %v do not share elements", m3, m1)
	}
	if !m1.Disjoint(empty) {
		t.Errorf("Disjoint failed: %v contains the empty set", m1)
	}
	if !empty.Disjoint(empty) {
		t.Errorf("Disjoint failed: the empty set is disjoint with itself")
	}
	if !empty.Disjoint(m1) {
		t.Errorf("Disjoint failed: the empty set is disjoint with other sets")
	}

}

func TestMapSetIntersection(t *testing.T) {
	m1 := set.NewMapSet("qq", 2, 3, 4)
	m2 := set.NewMapSet(9, false, 2, "true", 1, 3)
	m3 := set.NewMapSet(4, "banana", 7)
	empty := set.NewMapSet()

	i1 := m1.Intersection(m2)
	if i1.Len() != 2 || !i1.Contains(2, 3) {
		t.Errorf("Intersection failed for %v and %v", m1, m2)
	}

	i2 := m1.Intersection(m2, m3)
	if i2.Len() != 0 {
		t.Errorf("Intersection failed for %v, %v and %v", m1, m2, m3)
	}

	i3 := m1.Intersection(m1)
	if i3.Len() != m1.Len() {
		t.Error("Intersection failed for a set with itself")
	}
	if !i3.Contains(m1.Items()...) {
		t.Error("Intersection failed for a set with itself")
	}

	i4 := m1.Intersection(empty)
	if i4.Len() != 0 {
		t.Error("Intersection failed for the empty set")
	}
}

func TestMapSetDifference(t *testing.T) {
	m1 := set.NewMapSet("qq", 2, 3)
	m2 := set.NewMapSet(9, false, 2, "true", 1, 3)
	m3 := set.NewMapSet(4, "banana", 7, "true")
	empty := set.NewMapSet()

	d1 := m1.Difference(m2)
	if d1.Len() != 1 || !d1.Contains("qq") {
		t.Errorf("Difference failed for %v and %v: %v", m1, m2, d1)
	}
	d2 := m1.Difference(m1)
	if d2.Len() != 0 {
		t.Error("Difference failed for a set with itself")
	}
	d3 := m1.Difference(m3)
	if d3.Len() != 3 || !d3.Contains(m1.Items()...) {
		t.Errorf("Difference failed for %v and %v: %v", m1, m3, d3)
	}
	d4 := m1.Difference(empty)
	if d4.Len() != 3 || !d3.Contains(m1.Items()...) {
		t.Error("Difference failed for a set with the empty set")
	}
	d5 := empty.Difference(m1)
	if d5.Len() != 0 {
		t.Error("Difference failed for the empty set")
	}
	d6 := m2.Difference(m1, m3)
	if d6.Len() != 3 || !d6.Contains(9, false, 1) {
		t.Errorf("Difference failed for %v %v, %v: %v", m2, m1, m3, d6)
	}
}

func TestMapSetString(t *testing.T) {
	m := set.NewMapSet()
	if m.String() != "MapSet{}" {
		t.Errorf("String: %v != 'MapSet{}'", m)
	}
	m2 := set.NewMapSet(1, 2, "qq")
	if !strings.HasPrefix(m2.String(), "MapSet{") {
		t.Errorf("String: %v does not start with 'MapSet{'", m2)
	}
	if !strings.HasSuffix(m2.String(), "}") {
		t.Errorf("String: %v does not end with '}'", m2)
	}
	if !strings.Contains(m2.String(), "\"qq\"") {
		t.Errorf("String: %v does not contain 'qq'", m2)
	}
}
