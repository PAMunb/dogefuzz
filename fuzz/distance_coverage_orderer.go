package fuzz

import (
	"math"
	"sort"

	"github.com/dogefuzz/dogefuzz/pkg/dto"
)

const WEIGHT_1 = 2
const WEIGHT_2 = 4
const WEIGHT_3 = 4

type distanceCoverageBasedOrderer struct {
	contract *dto.ContractDTO
}

func newDistanceCoverageBasedOrderer(contract *dto.ContractDTO) *distanceCoverageBasedOrderer {
	return &distanceCoverageBasedOrderer{contract}
}

func (o *distanceCoverageBasedOrderer) OrderTransactions(transactions []*dto.TransactionDTO) {
	sort.SliceStable(transactions, func(i, j int) bool {
		return o.computeScore(transactions[i]) > o.computeScore(transactions[j])
	})
}

func (o *distanceCoverageBasedOrderer) computeScore(transaction *dto.TransactionDTO) float64 {
	// 2, 4, 4
	return WEIGHT_1*o.computeCriticalInstructionsHits(transaction) +
		WEIGHT_2*o.computeDistance(transaction) +
		WEIGHT_3*o.computeCoverage(transaction)
}

func (o *distanceCoverageBasedOrderer) computeCoverage(transaction *dto.TransactionDTO) float64 {
	var totalInstructions = len(o.contract.CFG.Instructions)
	var executedInstructions = len(transaction.ExecutedInstructions)

	if totalInstructions != 0 {
		return float64(executedInstructions) / float64(totalInstructions)
	}

	return 0
}

func (o *distanceCoverageBasedOrderer) computeDistance(transaction *dto.TransactionDTO) float64 {
	var maxDistance map[string]uint32
	var distanceSum int64 = 0
	var distancePercentage float64 = 0
	var minDistance uint64 = transaction.DeltaMinDistance

	for _, distance := range o.contract.DistanceMap {
		if maxDistance == nil {
			maxDistance = make(map[string]uint32)
			for pc := range distance {
				maxDistance[pc] = 0
			}
		}

		for instr := range maxDistance {
			if val, ok := distance[instr]; ok {
				if val != math.MaxUint32 && val > maxDistance[instr] {
					maxDistance[instr] = val
				}
			}
		}
	}

	for _, distance := range maxDistance {
		distanceSum += int64(distance)
	}

	if minDistance >= uint64(math.MaxUint32) {
		minDistance -= math.MaxUint32
	}

	if distanceSum != 0 {
		distancePercentage = float64(minDistance) / float64(distanceSum)
	}

	return distancePercentage
}

func (o *distanceCoverageBasedOrderer) computeCriticalInstructionsHits(transaction *dto.TransactionDTO) float64 {
	if o.contract.TargetInstructionsFreq != 0 {
		return float64(transaction.CriticalInstructionsHits) / float64(o.contract.TargetInstructionsFreq)
	}

	return 0
}
