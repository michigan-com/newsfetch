package lib

import (
	"math/rand"
	"time"
)

func RandomInt(max int) int {
	startTime := time.Now()
	rand.Seed(int64(startTime.Nanosecond()))

	return rand.Intn(max)
}
