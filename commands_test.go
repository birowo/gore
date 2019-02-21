package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActionDoc(t *testing.T) {
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)

	s, err := NewSession(stdout, stderr)
	defer s.Clear()
	require.NoError(t, err)

	err = actionImport(s, "encoding/json")
	require.NoError(t, err)
	err = actionImport(s, "fmt")
	require.NoError(t, err)

	test := func() {
		err = actionDoc(s, "fmt")
		require.NoError(t, err)

		err = actionDoc(s, "fmt.Print")
		require.NoError(t, err)

		err = actionDoc(s, "json.NewEncoder(nil).Encode")
		require.NoError(t, err)
	}

	test()

	// test :doc works after some code

	s.Eval("a := 1")
	s.Eval("fmt.Print()")

	test()

	assert.Contains(t, stdout.String(), "package fmt")
	assert.Contains(t, stdout.String(), "func Printf")
	assert.Equal(t, "", stderr.String())
}

func TestActionImport(t *testing.T) {
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	s, err := NewSession(stdout, stderr)
	defer s.Clear()
	require.NoError(t, err)

	require.NoError(t, actionImport(s, "encoding/json fmt"))

	require.NoError(t, s.Eval("fmt.Print"))
	require.NoError(t, s.Eval("json.Encoder{}"))

	assert.Contains(t, stdout.String(), "(func(...interface {}) (int, error))")
	assert.Contains(t, stdout.String(), "json.Encoder")
	assert.Equal(t, "", stderr.String())
}
