package reflectricity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	privateReflector = NewReflector(true)
	publicReflector  = NewReflector(false)
)

func TestPrintProperArrayType(t *testing.T) {
	i := &invites{
		Invites: []*invite{
			{
				To:   "baz",
				CC:   "Does BAZ",
				Text: "Hello",
			},
		},
	}

	expectation := `&reflectricity.invites{
	Invites: []*reflectricity.invite{
		&reflectricity.invite{
			To: "baz",
			CC: "Does BAZ",
			Text: "Hello",
		},
	},
}`
	out := publicReflector.PrintValue(i)
	assert.Equal(t, expectation, out)
}

func TestPrintInt8(t *testing.T) {
	i := int8(10)
	assert.Equal(t, "10", publicReflector.PrintValue(i))
}

func TestToPtr(t *testing.T) {
	p := new(***withPrivateField)
	wp := withPrivateField{
		Public: 10,
		foo:    "bar",
	}
	x := &wp
	y := &x
	z := &y
	p = &z
	s := privateReflector.PrintValue(p)
	expectation := `reflectricity.ToPtr(reflectricity.ToPtr(reflectricity.ToPtr(&reflectricity.withPrivateField{
	Public: 10,
	foo: "bar",
})))`
	second := ToPtr(ToPtr(ToPtr(&withPrivateField{
		Public: 10,
		foo:    "bar",
	})))
	assert.Equal(t, expectation, s)
	assert.Equal(t, expectation, privateReflector.PrintValue(second))
}

func TestPrintPrivateField(t *testing.T) {
	expectation := `reflectricity.withPrivateField{
	Public: 1,
	foo: "bar",
}`

	i := withPrivateField{
		Public: 1,
		foo:    "bar",
	}

	r := privateReflector
	result := r.PrintValue(i)
	assert.Equal(t, expectation, result)
}

func TestPrint(t *testing.T) {
	result := publicReflector.PrintValue(ts)
	expectation := `reflectricity.user{
	Username: "foo",
	Email: "bar@google.com",
	Friends: 0,
	Profile: &reflectricity.profile{
		Age: 30,
	},
	Invites: []reflectricity.invite{
		reflectricity.invite{
			To: "baz",
			CC: "Does BAZ",
			Text: "Hello",
		},
	},
	Invites2: []*reflectricity.invite{
		&reflectricity.invite{
			To: "baz",
			CC: "Does BAZ",
			Text: "Hello",
		},
	},
	StoredLocations: map[string]string{
		"home": "123 elm street",
	},
}`
	assert.Equal(t, expectation, result)
}

func TestCustomStringType(t *testing.T) {
	f := fString("hello")
	assert.Equal(t, "reflectricity.fString(\"hello\")", publicReflector.PrintValue(f))
}

func TestCustomBoolType(t *testing.T) {
	f := fbool(false)
	assert.Equal(t, "reflectricity.fbool(false)", publicReflector.PrintValue(f))
}

func TestCustomMapType(t *testing.T) {
	m := fmap(map[string]string{
		"foo": "bar",
	})
	assert.Equal(t, "reflectricity.fmap(map[string]string{\n\t\"foo\": \"bar\",\n})", publicReflector.PrintValue(m))
}

func TestPrintChannel(t *testing.T) {
	ch := make(chan fString)
	assert.Equal(t, "make(chan reflectricity.fString)", publicReflector.PrintValue(ch))
}

func TestPrintWithAlias(t *testing.T) {
	r := NewReflector(true)
	r.AddAlias("github.com/MichaelCombs28/reflectricity", "ref")
	result := r.PrintValue(map[string]*profile{
		"foo": {},
	})
	expectation := `map[string]*ref.profile{
	"foo": &ref.profile{
		Age: 0,
	},
}`
	assert.Equal(t, expectation, result)
}

func TestWrapBool(t *testing.T) {
	b := ToPtr(ToPtr(false))
	result := publicReflector.PrintValue(b)
	assert.Equal(t, "reflectricity.ToPtr(reflectricity.ToPtr(false))", result)
}

type fString string

type fbool bool

type fmap map[string]string

type invites struct {
	Invites []*invite
}

type user struct {
	Username        string
	Email           string
	Friends         int
	Profile         *profile
	Invites         []invite
	Invites2        []*invite
	StoredLocations map[string]string
}

type profile struct {
	Age int
}

type invite struct {
	To   string
	CC   string
	Text string
}

type withPrivateField struct {
	Public int
	foo    string
}

var ts = user{
	Username: "foo",
	Email:    "bar@google.com",
	Profile: &profile{
		Age: 30,
	},
	Invites: []invite{
		{
			To:   "baz",
			CC:   "Does BAZ",
			Text: "Hello",
		},
	},
	Invites2: []*invite{
		{
			To:   "baz",
			CC:   "Does BAZ",
			Text: "Hello",
		},
	},
	StoredLocations: map[string]string{
		"home": "123 elm street",
	},
}
