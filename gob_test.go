package reflectricity

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGobRegister(t *testing.T) {
	mp := map[string]Test{
		"Test": {
			R:       &R1{32},
			G:       new(G2),
			B:       B1{},
			Opacity: 0.3,
		},
	}

	RegisterGob(mp)
	buf := new(bytes.Buffer)
	assert.NoError(t, gob.NewEncoder(buf).Encode(mp))
	fmt.Println(string(buf.Bytes()))
	reader := bytes.NewReader(buf.Bytes())
	var m map[string]Test
	assert.NoError(t, gob.NewDecoder(reader).Decode(&m))
	assert.Equal(t, mp, m)
}

func TestGobRegisterArray(t *testing.T) {
	mp := []Test{
		{
			R:       &R1{32},
			G:       new(G2),
			B:       B1{},
			Opacity: 0.3,
		},
	}

	RegisterGob(mp)
	buf := new(bytes.Buffer)
	assert.NoError(t, gob.NewEncoder(buf).Encode(mp))
	fmt.Println(string(buf.Bytes()))
	reader := bytes.NewReader(buf.Bytes())
	var m []Test
	assert.NoError(t, gob.NewDecoder(reader).Decode(&m))
	assert.Equal(t, mp, m)
}

type Test struct {
	R
	G
	B
	Opacity float64
}

type R interface {
	r()
}

type R1 struct {
	Value float32
}

func (r *R1) r() {
}

type R2 struct {
}

func (r *R2) r() {
}

type G interface {
	g()
}

type G1 struct {
}

func (g *G1) g() {
}

type G2 struct {
}

func (g *G2) g() {
}

type B interface {
	b()
}

type B1 struct {
}

func (b B1) b() {
}
