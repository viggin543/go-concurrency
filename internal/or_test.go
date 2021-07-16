package internal

import (
	"fmt"
	"testing"
)

func Test_OrChannel(t *testing.T) {
	done := make(chan interface{})
	ints := Take(done, Repeat(done, 1, 2, 3), 10)
	defer close(done)
	for x := range Or(done, ints) {
		fmt.Println(x)
	}
}
