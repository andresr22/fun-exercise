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

func main() {
	reader := bufio.NewReader(os.Stdin)

	s := stores{
		storeInt:    make(map[string]int),
		storeString: make(map[string]string),
	}

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
		key := is[1]

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
		}
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

	if (len(is) >= 4) || (is[0] != "SET" && len(is) >= 3) {
		return errors.New("too many arguments")
	} else if _, ok := actions[is[0]]; !ok {
		return errors.New("invalid action")
	} else if (len(is) <= 1) || (is[0] == "SET" && len(is) <= 2) {
		return errors.New("missing arguments")
	}

	return nil
}

func (s stores) getValues(k string) (string, error) {
	if _, ok := s.storeInt[k]; ok {
		return fmt.Sprint(s.storeInt[k]), nil
	} else if _, ok := s.storeString[k]; ok {
		return s.storeString[k], nil
	} else {
		return "", errors.New(k + " does not exist")
	}
}

func (s stores) increments(k string) {
	_, okInt := s.storeInt[k]
	_, okString := s.storeString[k]
	if !okInt && !okString {
		s.storeInt[k] = 1
		fmt.Println(s.storeInt[k])
	} else if okInt && !okString {
		s.storeInt[k] += 1
		fmt.Println(s.storeInt[k])
	} else {
		fmt.Println("value of", k, "is string", "[", s.storeString[k], "]")
	}
}

func (s stores) decrements(k string) {
	_, okInt := s.storeInt[k]
	_, okString := s.storeString[k]
	if !okInt && !okString {
		s.storeInt[k] = -1
		fmt.Println(s.storeInt[k])
	} else if okInt && !okString {
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
	_, okInt := s.storeInt[k]
	_, okString := s.storeString[k]
	if !okInt && !okString {
		fmt.Println(k, "does not exist")
	} else {
		delete(s.storeInt, k)
		delete(s.storeString, k)
		fmt.Println("ok")
	}
}
