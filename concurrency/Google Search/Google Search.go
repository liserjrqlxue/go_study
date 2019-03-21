package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Result string

// START1 OMIT
func Google(query string) (results []Result) {
	c := make(chan Result)
	go func() {
		c <- Web(query)
	}()
	go func() {
		c <- Image(query)
	}()
	go func() {
		c <- Video(query)
	}()
	timeout := time.After(80 * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-timeout:
			fmt.Println("time out")
			return
		}
	}
	return
}

// STOP1 OMIT

// START2 OMIT
var (
	Web   = fakeSearch("web")
	Image = fakeSearch("image")
	Video = fakeSearch("video")
)

type Search func(query string) Result // HL

func fakeSearch(kind string) Search {
	return func(query string) Result {
		t := rand.Intn(100)
		fmt.Printf("use %d millisecond for %s\n", t, kind)
		time.Sleep(time.Duration(t) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q\t%d millisecond\n", kind, query, t))
	}
}

// STOP2 OMIT

func main() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	results := Google("golang")
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}
