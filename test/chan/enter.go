package main

import "fmt"

func main() {
	c := make(chan int, 2)
	c <- 1
	close(c)
	for i := 0; i < 3; i++ {
		fmt.Println(i)
		v, ok := <-c
		fmt.Println(v,ok)
	}

}
