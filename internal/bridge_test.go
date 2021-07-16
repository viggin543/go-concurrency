package internal

import (
	"fmt"
	"testing"
)

func Test_Bridge(t *testing.T) {
	genVals := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)
			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i // this blocks until some one will read from this chan
				close(stream)
				chanStream <- stream
			}
		}()
		return chanStream
	}

	for v := range Bridge(nil, genVals()) {
		fmt.Printf("%v ", v)
	}
}
