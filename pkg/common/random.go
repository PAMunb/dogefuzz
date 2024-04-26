package common

import (
	"math/rand"
	"time"
)

func RandomChoice[T any](slice []T) T {
	rand.Seed(time.Now().UnixNano())
	rndIdx := rand.Intn(len(slice))
	return slice[rndIdx]
}

func RandomFloat64() float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()
}

