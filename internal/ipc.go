package internal

import "fmt"

type KeyValue interface {
	Key() string
	fmt.Stringer
}
