// Package timer with pause /**
package timer

import (
	"fmt"
	"time"
)

// todo
func main() {
	timer := time.NewTimer(2 * time.Second)

	<-timer.C
	fmt.Println("Timer fired!")
}
