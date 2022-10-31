# fun-exercise
This is just a fun golang exercise.

# Instructions
Write a functional key/value store that implements the following functions:

* INCR <key>: With a value starting at zero, INCR increments the value associated to the key by 1.
* DECR <key>: With a value starting at zero, DECR decrements the value associated to the key by 1.
* GET <key>: Returns the value associated with that key.  If the key does not exist, it returns an error to the user.
* SET <key> <value>: Sets a key to be equal the value.

In addition to the functions, the key/value store must support transactions:

* BEGIN: Starts a transaction
* COMMIT: Commits (saves) a transaction.  If a transaction has not been started it returns an error.
* ROLLBACK: Rolls the key/value state back to the beginning of a transaction disregarding all of the changes that have been made. If a transaction has not been started it returns an error.

All input is taken from the command line.

Examples:
```
> GET mykey
mykey does not exist
> INCR mykey
1
> INCR mykey
2
> GET mykey
2
> SET anotherkey 5
5
> GET anotherkey
5
> BEGIN
ok
> SET anotherkey 10
10
> GET anotherkey
10
> DELETE mykey
ok
> GET mykey
mykey does not exist
> ROLLBACK
ok
> GET anotherkey
5
> GET mykey
2
> BEGIN
ok
> SET anotherkey 15
15
> COMMIT
ok
> GET anotherkey
15
```


## Nested transactions

```
INCR key
1
BEGIN
ok
INCR key
2
SET hello world
world
BEGIN
ok
INCR key
3
SET anotherstring hello
hello
ROLLBACK
ok
GET key
2
GET hello
world
GET anotherstring
anotherstring does not exist
COMMIT
```

```go
// type Store map[string]interface{}
type Store map[string]Element

type Element struct {
  Value interface{} // the actual (or committed) value.
  Transactions map[string]interface{} // map of transaction ids pointing to values associated with the transaction.
}

type Transaction struct {
  Id string
  Keys map[string]bool
}

type TransactionStack []Transaction
```

```go
type Store struct {
  data map[string]interface{}
  next *Store
  prev *Store
}

func NewStore(parent *Store) *Store {
  store := &Store{
    prev: parent,
  }

  if parent != nil {
    parent.next(store)
  }

  return store
}

func (s *Store) Get(key string) interface{} {
  if _, ok := s.data[key]; !ok {
    return s.prev.Get(key)
  }
  return s.data[key]
}

func (s *Store) Set(key string, val interface{}) {
  s.data[key] = val
}

func (s *Store) EndTransaction() *Store {
  for k, v := range s.data {
    s.prev[k] = v
  }
  s.prev.next = nil
  return s.prev
}


store := NewStore(nil)

storev1 := NewStore(store)
```

INCR key
1
SET other 10
BEGIN: creates a transaction id, pushes the transaction id to the stack.
INCR key: Do we have a transaction? yes, get the id, get the 'key' element, see if we have a transaction already in the transaction map? no - incr value and place in transaction map. elem['key'].Transactions['id'] = 2
INCR key: -
COMMIT