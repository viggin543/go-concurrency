package internal

import (
	"fmt"
	"testing"
)
//As a refresher, remember that Go follows a fork-join model for concurrency. Forks are when goroutines are started, and join points are when two or more goroutines are synchronized through channels or types in the sync package. The work stealing algorithm follows a few basic rules. Given a thread of execution:
//
//    1 - At a fork point, add tasks to the tail of the deque associated with the thread.
//
//    2 - If the thread is idle, steal work from the head of deque associated with some other random thread.
//
//    3 - At a join point that cannot be realized yet (i.e., the goroutine it is synchronized with has not completed yet), pop work off the tail of the thread’s own deque.
//
//    5 - If the thread’s deque is empty, either:
//
//        Stall at a join.
//
//        Steal work from the head of a random thread’s associated deque.

func fib(n int) <-chan int {
	result := make(chan int)
	go func() { // task
		defer close(result)
		if n <= 2 {
			result <- 1
			return
		}
		result <- <-fib(n-1) + <-fib(n-2) // join point
	}()
	return result // continuation
}

// goroutines scheduling is continuation scheduling

func TestWorkStealing(t *testing.T) {
	fmt.Printf("fib(4) = %d", <-fib(16))
}
//What Go does in this situation is dissociate the context from the OS thread so that the context can be handed off to another,
//unblocked, OS thread.
//This allows the context to schedule further goroutines, which allows the runtime to keep the host machine’s CPUs active.
//The blocked goroutine remains associated with the blocked thread.
//
//When the goroutine eventually becomes unblocked, the host OS thread attempts to steal back a context from one of the other OS threads
//so that it can continue executing the previously blocked goroutine. However, sometimes this is not always possible.
//In this case, the thread will place its goroutine on a global context, the thread will go to sleep,
//and it will be put into the runtime’s thread pool for future use (for instance, if a goroutine becomes blocked again).
//
//The global context we just mentioned doesn’t fit into our prior discussions of abstract work-stealing algorithms.
//It’s an implementation detail that is necessitated by how Go is optimizing CPU utilization.
//To ensure that goroutines placed into the global context aren’t there perpetually,
//a few extra steps are added into the work-stealing algorithm. Periodically,
//a context will check the global context to see if there are any goroutines there,
//and when a context’s queue is empty, it will first check the global context for work to steal before checking other OS threads’ contexts.
//
//Other than input/output and system calls,
//Go also allows goroutines to be preempted during any function call.
//This works in tandem with Go’s philosophy of preferring very fine-grained concurrent tasks by ensuring the runtime can efficiently schedule work.
//One notable exception that the team has been trying to solve is goroutines that perform no input/output, system calls,
//or function calls. Currently,
//these kinds of goroutines are not preemptable and can cause significant issues like long GC waits,
//or even deadlocks. Fortunately, from an anecdotal perspective, this is a vanishingly small occurrence.

// 			A VERY NICE POST ON THIS
//https://morsmachine.dk/go-scheduler

// GOLANG IO
//https://morsmachine.dk/netpoller

// GR8 posts on event loops
//https://medium.com/ing-blog/how-does-non-blocking-io-work-under-the-hood-6299d2953c74
//https://dev.to/frevib/a-tcp-server-with-kqueue-527