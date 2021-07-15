package internal

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

func Test_pipe(t *testing.T) {
	generator := func(done <-chan interface{}, integers ...int) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for _, i := range integers {
				select {
				case <-done:
					return
				case intStream <- i:
				}
			}
		}()
		return intStream
	}

	multiply := func(
		done <-chan interface{},
		intStream <-chan int,
		multiplier int,
	) <-chan int {
		multipliedStream := make(chan int)
		go func() {
			defer close(multipliedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case multipliedStream <- i * multiplier:
				}
			}
		}()
		return multipliedStream
	}

	add := func(
		done <-chan interface{},
		intStream <-chan int,
		additive int,
	) <-chan int {
		addedStream := make(chan int)
		go func() {
			defer close(addedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case addedStream <- i + additive:
				}
			}
		}()
		return addedStream
	}

	done := make(chan interface{})
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4)
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	for v := range pipeline {
		fmt.Println(v)
	}
}

var Take = func(
	done <-chan interface{},
	valueStream <-chan interface{},
	num int,
) <-chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()
	return takeStream
}
var Repeat = func(
	done <-chan interface{},
	values ...interface{},
) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case valueStream <- v:
				}
			}
		}
	}()
	return valueStream
}

var repeatFn = func(
	done <-chan interface{},
	fn func() interface{},
) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for {
			select {
			case <-done:
				return
			case valueStream <- fn():
			}
		}
	}()
	return valueStream
}

var toString = func(
	done <-chan interface{},
	valueStream <-chan interface{},
) <-chan string {
	stringStream := make(chan string)
	go func() {
		defer close(stringStream)
		for v := range valueStream {
			select {
			case <-done:
				return
			case stringStream <- v.(string):
			}
		}
	}()
	return stringStream
}

func Test_geenrators(t *testing.T) {

	done := make(chan interface{})
	for num := range Take(done, Repeat(done, 1), 10) {
		fmt.Printf("%v ", num)
	}

	random := func() interface{} { return rand.Int() }

	for num := range Take(done, repeatFn(done, random), 10) {
		fmt.Println(num)
	}
}

func BenchmarkGeneric(b *testing.B) {
	done := make(chan interface{})
	defer close(done)

	b.ResetTimer()
	for range toString(done, Take(done, Repeat(done, "a"), b.N)) {
	}
}

func BenchmarkTyped(b *testing.B) {
	repeat := func(done <-chan interface{}, values ...string) <-chan string {
		valueStream := make(chan string)
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}()
		return valueStream
	}

	take := func(
		done <-chan interface{},
		valueStream <-chan string,
		num int,
	) <-chan string {
		takeStream := make(chan string)
		go func() {
			defer close(takeStream)
			for i := num; i > 0 || i == -1; {
				if i != -1 {
					i--
				}
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()
		return takeStream
	}

	done := make(chan interface{})
	defer close(done)

	b.ResetTimer()
	for range take(done, repeat(done, "a"), b.N) {
	}
}

var toInt = func(
	done <-chan interface{},
	valueStream <-chan interface{},
) <-chan int {
	stringStream := make(chan int)
	go func() {
		defer close(stringStream)
		for v := range valueStream {
			select {
			case <-done:
				return
			case stringStream <- v.(int):
			}
		}
	}()
	return stringStream
}

func Test_PrimeFinder(t *testing.T) {
	rand := func() interface{} { return rand.Intn(7) }
	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	randIntStream := toInt(done, repeatFn(done, rand))
	fmt.Println("Primes:")
	for prime := range Take(done, primeFinder(done, randIntStream), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v", time.Since(start))
}

func primeFinder(done chan interface{}, stream <-chan int) <-chan interface{} {
	primes := make(chan interface{})
	go func() {
		for {
			select {
			case candidate := <-stream:
				isPrime := true
				for prev := candidate - 1; prev > 2; prev-- {
					if candidate%prev == 0 {
						isPrime = false
						break
					}
				}
				if isPrime && candidate > 1 {
					primes <- candidate
				}

			case <-done:
				return
			}
		}
	}()
	return primes
}

var fanIn = func(
	done <-chan interface{},
	channels ...<-chan interface{},
) <-chan interface{} {
	var wg sync.WaitGroup
	multiplexedStream := make(chan interface{})
	multiplex := func(c <-chan interface{}) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- i:
			}
		}
	}

	// Select from all the channels
	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	// Wait for all the reads to complete
	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

func TestFunOut_FanIn(t *testing.T) {
	done := make(chan interface{})
	rand := func() interface{} { return rand.Intn(50000000) }

	defer close(done)

	randIntStream := toInt(done, repeatFn(done, rand))

	numFinders := runtime.NumCPU()
	finders := make([]<-chan interface{}, numFinders)
	for i := 0; i < numFinders; i++ { // fan out
		finders[i] = primeFinder(done, randIntStream)
	}

	for prime := range Take(done, fanIn(done, finders...), 100) {
		fmt.Printf("\t%d\n", prime)
	}

}
