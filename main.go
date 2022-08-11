package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var storeInt = make(map[string]int)

// var storeString = make(map[string]string)

func main() {
	reader := bufio.NewReader(os.Stdin)

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
			increments(key)
		case "DECR":
			decrements(key)
		case "GET":
			i, err := get(key)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(i)
			}
		case "SET":
			value := is[2]
			i, _ := strconv.Atoi(value)
			sets(key, i)
		}
	}
}

func validateInput(is []string) error {
	actions := map[string]int{
		"INCR": 1,
		"DECR": 2,
		"GET":  3,
		"SET":  4,
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

func get(k string) (int, error) {
	if _, ok := storeInt[k]; !ok {
		return 0, errors.New(k + " does not exist")
	}

	return storeInt[k], nil
}

func increments(k string) {
	if _, ok := storeInt[k]; !ok {
		storeInt[k] = 1
		fmt.Println(storeInt[k])
	} else {
		storeInt[k] += 1
		fmt.Println(storeInt[k])
	}
}

func decrements(k string) {
	if _, ok := storeInt[k]; !ok {
		storeInt[k] = -1
		fmt.Println(storeInt[k])
	} else {
		storeInt[k] -= 1
		fmt.Println(storeInt[k])
	}
}

func sets(k string, v int) {
	storeInt[k] = v
	fmt.Println(storeInt[k])
}
