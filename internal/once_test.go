package internal

import (
	"fmt"
	"sync"
	"testing"
)

func Test_GOLANG_ONCE(t *testing.T) {
	var count int

	increment := func() {
		count++
	}

	var once sync.Once

	var increments sync.WaitGroup
	increments.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer increments.Done()
			once.Do(increment)
		}()
	}

	increments.Wait()
	fmt.Printf("Count is %d\n", count)
}

func Test_ONCE_DEADLOCK(t *testing.T) {
	var onceA, onceB sync.Once
	var initB func()
	initA := func() { onceB.Do(initB) }
	initB = func() { onceA.Do(initA) } // 1
	onceA.Do(initA) // 2
	//This program will deadlock because the call to Do at 1 won’t proceed until the call to Do at 2
}

