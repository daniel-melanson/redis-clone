package redis

// import (
//   "bufio"
//   "fmt"
//   "io"
//   "strconv"
// )

const (
  STRING = '+'
  ERROR = '-'
  INTEGER = ':'
  BULK = '$'
  ARRAY = '*'
)

type Value struct {
  kind string
  str string
  num int
  bulk string
  array []Value
}
