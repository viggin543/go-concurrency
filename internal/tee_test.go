package internal

import (
	"fmt"
	"testing"
)


func Test_Tee(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	out1, out2 := Tee(done, Take(done, Repeat(done, 1, 2), 4))

	for val1 := range out1 {
		fmt.Println()
		fmt.Printf("out1: %v, out2: %v\n", val1, <-out2)
	}
}
