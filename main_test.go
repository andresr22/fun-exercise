package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyExists(t *testing.T) {
	m := map[string]interface{}{
		"test1": 1,
	}

	tests := []struct {
		key    string
		exists bool
	}{
		{"test", false},
		{"test1", true},
	}

	for i, tt := range tests {
		if tt.exists {
			assert.True(t, KeyExists(m, tt.key), "test[%d]: %s", i, tt.key)
		} else {
			assert.False(t, KeyExists(m, tt.key), "test[%d]: %s", i, tt.key)
		}
	}
}

func TestValidateInput(t *testing.T) {
	tests := []struct {
		input     string
		hasError  bool
		errorKind error
	}{
		{"INCR key", false, nil},
		{"INCR key value", true, TooManyArguments},
		{"DEC key", true, InvalidAction},
		{"SET key", true, MissingArguments},
	}

	for i, tt := range tests {
		err := ValidateInput(strings.Split(tt.input, " "))
		if tt.hasError {
			assert.Error(t, err, "test[%d]: %s", i, tt.input)
		} else {
			assert.NoError(t, err, "test[%d]: %s", i, tt.input)
			assert.IsType(t, tt.errorKind, err)
		}
	}
}

func TestCopyMap(t *testing.T) {

	tests := []struct {
		s        StorageInterface
		fromelem bool
	}{
		{
			StorageInterface{
				elem: map[string]interface{}{
					"test1": 1,
					"test2": "test",
				},
				transaction: TransactionInterface{
					begin: false,
					elem:  make(map[string]interface{}),
				},
			},
			true,
		},

		{
			StorageInterface{
				elem: make(map[string]interface{}),
				transaction: TransactionInterface{
					begin: false,
					elem: map[string]interface{}{
						"test1": 1,
						"test2": "test",
					},
				},
			},
			false,
		},
	}

	for i, tt := range tests {
		tt.s.CopyMap(tt.fromelem)
		assert.Equal(t, tt.s.elem, tt.s.transaction.elem, "test[%d]: elem(%s) -> transaction.elem(%s)", i, tt.s.elem, tt.s.transaction.elem)
	}
}

func TestSet(t *testing.T) {
	s := StorageInterface{
		elem: make(map[string]interface{}),
		transaction: TransactionInterface{
			begin: false,
			elem:  make(map[string]interface{}),
		},
	}
	tests := []struct {
		key   string
		value string
	}{
		{"string", "value"},
		{"int", "1"},
	}

	for i, tt := range tests {
		s.Set(tt.key, tt.value)
		switch s.elem[tt.key].(type) {
		case string:
			assert.Equal(t, "value", s.elem[tt.key], "test[%d]: {%s: %s}", i, tt.key, tt.value)
		case int:
			assert.Equal(t, 1, s.elem[tt.key], "test[%d]: {%s: %s}", i, tt.key, tt.value)
		}
	}
}
