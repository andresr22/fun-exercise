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

func TestSet(t *testing.T) {

	store := new(Store)
	store.Push(make(Data))
	store.Push(make(Data))
	tests := []struct {
		key   string
		value string
	}{
		{"string", "value"},
		{"int", "1"},
	}

	for i, tt := range tests {
		store.Set(tt.key, tt.value)
		switch store.head.data[tt.key].(type) {
		case string:
			assert.Equal(t, "value", store.head.data[tt.key], "test[%d]: {%s: %s}", i, tt.key, tt.value)
		case int:
			assert.Equal(t, 1, store.head.data[tt.key], "test[%d]: {%s: %s}", i, tt.key, tt.value)
		}
	}
}

func TestGet(t *testing.T) {
	store := new(Store)
	store.Push(make(Data))
	store.Push(make(Data))
	store.Set("test1", "1")
	tests := []struct {
		s         *Store
		key       string
		hasError  bool
		errorKind error
	}{
		{store, "test1", false, nil},
		{store, "key", true, DoesNotExist},
	}

	for i, tt := range tests {
		err := store.Get(tt.key)
		if tt.hasError {
			assert.Error(t, err, "test[%d]: %s", i, tt.key)
			assert.IsType(t, tt.errorKind, err)
		} else {
			assert.True(t, KeyExists(tt.s.head.data, tt.key), "test[%d]: %s", i, tt.key)
		}
	}
}

func TestDelete(t *testing.T) {
	store := new(Store)
	store.Push(make(Data))
	store.Push(make(Data))
	store.Set("test1", "1")
	tests := []struct {
		s         *Store
		key       string
		hasError  bool
		errorKind error
	}{
		{store, "test1", false, nil},
		{store, "key", true, DoesNotExist},
	}

	for i, tt := range tests {
		err := store.Delete(tt.key)
		if tt.hasError {
			assert.Error(t, err, "test[%d]: %s", i, tt.key)
			assert.IsType(t, tt.errorKind, err)
		} else {
			assert.False(t, KeyExists(tt.s.head.data, tt.key), "test[%d]: %s", i, tt.key)
		}
	}
}

func TestIncr(t *testing.T) {
	store := new(Store)
	store.Push(make(Data))
	store.Push(make(Data))
	store.Set("test1", "1")
	store.Set("test2", "test")
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
		err := store.Incr(tt.key)
		if tt.hasError {
			assert.Error(t, err, "test[%d]: %s", i, tt.key)
			assert.IsType(t, tt.errorKind, err)
		} else {
			assert.Equal(t, store.head.data[tt.key], tt.value, "test[%d]: %s", i, tt.key)
		}
	}
}

func TestDecr(t *testing.T) {
	store := new(Store)
	store.Push(make(Data))
	store.Push(make(Data))
	store.Set("test1", "1")
	store.Set("test2", "test")
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
		err := store.Decr(tt.key)
		if tt.hasError {
			assert.Error(t, err, "test[%d]: %s", i, tt.key)
			assert.IsType(t, tt.errorKind, err)
		} else {
			assert.Equal(t, store.head.data[tt.key], tt.value, "test[%d]: %s", i, tt.key)
		}
	}
}
