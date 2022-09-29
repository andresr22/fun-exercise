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
			assert.IsType(t, tt.errorKind, err)
		} else {
			assert.NoError(t, err, "test[%d]: %s", i, tt.input)
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
				transaction: NewStack(),
			},
			true,
		},
		{
			StorageInterface{
				elem:        make(map[string]interface{}),
				transaction: NewStack(),
			},
			false,
		},
	}

	for i, tt := range tests {
		transaction := map[string]interface{}{"key": 1, "key1": "string"}
		tt.s.transaction.Push(&Node{transaction})
		switch tt.fromelem {
		case true:
			tt.s.CopyMap(tt.fromelem)
			transaction = tt.s.transaction.Pop().value
		default:
			tt.s.CopyMap(tt.fromelem)
		}
		assert.Equal(t, tt.s.elem, transaction, "test[%d]: elem(%s) -> transaction.elem(%s)", i, tt.s.elem, transaction)
	}
}

func TestSet(t *testing.T) {
	s := StorageInterface{
		elem:        make(map[string]interface{}),
		transaction: NewStack(),
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

func TestGet(t *testing.T) {
	s := StorageInterface{
		elem: map[string]interface{}{
			"test1": 1,
			"test2": "test",
		},
		transaction: NewStack(),
	}
	tests := []struct {
		s         StorageInterface
		key       string
		hasError  bool
		errorKind error
	}{
		{s, "test1", false, nil},
		{s, "key", true, DoesNotExist},
	}

	for i, tt := range tests {
		err := s.Get(tt.key)
		if tt.hasError {
			assert.Error(t, err, "test[%d]: %s", i, tt.key)
			assert.IsType(t, tt.errorKind, err)
		} else {
			assert.True(t, KeyExists(tt.s.elem, tt.key), "test[%d]: %s", i, tt.key)
		}
	}
}

func TestDelete(t *testing.T) {
	s := StorageInterface{
		elem: map[string]interface{}{
			"test1": 1,
			"test2": "test",
		},
		transaction: NewStack(),
	}
	tests := []struct {
		s         StorageInterface
		key       string
		hasError  bool
		errorKind error
	}{
		{s, "test1", false, nil},
		{s, "key", true, DoesNotExist},
	}

	for i, tt := range tests {
		err := s.Delete(tt.key)
		if tt.hasError {
			assert.Error(t, err, "test[%d]: %s", i, tt.key)
			assert.IsType(t, tt.errorKind, err)
		} else {
			assert.False(t, KeyExists(tt.s.elem, tt.key), "test[%d]: %s", i, tt.key)
		}
	}
}

func TestIncr(t *testing.T) {
	s := StorageInterface{
		elem: map[string]interface{}{
			"test1": 1,
			"test2": "test",
		},
		transaction: NewStack(),
	}
	tests := []struct {
		key       string
		value     int
		hasError  bool
		errorKind error
	}{
		{"test1", 2, false, nil},
		{"test2", 0, true, CannotIncrStr},
	}

	for i, tt := range tests {
		err := s.Incr(tt.key)
		if tt.hasError {
			assert.Error(t, err, "test[%d]: %s", i, tt.key)
			assert.IsType(t, tt.errorKind, err)
		} else {
			assert.Equal(t, s.elem[tt.key], tt.value, "test[%d]: %s", i, tt.key)
		}
	}
}

func TestDecr(t *testing.T) {
	s := StorageInterface{
		elem: map[string]interface{}{
			"test1": 1,
			"test2": "test",
		},
		transaction: NewStack(),
	}
	tests := []struct {
		key       string
		value     int
		hasError  bool
		errorKind error
	}{
		{"test1", 0, false, nil},
		{"test2", 0, true, CannotIncrStr},
	}

	for i, tt := range tests {
		err := s.Decr(tt.key)
		if tt.hasError {
			assert.Error(t, err, "test[%d]: %s", i, tt.key)
			assert.IsType(t, tt.errorKind, err)
		} else {
			assert.Equal(t, s.elem[tt.key], tt.value, "test[%d]: %s", i, tt.key)
		}
	}
}
