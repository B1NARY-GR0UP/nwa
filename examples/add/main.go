package main

import "fmt"

func main() {
	//	ch1 := make(chan struct{})
	//	go func() {
	//		time.Sleep(time.Second * 3)
	//		ch1 <- struct{}{}
	//	}()
	//LOOP:
	//	for {
	//		select {
	//		case <-ch1:
	//			fmt.Println("<-ch1")
	//			break LOOP
	//		default:
	//			time.Sleep(time.Second * 1)
	//			fmt.Println("default")
	//		}
	//		fmt.Println("for")
	//	}
	ch := make(chan int)
	select {
	case ch <- 1:
		fmt.Println("send success")
	default:
		fmt.Println("send failed")
	}
	fmt.Println(333)
}
