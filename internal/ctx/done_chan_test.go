package ctx

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

//var Canceled = errors.New("context canceled")
//var DeadlineExceeded error = deadlineExceededError{}
//
//type CancelFunc
//type Context
//
//func Background() Context
//func TODO() Context
//func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
//func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)
//func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
//func WithValue(parent Context, key, val interface{}) Context

func Test_before_context_demonstration(t *testing.T) {
	var wg sync.WaitGroup
	done := make(chan interface{})
	defer close(done)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printGreeting(done); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewell(done); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}()

	wg.Wait()

}

func printGreeting(done <-chan interface{}) error {
	greeting, err := genGreeting(done)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", greeting)
	return nil
}

func printFarewell(done <-chan interface{}) error {
	farewell, err := genFarewell(done)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", farewell)
	return nil
}

func genGreeting(done <-chan interface{}) (string, error) {
	switch locale, err := locale(done); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func genFarewell(done <-chan interface{}) (string, error) {
	switch locale, err := locale(done); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func locale(done <-chan interface{}) (string, error) {
	select {
	case <-done:
		return "", fmt.Errorf("canceled")
	case <- time.After(1 * time.Minute):
	}
	return "EN/US", nil
}
