package fuzz

import (
	"math"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

const ZERO float64 = 0.0
const POWER float64 = 4.0

var targetInstructions = []string{"CALL", "SELFDESTRUCT", "CALLCODE", "DELEGATECALL"}

func ComputeEnergy(v float64) float64 {
	return 1 / math.Pow(v, POWER)
}

func FindBlocksContainingTargetInstructions(cfg common.CFG, targetInstructions []string) []string {
	targetBlocks := make([]string, 0)
	for blockPC, block := range cfg.Blocks {
		for _, instr := range block.Instructions {
			if common.Contains(targetInstructions, instr) {
				targetBlocks = append(targetBlocks, blockPC)
				break
			}
		}
	}

	return targetBlocks
}

func MeanMinDistances(mapDistanceToTargetInstruction map[string]float64) float64 {
	var score = ZERO

	for _, distance := range mapDistanceToTargetInstruction {
		score += distance
	}

	if score == ZERO {
		return math.MaxFloat64
	}

	return score / float64(len(mapDistanceToTargetInstruction))
}
