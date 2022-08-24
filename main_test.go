package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestUtils(t *testing.T) {

	/*** stringInMap Tests ***/
	m := map[string]int{
		"test1": 1,
		"test2": 2,
	}

	// The key 'test' shouldn't exists
	if stringInMap(m, "test") {
		t.Errorf("[stringInMap shouldn't exists] Expected false, but got %v", stringInMap(m, "test"))
	}

	// The key 'test1' should exists
	if !stringInMap(m, "test1") {
		t.Errorf("[stringInMap should exists] Expected true, but got %v", stringInMap(m, "test1"))
	}
	/*** stringInMap Tests ***/

	/*** copyMaps Tests ***/
	s := stores{
		storeInt:    make(map[string]int),
		storeString: make(map[string]string),
	}

	transactions := transaction{
		begin: false,
		stores: stores{
			storeInt: map[string]int{
				"test": 1,
			},
			storeString: map[string]string{
				"test": "test",
			},
		},
	}
	transactionPointer := &transactions

	transactionPointer.copyMaps(s)

	if !reflect.DeepEqual(transactions.stores, s) {
		t.Errorf("[copyMaps] Expected two equal maps, but got map1 %v != map2 %v", transactions, s)
	}
	/*** copyMaps Tests ***/

	/*** validateInput Tests ***/
	err := validateInput(strings.Split("INCR key", " "))
	if err != nil {
		t.Errorf("[validateInput 'INCR key'] Expected nil, but got '%v'", err)
	}

	err = validateInput(strings.Split("DECR key", " "))
	if err != nil {
		t.Errorf("[validateInput 'DECR key'] Expected nil, but got '%v'", err)
	}

	err = validateInput(strings.Split("GET key", " "))
	if err != nil {
		t.Errorf("[validateInput 'GET key'] Expected nil, but got '%v'", err)
	}

	err = validateInput(strings.Split("SET key value", " "))
	if err != nil {
		t.Errorf("[validateInput 'SET key value'] Expected nil, but got '%v'", err)
	}

	err = validateInput(strings.Split("DELETE key", " "))
	if err != nil {
		t.Errorf("[validateInput 'DELETE key'] Expected nil, but got '%v'", err)
	}

	err = validateInput(strings.Split("INCR key value", " "))
	if err.Error() != "too many arguments" {
		t.Errorf("[validateInput 'INCR key value'] Expected 'too many arguments', but got '%v'", err)
	}

	err = validateInput(strings.Split("INCR", " "))
	if err.Error() != "missing arguments" {
		t.Errorf("[validateInput 'INCR'] Expected 'missing arguments', but got '%v'", err)
	}

	err = validateInput(strings.Split("SET key", " "))
	if err.Error() != "missing arguments" {
		t.Errorf("[validateInput 'SET key'] Expected 'missing arguments', but got '%v'", err)
	}

	err = validateInput(strings.Split("QWERTY", " "))
	if err.Error() != "invalid action" {
		t.Errorf("[validateInput 'QWERTY'] Expected 'invalid action', but got '%v'", err)
	}
	/*** validateInput Tests ***/

}

func TestActions(t *testing.T) {

	s := stores{
		storeInt: map[string]int{
			"key": 0,
		},
		storeString: map[string]string{
			"key_str": "str",
		},
	}

	/*** getValues Tests ***/
	n, err := s.getValues("key")
	if n != "0" {
		t.Errorf("[getValues 'key: 0'] Expected '0', but got '%v'", n)
	}
	if err != nil {
		t.Errorf("[getValues 'key: 0'] Expected nil, but got '%v'", err)
	}

	// need to check this test
	_, err = s.getValues("anotherkey")
	if err.Error() != "anotherkey does not exist" {
		t.Errorf("[getValues 'anotherkey shouldn't exists'] Expected 'anotherkey does not exist', but got '%v'", err)
	}

	str, err_str := s.getValues("key_str")
	if str != "str" {
		t.Errorf("[getValues 'key_str: str'] Expected 'str', but got '%v'", str)
	}
	if err_str != nil {
		t.Errorf("[getValues 'key_str: str'] Expected nil, but got '%v'", err_str)
	}
	/*** getValues Tests ***/

	/*** increments Tests ***/
	v := s.increments("key")
	if v != "1" {
		t.Errorf("[increments 'key'] Expected 1, but got '%v'", v)
	}

	v = s.increments("key_str")
	if v != "value of key_str is string [str]" {
		t.Errorf("[increments 'key_str'] Expected 'value of key_str is string [str]', but got '%v'", v)
	}
	/*** increments Tests ***/

	/*** decrements Tests ***/
	v = s.decrements("key")
	if v != "0" {
		t.Errorf("[decrements 'key'] Expected 0, but got '%v'", v)
	}

	v = s.decrements("key_str")
	if v != "value of key_str is string [str]" {
		t.Errorf("[increments 'key_str'] Expected 'value of key_str is string [str]', but got '%v'", v)
	}
	/*** decrements Tests ***/

	/*** sets Tests ***/
	if s.sets("new_key", "1234"); s.storeInt["new_key"] != 1234 {
		t.Errorf("[sets 'new_key: 1234'] Expected '1234', but got '%v'", s.storeInt["new_key"])
	}

	if s.sets("new_key_str", "test"); s.storeString["new_key_str"] != "test" {
		t.Errorf("[sets 'new_key_str: test'] Expected 'test', but got '%v'", s.storeString["new_key_str"])
	}
	/*** sets Tests ***/

	/*** delete Tests ***/
	if s.delete("new_key"); stringInMap(s.storeInt, "new_key") {
		t.Errorf("[delete 'new_key'] Expected 'false', but got '%v', %v", stringInMap(s.storeInt, "new_key"), s.storeInt)
	}

	s.delete("new_key_str")
	if _, ok := s.storeString["new_key_str"]; ok {
		t.Errorf("[delete 'new_key_str'] Expected 'false', but got '%v', %v", ok, s.storeString)
	}
	/*** delete Tests ***/

}
