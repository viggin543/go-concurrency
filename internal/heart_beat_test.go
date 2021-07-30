package internal

import (
	"fmt"
	"testing"
	"time"
)

func doHeartBeat(
	done <-chan interface{},
	pulseInterval time.Duration,
) (<-chan interface{}, <-chan time.Time) {
	heartbeat := make(chan interface{})
	results := make(chan time.Time)
	go func() {
		defer close(heartbeat)
		defer close(results)

		pulse := time.Tick(pulseInterval)
		workGen := time.Tick(2 * pulseInterval)

		sendPulse := func() {
			select {
			case heartbeat <- struct{}{}: //select causing this to not block, if no one reads from heartbeat
			default:
			}
		}
		sendResult := func(r time.Time) {
			for { // why the loop ?
				select {
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case results <- r:
					return
				}
			}
		}

		for {
			select {
			case <-done:
				return
			case <-pulse:
				sendPulse()
			case r := <-workGen:
				sendResult(r)
			}
		}
	}()
	return heartbeat, results
}

func TestConsumeHeartBeat(t *testing.T) {
	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2*time.Second
	heartbeat, results := doHeartBeat(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if ok == false {
				return
			}
			fmt.Printf("results %v\n", r.Second())
		case <-time.After(timeout):
			return
		}
	}
}
func TestName(t *testing.T) {
	heartbeat := make(chan interface{})

	sendPulse := func() {
		select {
		case heartbeat <- struct{}{}: // can this be not a select ( NO !)
		default: // The default case in a select is run if no other case is ready.
			fmt.Println("will this happen?")
		}
	}
	sendPulse()
	fmt.Println("opa")
}
