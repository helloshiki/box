package rand

import (
	"math/rand"
	"time"
)

func Between(min, max int) int {
	rand.Seed(time.Now().Unix())
	return min + rand.Intn(max - min)
}

