package internal

import (
	"fmt"

	"sync"
	"testing"

)

func TestChannel(t *testing.T) {
	var dataStream chan interface{}
	dataStream = make(chan interface{})
	dataStream <- ""
}

func TestReadOnlyChan(t *testing.T) {
	var dataStream <-chan interface{}
	dataStream = make(<-chan interface{})
	<-dataStream
}

func TestWriteOnlyChan(t *testing.T) {
	var dataStream chan<- interface{}
	dataStream = make(chan<- interface{})
	dataStream <- ""
}

func TestAutoCastingChannels(t *testing.T) {
	var receiveChan <-chan interface{}
	var sendChan chan<- interface{}
	dataStream := make(chan interface{})

	// Valid statements:
	receiveChan = dataStream
	sendChan = dataStream

	sendChan <- receiveChan
}

func TestChanDeadLock(t *testing.T) {
	stringStream := make(chan string)
	go func() {
		if true {
			return
		}
		stringStream <- "Hello channels!"
	}()
	str, ok := <-stringStream
	fmt.Println(str, ok)
	//ok is false if chan was closed
}

func TestReadingFromClosedChannel(t *testing.T) {
	intStream := make(chan int)
	close(intStream)
	integer, ok := <-intStream
	fmt.Printf("(%v): %v", ok, integer)
	// ok is false
	// integer is 0
	// you can read many times from a closed chan
}

func TestRangingOverAStream(t *testing.T) {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for i := 1; i <= 5; i++ {
			intStream <- i
		}
	}()

	for integer := range intStream {
		fmt.Printf("%v ", integer)
	}
}

func TestBroadcastClose(t *testing.T) {
	begin := make(chan interface{})
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin // waits for close
			fmt.Printf("%v has begun\n", i)
		}(i)
	}

	fmt.Println("Unblocking goroutines...")
	close(begin) // will unleash all waiters ( shorter than cond.Broadcast() )
	wg.Wait()
}

func TestBufferedChannel(t *testing.T) {
	var dataStream chan interface{}
	dataStream = make(chan interface{}, 4)
	//even if no reads are performed on the channel,
	//a goroutine can still perform 4 writes,
	dataStream <- ""
}
