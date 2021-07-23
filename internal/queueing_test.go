package internal

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func Sleep(
	done <-chan interface{},
	duration time.Duration,
	in <-chan interface{}) <-chan interface{} {
	ret := make(chan interface{})
	go func() {
		defer close(ret)
		for {
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				time.Sleep(duration)
				ret <- val
			}
		}
	}()
	return ret
}

func Buffer(done <-chan interface{}, size int, in <-chan interface{}) <- chan interface{} {
	ret := make(chan interface{	}, size)
	go func() {
		defer close(ret)
		input := OrDone(done,in)
		for v := range input {
			ret <- v
		}
	}()
	return ret
}

func Test_Queueing(t *testing.T) {
	done := make(chan interface{})
	defer close(done)
	zeros := Take(done, Repeat(done, 0), 3)
	short := Sleep(done, 1*time.Second, zeros)
	buffered := Buffer(done, 2, short)
	long := Sleep(done, 4*time.Second, buffered)
	pipeline := long
	for x := range pipeline {
		fmt.Println(x)
	}
}


func BenchmarkUnbufferedWrite(b *testing.B) {
	performWrite(b, tmpFileOrFatal())
}

func BenchmarkBufferedWrite(b *testing.B) {
	//chunking is faster
	performWrite(b, bufio.NewWriter(tmpFileOrFatal()))
}

func tmpFileOrFatal() *os.File {
	file, err := ioutil.TempFile("", "tmp")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return file
}

func performWrite(b *testing.B, writer io.Writer) {
	done := make(chan interface{})
	defer close(done)

	b.ResetTimer()
	for bt := range Take(done, Repeat(done, byte(0)), b.N) {
		writer.Write([]byte{bt.(byte)})
	}
}