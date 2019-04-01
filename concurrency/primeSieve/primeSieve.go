// A concurrent prime sieve

package main

import "time"

// Send the sequence 2, 3, 4, ... to channel 'ch'.
func Generate(ch chan<- int) {
	for i := 2; ; i++ {
		ch <- i // Send 'i' to channel 'ch'
	}
}

// Copy the values from channel 'in' to channel 'out'
// removing those divisible by 'prime'
func Filter(in <-chan int, out chan<- int, prime int) {
	for {
		i := <-in // Reseive value from 'in'.
		if i%prime != 0 {
			out <- i // Send 'i' to 'out'
		}
	}
}

// The prime sieve: Daisy-chain Filter processes.
func main() {
	t0 := time.Now()
	ch := make(chan int) // Create a new channel.
	go Generate(ch)      // Launch Generate goroutine.
	for i := 0; i < 10000; i++ {
		prime := <-ch
		//print(prime,"\n")
		ch1 := make(chan int)
		go Filter(ch, ch1, prime)
		ch = ch1
	}
	count := <-ch
	d := time.Now().Sub(t0)
	print(count, "\n")
	print(d.Nanoseconds()/int64(count), "\n")
	print(d.Nanoseconds(), "\n")
	print(d.String(), "\n")
}
