// Timers and Ticker
package ch07

import (
	"fmt"
	"time"
)

func timelyFixed() {
	timer := time.NewTimer(5 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop() // Be sure to stop the ticker!

	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				fmt.Println("tick!")
			}
		}
	}()

	<-timer.C
	fmt.Println("timer done!")
	close(done)
}
