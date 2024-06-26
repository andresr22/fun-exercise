package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CacheError string

const (
	TooManyArguments     CacheError = "too many arguments"
	InvalidAction        CacheError = "invalid action"
	MissingArguments     CacheError = "missing arguments"
	CannotIncrStr        CacheError = "cannot increment string"
	CannotDecrStr        CacheError = "cannor decrement string"
	DoesNotExist         CacheError = "key does not exist"
	NoTransactionStarted CacheError = "there is no transaction started"
)

var actions = map[string]interface{}{
	"INCR":   1,
	"DECR":   2,
	"GET":    3,
	"SET":    4,
	"DELETE": 5,
}

var transactions = map[string]interface{}{
	"BEGIN":    1,
	"COMMIT":   2,
	"ROLLBACK": 3,
	"PRINT":    4,
}

func (e CacheError) Error() string {
	return string(e)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	store := new(Store)
	store.Push(make(Data))
	store.Push(make(Data))

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")
		if input == "exit" {
			os.Exit(0)
		}

		is := strings.Split(input, " ")
		err := ValidateInput(is)

		if err != nil {
			fmt.Println(err)
			continue
		}

		action := is[0]
		key := ""
		if len(is) > 1 {
			key = is[1]
		}

		switch action {
		case "INCR":
			err = store.Incr(key)
		case "DECR":
			err = store.Decr(key)
		case "GET":
			err = store.Get(key)
		case "SET":
			value := is[2]
			err = store.Set(key, value)
		case "DELETE":
			err = store.Delete(key)
		case "BEGIN":
			err = store.Begin()
		case "COMMIT":
			err = store.Commit()
		case "ROLLBACK":
			err = store.Rollback()
		case "PRINT":
			store.Print()
			// fmt.Println()
		}

		store.PrintResult(key, err)
	}
}

func (s Store) PrintResult(key string, err error) {
	switch err {
	case nil:
		switch KeyExists(s.head.data, key) {
		case true:
			fmt.Println(s.head.data[key])
		default:
			fmt.Println("ok")
		}
	default:
		fmt.Println(err)
	}
}

func KeyExists(m map[string]interface{}, s string) bool {
	_, ok := m[s]
	return ok
}

func ValidateInput(s []string) error {
	if (len(s) >= 4) || (s[0] != "SET" && len(s) >= 3) {
		return TooManyArguments
	} else if !KeyExists(actions, s[0]) && !KeyExists(transactions, s[0]) {
		return InvalidAction
	} else if ((len(s) <= 1) || (s[0] == "SET" && len(s) <= 2)) && !KeyExists(transactions, s[0]) {
		return MissingArguments
	}
	return nil
}

func (s *Store) Set(key string, value interface{}) error {
	i, err := strconv.Atoi(fmt.Sprintf("%v", value))
	switch err {
	case nil:
		s.head.data[key] = i
	default:
		s.head.data[key] = value
	}
	s.tail.data[key] = s.head.data[key]
	return nil
}

func (s *Store) Get(key string) error {
	switch KeyExists(s.head.data, key) {
	case true:
		return nil
	default:
		return DoesNotExist
	}
}

func (s *Store) Delete(key string) error {
	switch KeyExists(s.head.data, key) {
	case true:
		delete(s.tail.data, key)
		delete(s.head.data, key)
	default:
		return DoesNotExist
	}
	return nil
}

func (s *Store) Incr(key string) error {
	switch v := s.head.data[key].(type) {
	case nil:
		s.head.data[key] = 1
	case int:
		s.head.data[key] = v + 1
	default:
		return CannotIncrStr
	}
	s.tail.data[key] = s.head.data[key]
	return nil
}

func (s *Store) Decr(key string) error {
	switch v := s.head.data[key].(type) {
	case nil:
		s.head.data[key] = -1
	case int:
		s.head.data[key] = v - 1
	default:
		return CannotDecrStr
	}
	s.tail.data[key] = s.head.data[key]
	return nil
}

func (s *Store) Begin() error {
	s.Push(make(Data))
	return nil
}

func (s *Store) Commit() error {
	if s.Size() <= 2 {
		return NoTransactionStarted
	}
	s.Push(make(Data))
	return nil
}

func (s *Store) Rollback() error {
	d, err := s.Pop()
	if err != nil {
		return err
	}
	for k, v := range s.tail.data {
		s.head.data[k] = v
	}
	for k := range d {
		if !KeyExists(s.tail.data, k) {
			delete(s.head.data, k)
		}
	}
	return nil
}
