package genetic_operators

import "github.com/dogefuzz/dogefuzz/pkg/common"

const MUTATION_CHANCE = 0.1

func Mutation(mutationFunction []func()) {
	rnd := common.RandomFloat64()

	if rnd < MUTATION_CHANCE {
		mutationFunction := common.RandomChoice(mutationFunction)
		mutationFunction()
	}
}
