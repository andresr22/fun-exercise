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
		err = StoreEmpty
	} else {
		d = s.tail.data
		s.tail = s.tail.prev
		if s.tail == nil {
			s.head = nil
		}
	}
	return d, err
}

func (s Store) Print() {
	dashes := strings.Repeat("-", 50)
	fmt.Println(dashes)
	for n := s.First(); n != nil; n = n.Next() {
		fmt.Printf("%v\n", n.data)
	}
}

// func main() {
// 	dashes := strings.Repeat("-", 50)
// 	s := new(Store)

// 	d := make(Data)
// 	d["key"] = 1
// 	s.Push(d)

// 	d1 := make(Data)
// 	d1["hello"] = "world"
// 	s.Push(d1)

// 	d2 := make(Data)
// 	d2["anotherkey"] = 2
// 	s.Push(d2)

// 	processed := make(map[*Element]bool)

// 	fmt.Println("First time through Store...")
// 	for n := s.First(); n != nil; n = n.Next() {
// 		fmt.Printf("%v\n", n.data)
// 		if processed[n] {
// 			fmt.Printf("%s as been processed\n", n.data)
// 		}
// 		processed[n] = true
// 	}
// 	fmt.Println(dashes)
// 	fmt.Println("Second time through Store...")
// 	for n := s.First(); n != nil; n = n.Next() {
// 		fmt.Printf("%v", n.data)
// 		if processed[n] {
// 			fmt.Println(" has been processed")
// 		} else {
// 			fmt.Println()
// 		}
// 		processed[n] = true
// 	}

// 	fmt.Println(dashes)
// 	fmt.Println("* Pop each value off Store...")
// 	for v, err := s.Pop(); err == nil; v, err = s.Pop() {
// 		fmt.Printf("%v\n", v)
// 	}
// 	fmt.Println(s.Pop())
// }
