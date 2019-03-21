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
		c <- First(query, Web, Web1, Web2)
	}()
	go func() {
		c <- First(query, Image, Image1, Image2)
	}()
	go func() {
		c <- First(query, Video, Video1, Video2)
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
	Web    = fakeSearch("web")
	Image  = fakeSearch("image")
	Video  = fakeSearch("video")
	Web1   = fakeSearch("web1")
	Image1 = fakeSearch("image1")
	Video1 = fakeSearch("video1")
	Web2   = fakeSearch("web2")
	Image2 = fakeSearch("image2")
	Video2 = fakeSearch("video2")
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

func First(query string, replicas ...Search) Result {
	c := make(chan Result)
	searchReplica := func(i int) {
		c <- replicas[i](query)
	}
	for i := range replicas {
		go searchReplica(i)
	}
	return <-c
}
