package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirst(t *testing.T) {
	store := new(Store)
	store.Push(Data{"key1": 1})
	store.Push(Data{"key2": 2})
	store.Push(Data{"key3": 3})

	tests := []struct {
		key   string
		value interface{}
	}{
		{"key1", 1},
		{"key2", 2},
	}

	for i, tt := range tests {
		first := store.First()
		fmt.Println(first.data[tt.key])
		assert.Equal(t, tt.value, first.data[tt.key], "test[%d]: %s", i, tt.key)
	}

}
