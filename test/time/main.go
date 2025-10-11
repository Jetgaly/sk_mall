package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now().Add(-5 * time.Minute)
	fmt.Println(now)
}
