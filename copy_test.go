package reflectricity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrivateCopy(t *testing.T) {
	p := withPrivateField{
		Public: 10,
		foo:    "bar",
	}

	cpy := DeepCopy(p, true)
	assert.Equal(t, p, cpy)
	// WOOOO
	assert.Equal(t, p.foo, cpy.foo)
	cpy.foo = "baz"
	assert.NotEqual(t, p.foo, cpy.foo)
}

func TestDeepCopyLinkedList(t *testing.T) {
	n := &node{
		value: 1,
		next: &node{
			value: 2,
			next: &node{
				value: 3,
			},
		},
	}
	cpy := DeepCopy(n, true)
	assert.Equal(t, n, cpy)
	cpy.next.next = &node{
		value: 4,
	}
	assert.NotEqual(t, n, cpy)
}

func TestCopyChannel(t *testing.T) {
	c1 := make(chan string, 10)
	c2 := DeepCopy(c1, true)
	assert.Equal(t, cap(c2), 10)
}

type node struct {
	value int
	next  *node
}

func TestCopyIface(t *testing.T) {
	m := map[string]cpIfaceTest{
		"foo": new(testS),
		"bar": new(testX),
	}
	x := DeepCopy(m, false)
	assert.Equal(t, m, x)
}

type testS struct{}

func (t *testS) Tst() {}

type testX struct{}

func (t *testX) Tst() {}

type cpIfaceTest interface {
	Tst()
}
