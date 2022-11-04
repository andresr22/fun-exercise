package main

import (
	"fmt"
	"strings"
)

type Data map[string]interface{}

type Element struct {
	data Data
	prev *Element
	next *Element
}

type Store struct {
	head *Element
	tail *Element
}

func (s *Store) First() *Element {
	return s.head
}

func (e *Element) Next() *Element {
	return e.next
}

func (e *Element) Prev() *Element {
	return e.prev
}

func (s *Store) Push(d Data) *Store {
	elem := &Element{data: d}
	if s.head == nil {
		s.head = elem
	} else {
		s.tail.next = elem
		elem.prev = s.tail
	}
	s.tail = elem
	return s
}

func (s *Store) Pop() (d Data, err error) {
	if s.tail == nil {
		err = NoTransactionStarted
	} else {
		d = s.tail.data
		if s.tail == s.head {
			err = NoTransactionStarted
		} else {
			s.tail = s.tail.prev
		}
	}
	return d, err
}

func (s Store) Size() int {
	size := 0
	for n := s.First(); n != nil; n = n.Next() {
		size += 1
	}
	return size
}

func (s Store) Print() {
	dashes := strings.Repeat("-", 50)
	fmt.Println(dashes)
	fmt.Println("store head: ", s.head)
	fmt.Println("store tail: ", s.tail)
	fmt.Println(dashes)
	for n := s.First(); n != nil; n = n.Next() {
		fmt.Printf("%v\n", n.data)
	}
	fmt.Println(dashes)
	fmt.Println("Size: ", s.Size())
}
