package internal

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestSelectStatements(t *testing.T) {
	start := time.Now()
	c := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(c)
	}()

	fmt.Println("Blocking on read...")
	select {
	case <-c:
		fmt.Printf("Unblocked %v later.\n", time.Since(start))
	}
	//If none of the channels are ready, the entire select statement blocks.
	//Then when one the channels is ready,
	//that operation will proceed, and its corresponding statements will execute
}

func TestSelectParallel(t *testing.T) {
	c1 := make(chan interface{})
	close(c1)
	c2 := make(chan interface{})
	close(c2)

	var c1Count, c2Count int
	for i := 1000; i >= 0; i-- {
		select {
		case <-c1:
			c1Count++
		case <-c2:
			c2Count++
		}
	}
	//The Go runtime will perform a pseudo-random uniform selection over the set of case statements
	//the set of case statements, each has an equal chance of being selected as all the others
	fmt.Printf("c1Count: %d\nc2Count: %d\n", c1Count, c2Count)
}

func TestChanTimeout(t *testing.T) {
	var c <-chan int
	select {
	case <-c:
	case <-time.After(1 * time.Second):
		fmt.Println("Timed out.")
	}
}

func TestSelectDefaultBranch(t *testing.T) {
	start := time.Now()
	var c1, c2 <-chan int
	select {
	case <-c1:
	case <-c2:
	default:
		fmt.Printf("In default after %v\n\n", time.Since(start))
	}
}
func TestSelectDefault(t *testing.T) {
	done := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	workCounter := 0
loop:
	for {
		select {
		case <-done:
			break loop
		default:
		}

		// Simulate work
		workCounter++
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("Achieved %v cycles of work before signalled to stop.\n", workCounter)
}

func TestSelectBlockForEver(t *testing.T) {
	select {} // block forever
}

func TestControlWorkerThreadPool(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

