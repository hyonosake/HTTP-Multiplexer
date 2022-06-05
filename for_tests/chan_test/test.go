package main

import "fmt"

type Job struct {
	id    int
	stuff []string
}

func main() {

	TestSomeStuff()
}

func TestSomeStuff() {

	jobz := make(chan int, 10)
	for i := 1; i < 36; i++ {
		select {
		case jobz <- i:
		default:
			close(jobz)
			collect(jobz)
			jobz = make(chan int, 10)
		}
	}
	close(jobz)
	collect(jobz)
}

func collect(in <-chan int) {
	for id := range in {
		fmt.Println("in + ", id)
	}

}
