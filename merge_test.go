package reflectricity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeStruct(t *testing.T) {
	t1 := testStruct{
		PublicField: &testStruct{
			Foo: 2,
		},
		Foo: 1,
		Bar: 2,
		Baz: s("bar"),
	}
	t2 := testStruct{
		PublicField: &testStruct{
			Foo: 1,
		},
		Foo: 3,
		Bar: 4,
		Baz: s("buz"),
	}

	result := MergeLeftWithOptions(t1, t2, false, CONCAT)
	assert.Equal(t, testStruct{
		PublicField: &testStruct{
			Foo: 1,
		},
		Foo: 3,
		Bar: 4,
		Baz: s("buz"),
	}, result)
}

func TestMergeStructWithPrivate(t *testing.T) {
	t1 := testStruct{
		PublicField: &testStruct{
			Foo: 2,
		},
		Foo:          1,
		Bar:          2,
		Baz:          s("bar"),
		privateField: s("foo"),
	}
	t2 := testStruct{
		PublicField: &testStruct{
			Foo: 1,
		},
		Foo:          3,
		Bar:          4,
		Baz:          s("buz"),
		privateField: s("bar"),
	}

	result := MergeLeftWithOptions(t1, t2, true, CONCAT)
	assert.Equal(t, testStruct{
		PublicField: &testStruct{
			Foo: 1,
		},
		Foo:          3,
		Bar:          4,
		Baz:          s("buz"),
		privateField: s("bar"),
	}, result)
}

func TestMergeStructWithPrivateNil(t *testing.T) {
	t1 := testStruct{
		PublicField: &testStruct{
			Foo: 2,
		},
		Foo:          1,
		Bar:          2,
		Baz:          s("bar"),
		privateField: s("foo"),
	}
	t2 := testStruct{
		PublicField: &testStruct{
			Foo: 1,
		},
		Foo: 3,
		Bar: 4,
	}

	result := MergeLeftWithOptions(t1, t2, true, CONCAT)
	assert.Equal(t, testStruct{
		PublicField: &testStruct{
			Foo: 1,
		},
		Foo:          3,
		Bar:          4,
		Baz:          s("bar"),
		privateField: s("foo"),
	}, result)
}

func TestMergeArrayConcat(t *testing.T) {
	i1 := []int{1, 2, 3}
	i2 := []int{4, 5, 6, 7}

	arr := MergeLeftWithOptions(i1, i2, false, CONCAT)
	a, ok := arr.([]int)
	assert.True(t, ok)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, a)
}

func TestMergeArrayFullReplace(t *testing.T) {
	i1 := []int{1, 2, 3}
	i2 := []int{4, 5, 6, 7}

	arr := MergeLeftWithOptions(i1, i2, false, REPLACE)
	a, ok := arr.([]int)
	assert.True(t, ok)
	assert.Equal(t, []int{4, 5, 6, 7}, a)
}

func TestMergeMap(t *testing.T) {
	m1 := map[string]string{
		"foo": "bar",
	}
	m2 := map[string]string{
		"baz": "buz",
	}

	result := MergeLeft(m1, m2)
	assert.Equal(t, map[string]string{
		"foo": "bar",
		"baz": "buz",
	}, result)
}

type testStruct struct {
	PublicField  *testStruct
	Foo          int
	Bar          int
	Baz          *string
	privateField *string
}

func s(st string) *string {
	return &st
}
