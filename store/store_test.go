package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyStore(t *testing.T) {
	store := NewStore(nil)
	store.Set("a", "NEW A")
	res, err := store.Replace("ab{{.a}}bb")

	assert.Nil(t, err)
	assert.Equal(t, "abNEW Abb", res)
}

func TestFilledStore(t *testing.T) {
	store := NewStore(map[string]string{
		"b": "NEW B",
	})
	store.Set("a", "NEW A")
	res, err := store.Replace("ab{{.a}}b{{.b}}b")

	assert.Nil(t, err)
	assert.Equal(t, "abNEW AbNEW Bb", res)
}
