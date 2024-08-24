package main

import (
	"fmt"
	"math/rand"
	"time"
)

func f01(i int) {
	for j :=0; j < 3; j++ {
		fmt.Println(i, ":", j)
	}
}

func testConcSample01() {
	for i := 0; i < 8; i++ {
		go f01(i)
	}

	var input string
	fmt.Scanln(&input)
}

func f02(i int) {
	for j := 0; j < 5; j++ {
		fmt.Println(i, ":", j)
		amt := time.Duration(rand.Intn(100))
		time.Sleep(time.Millisecond * amt)
	}
}

func testConcSample02(i int) {
	for i := 0; i < 5; i++ {
		go f02(i)
		amt := time.Duration(rand.Intn(250))
		time.Sleep(time.Millisecond * amt)
	}
	var input string
	fmt.Scanln(&input)
}

func pinger(c chan string) {
	for i := 0; ; i++ {
		c <- "ping"
	}
}

func printer(c chan string) {
	for {
		msg := <- c
		fmt.Println(msg)
		time.Sleep(time.Second * 1)
	}
}

func testConcSample03() {
	var c chan string = make(chan string)
	go pinger(c)
	go printer(c)
	var input string
	fmt.Scanln(&input)
}

func ponger(c chan string) {
	for i := 0; ; i++ {
		c <- "pong"
	}
}

func printer2(c chan string) {
	for {
		msg := <- c
		fmt.Println(msg)
		time.Sleep(time.Second * 1)
	}
}

func testConcSample04() {
	var c chan string = make(chan string)
	go pinger(c)
	go ponger(c)
	go printer2(c)
	var input string
	fmt.Scanln(&input)
}

func testConcSample06() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		for {
			c1 <- "C1F1"
			time.Sleep(time.Millisecond * 250)
		}
	}()

	go func() {
		for {
			c1 <- "C1F2"
			time.Sleep(time.Millisecond * 500)
		}
	}()

	go func() {
		for {
			c1 <- "C1F3"
			time.Sleep(time.Millisecond * 500)
		}
	}()

	go func(){
		for {
			c2 <- "C2F1"
			time.Sleep(time.Millisecond * 250)
		}
	}()

	go func() {
		for {
			c2 <- "C2F2"
			time.Sleep(time.Millisecond * 500)
		}
	}()

	go func() {
		for {
			c2 <- "C2F3"
			time.Sleep(time.Millisecond * 500)
		}
	}()

	go func() {
		for {
			select {
			case msg1 := <- c1:
				fmt.Println(msg1)
			case msg2 := <- c2:
				fmt.Println(msg2)
			}
		}
	}()

	var input string
	fmt.Scanln(&input)
}

func main() {
	// Remember to run this for
	// i=4, j=3
	// i=8, j=3
	// i=5, j=5
	// testConcSample01()

	// Remember to run this for
	// d = 100, d = 500, d = 1000
	// testConcSample02(100)

	//testConcSample03()

	//testConcSample04()

	testConcSample06()
}
