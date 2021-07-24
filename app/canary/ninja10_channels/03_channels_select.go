package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	beeps := beep()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGKILL)

	fmt.Println("Listening for the beeps")
	for {
		select {
		case message := <-beeps:
			fmt.Printf("Message received: %s\n", message)
		case <-quit:
			fmt.Println("That Beep is anoying, quit listening!")
			return
		}
	}
}

func beep() <-chan string {
	beeps := make(chan string, 1024)
	go func() {
		for i := 0; ; i++ {
			beeps <- fmt.Sprintf("beep %d", i)
			time.Sleep(10000 * time.Millisecond)
		}
	}()
	return beeps
}
