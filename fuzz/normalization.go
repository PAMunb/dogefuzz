package fuzz

import (
	"math"
	"strconv"
)

/*
* Normalization: A Preprocessing Stage
* https://arxiv.org/pdf/1503.06462
 */
func IntegerScalingNormalization(score uint64) float64 {
	s := strconv.FormatUint(score, 10)
	n := len(s)
	firstN, _ := strconv.ParseFloat(s[:1], 64)
	pTen := math.Pow(10, float64(n-1))

	return (float64(score) - pTen*firstN) / pTen
}
