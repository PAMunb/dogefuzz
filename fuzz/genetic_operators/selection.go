package genetic_operators

import (
	"math"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

// consts to define chances to select some seeds
const FIRST_LIMIT = 0.5
const SECOND_LIMIT = 0.8

// consts to define what range interval will used to prioritize some seeds
const FIRST_RANGE = 0.4
const SECOND_RANGE = 0.7

func Selection(seedsList [][]string) []string {
	rnd := common.RandomFloat64()
	lengthList := len(seedsList)
	var seedsSlice [][]string

	if rnd < FIRST_LIMIT {
		seedsSlice = seedsList[0:maxIndex(lengthList, FIRST_RANGE)]
	} else if rnd < SECOND_LIMIT {
		seedsSlice = seedsList[minIndex(lengthList, FIRST_RANGE):maxIndex(lengthList, SECOND_RANGE)]
	} else {
		// [minIndex:len(seedsList)]
		seedsSlice = seedsList[minIndex(lengthList, SECOND_RANGE):]
	}

	return common.RandomChoice(seedsSlice)
}

func maxIndex(length int, factor float64) int {
	return int(math.Ceil(float64(length) * factor))
}

func minIndex(length int, factor float64) int {
	return int(math.Floor(float64(length) * factor))
}
