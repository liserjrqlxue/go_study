package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	var ts []time.Time
	ts = append(ts, time.Now())
	//c:=make(chan string)
	//go boring("boring!",c)
	joe := boring("Joe")
	ann := boring("Ann")
	c := fanIn(joe, ann)
	fmt.Println("I'm listening.")
	//time.Sleep(2*time.Second)
	for i := 0; i < 50; i++ {
		//fmt.Printf("You say: %q\n",<-joe)
		//fmt.Printf("You say: %q\n",<-ann)
		fmt.Printf("You say: %q\n", <-c)
	}
	fmt.Println("You're boring; I'm leaving.")
	ts = append(ts, time.Now())
	fmt.Printf("total time took %7.3fs\n", ts[1].Sub(ts[0]).Seconds())
}

func boring(msg string) <-chan string { // Returns receive-only channel of strings.
	c := make(chan string)
	go func() { // launch goroutine from inside the function
		sum := 0
		for i := 0; ; i++ {
			t := rand.Intn(1e3)
			sum += t
			c <- fmt.Sprintf("%s %d %d", msg, i, t)
			//time.Sleep(time.Second)
			time.Sleep(time.Duration(t) * time.Millisecond)
		}
		fmt.Println("sum", sum)
	}()
	return c // Return the channel to the caller
}

func fanIn(input1, input2 <-chan string) <-chan string {
	c := make(chan string)
	go func() {
		for {
			c <- <-input1
		}
	}()
	go func() {
		for {
			c <- <-input2
		}
	}()
	return c
}
