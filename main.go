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

type StorageInterface struct {
	elem        map[string]interface{}
	transaction *Stack
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	s := StorageInterface{
		elem:        make(map[string]interface{}),
		transaction: NewStack(),
	}

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
			err = s.Incr(key)
		case "DECR":
			err = s.Decr(key)
		case "GET":
			err = s.Get(key)
		case "SET":
			value := is[2]
			err = s.Set(key, value)
		case "DELETE":
			err = s.Delete(key)
		case "BEGIN":
			err = s.Begin()
		case "COMMIT":
			err = s.Commit()
		case "ROLLBACK":
			err = s.Rollback()
		case "PRINT":
			fmt.Println("s.elem: ", s.elem)
			fmt.Print("s.transaction: ")
			s.transaction.Print()
			fmt.Println()
		}

		s.PrintResult(key, err)
	}
}

func (s StorageInterface) PrintResult(key string, err error) {
	switch err {
	case nil:
		switch KeyExists(s.elem, key) {
		case true:
			fmt.Println(s.elem[key])
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

func (s *StorageInterface) CopyMap(b bool) {
	// if b is true then enqueue an elem
	// otherwise dequeue and the last elem
	switch b {
	case true:
		m := make(map[string]interface{})
		for k, v := range s.elem {
			m[k] = v
		}
		s.transaction.Push(&Node{m})
	default:
		elem := s.transaction.Pop()
		s.elem = make(map[string]interface{})
		for k, v := range elem.value {
			s.elem[k] = v
		}
	}
}

func (s *StorageInterface) Set(key string, value interface{}) error {
	i, err := strconv.Atoi(fmt.Sprintf("%v", value))
	switch err {
	case nil:
		s.elem[key] = i
	default:
		s.elem[key] = value
	}
	return nil
}

func (s *StorageInterface) Get(key string) error {
	switch KeyExists(s.elem, key) {
	case true:
		return nil
	default:
		return DoesNotExist
	}
}

func (s *StorageInterface) Delete(key string) error {
	switch KeyExists(s.elem, key) {
	case true:
		delete(s.elem, key)
	default:
		return DoesNotExist
	}
	return nil
}

func (s *StorageInterface) Incr(key string) error {
	switch v := s.elem[key].(type) {
	case nil:
		s.elem[key] = 1
	case int:
		s.elem[key] = v + 1
	default:
		return CannotIncrStr
	}
	return nil
}

func (s *StorageInterface) Decr(key string) error {
	switch v := s.elem[key].(type) {
	case nil:
		s.elem[key] = -1
	case int:
		s.elem[key] = v - 1
	default:
		return CannotDecrStr
	}
	return nil
}

func (s *StorageInterface) Begin() error {
	s.CopyMap(true)
	return nil
}

func (s *StorageInterface) Commit() error {
	if s.transaction.size == 0 {
		return NoTransactionStarted
	}
	s.CopyMap(true)
	return nil
}

func (s *StorageInterface) Rollback() error {
	if s.transaction.size == 0 {
		return NoTransactionStarted
	}
	s.CopyMap(false)
	return nil
}

// func (s *StorageInterface) IsString(key string) bool {
// 	switch s.elem[key].(type) {
// 	case string:
// 		return true
// 	default:
// 		return false
// 	}
// }
