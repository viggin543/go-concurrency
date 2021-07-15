package internal

import (
	"fmt"
	"sync"
	"testing"
)

func TestWaitGroupSimple(t *testing.T) {
	var wg sync.WaitGroup
	salutation := "hello"
	wg.Add(1)
	go func() {
		defer wg.Done()
		salutation = "welcome"
	}()
	wg.Wait()
	fmt.Println(salutation) // will print welcome
}

func TestWaitGroupLoop(t *testing.T) {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(salutation) // good day x 3
			//todo: use idea to create intermediate variable
		}()
	}
	wg.Wait()
}
