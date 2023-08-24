package main

import "fmt"

func main() {
	ch := make(chan int)
	select {
	case ch <- 1:
		fmt.Println("send success")
	default:
		fmt.Println("send failed")
	}
	fmt.Println(333)
}
