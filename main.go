package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type stores struct {
	storeInt    map[string]int
	storeString map[string]string
}

type transaction struct {
	begin bool
	stores
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	s := stores{
		storeInt:    make(map[string]int),
		storeString: make(map[string]string),
	}
	storesPointer := &s

	t := transaction{
		begin: false,
		stores: stores{
			storeInt:    make(map[string]int),
			storeString: make(map[string]string),
		},
	}
	transactionPointer := &t

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")
		if input == "exit" {
			os.Exit(0)
		}
		is := strings.Split(input, " ")
		err := validateInput(is)

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
			s.increments(key)
		case "DECR":
			s.decrements(key)
		case "GET":
			s, err := s.getValues(key)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(s)
			}
		case "SET":
			value := is[2]
			s.sets(key, value)
		case "DELETE":
			s.delete(key)
		case "BEGIN":
			s.beginTransaction(transactionPointer)
		case "COMMIT":
			err = s.commitTransaction(transactionPointer)
			if err != nil {
				fmt.Println(err)
			}
		case "ROLLBACK":
			err = storesPointer.rollbackTransaction(transactionPointer)
			if err != nil {
				fmt.Println(err)
			}
		case "PRINT":
			t.printTransaction()
			s.printSores()
		}
	}
}

func stringInMap(m map[string]int, s string) bool {
	_, ok := m[s]
	return ok
}

func (pointerToTransaction *transaction) copyMaps(s stores) {
	(*pointerToTransaction).stores.storeInt = make(map[string]int)
	for k, v := range s.storeInt {
		(*pointerToTransaction).stores.storeInt[k] = v
	}

	(*pointerToTransaction).stores.storeString = make(map[string]string)
	for k, v := range s.storeString {
		(*pointerToTransaction).stores.storeString[k] = v
	}
}

func validateInput(is []string) error {

	actions := map[string]int{
		"INCR":   1,
		"DECR":   2,
		"GET":    3,
		"SET":    4,
		"DELETE": 5,
	}

	transactions := map[string]int{
		"BEGIN":    1,
		"COMMIT":   2,
		"ROLLBACK": 3,
		"PRINT":    4,
	}

	if (len(is) >= 4) || (is[0] != "SET" && len(is) >= 3) {
		return errors.New("too many arguments")
	} else if !stringInMap(actions, is[0]) && !stringInMap(transactions, is[0]) {
		return errors.New("invalid action")
	} else if ((len(is) <= 1) || (is[0] == "SET" && len(is) <= 2)) && !stringInMap(transactions, is[0]) {
		return errors.New("missing arguments")
	}

	return nil
}

func (s stores) getValues(k string) (string, error) {
	if stringInMap(s.storeInt, k) {
		return fmt.Sprint(s.storeInt[k]), nil
	} else if _, ok := s.storeString[k]; ok {
		return s.storeString[k], nil
	} else {
		return "", errors.New(k + " does not exist")
	}
}

func (s stores) increments(k string) {
	_, okString := s.storeString[k]
	if !stringInMap(s.storeInt, k) && !okString {
		s.storeInt[k] = 1
		fmt.Println(s.storeInt[k])
	} else if stringInMap(s.storeInt, k) && !okString {
		s.storeInt[k] += 1
		fmt.Println(s.storeInt[k])
	} else {
		fmt.Println("value of", k, "is string", "[", s.storeString[k], "]")
	}
}

func (s stores) decrements(k string) {
	_, okString := s.storeString[k]
	if !stringInMap(s.storeInt, k) && !okString {
		s.storeInt[k] = -1
		fmt.Println(s.storeInt[k])
	} else if stringInMap(s.storeInt, k) && !okString {
		s.storeInt[k] -= 1
		fmt.Println(s.storeInt[k])
	} else {
		fmt.Println("value of", k, "is string", "[", s.storeString[k], "]")
	}
}

func (s stores) sets(k string, v string) {
	i, err := strconv.Atoi(v)
	if err != nil {
		s.storeString[k] = v
		fmt.Println(s.storeString[k])
	} else {
		s.storeInt[k] = i
		fmt.Println(s.storeInt[k])
	}
}

func (s stores) delete(k string) {
	_, okString := s.storeString[k]
	if !stringInMap(s.storeInt, k) && !okString {
		fmt.Println(k, "does not exist")
	} else {
		delete(s.storeInt, k)
		delete(s.storeString, k)
		fmt.Println("ok")
	}
}

func (s stores) beginTransaction(pointerToTransaction *transaction) {
	(*pointerToTransaction).begin = true

	pointerToTransaction.copyMaps(s)
	fmt.Println("ok")
}

func (s stores) commitTransaction(pointerToTransaction *transaction) error {
	if !(*pointerToTransaction).begin {
		return errors.New("there is no transaction started")
	}

	pointerToTransaction.copyMaps(s)
	fmt.Println("ok")
	return nil
}

func (pointerToStores *stores) rollbackTransaction(pointerToTransaction *transaction) error {
	if !(*pointerToTransaction).begin {
		return errors.New("there is no transaction started")
	}

	(*pointerToStores).storeInt = make(map[string]int)
	for k, v := range (*pointerToTransaction).stores.storeInt {
		(*pointerToStores).storeInt[k] = v
	}

	(*pointerToTransaction).stores.storeString = make(map[string]string)
	for k, v := range (*pointerToTransaction).stores.storeString {
		(*pointerToStores).storeString[k] = v
	}

	fmt.Println("ok")
	return nil
}

func (t transaction) printTransaction() {
	fmt.Println(t)
}

func (s stores) printSores() {
	fmt.Println(s)
}
