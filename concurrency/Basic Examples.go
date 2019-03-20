package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Message struct {
	str  string
	wait chan bool
}

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
		//fmt.Printf("You say: %q\n", <-c)
		msg1 := <-c
		fmt.Println(msg1.str)
		msg2 := <-c
		fmt.Println(msg2.str)
		msg1.wait <- true
		msg2.wait <- true
	}
	fmt.Println("You're boring; I'm leaving.")
	ts = append(ts, time.Now())
	fmt.Printf("total time took %7.3fs\n", ts[1].Sub(ts[0]).Seconds())

	for i := 0; i < 5; i++ {

	}

}

func boring(msg string) <-chan Message { // Returns receive-only channel of strings.
	c := make(chan Message)
	waitForIt := make(chan bool) // Shared between all messages.
	go func() {                  // launch goroutine from inside the function
		sum := 0
		for i := 0; ; i++ {
			t := rand.Intn(1e3)
			sum += t
			c <- Message{
				fmt.Sprintf("%s: %d %d", msg, i, t),
				waitForIt,
			}
			time.Sleep(time.Duration(t) * time.Millisecond)
			<-waitForIt
		}
	}()
	return c // Return the channel to the caller
}

func fanIn(input1, input2 <-chan Message) <-chan Message {
	c := make(chan Message)
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
