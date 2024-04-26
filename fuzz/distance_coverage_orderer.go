package fuzz

import (
	"math"
	"sort"

	"github.com/dogefuzz/dogefuzz/pkg/dto"
)

const (
	HITS_WEIGHT_MO     = 0.1
	COVERAGE_WEIGHT_MO = 0.6
	DISTANCE_WEIGHT_MO = 0.3
)

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
	return HITS_WEIGHT_MO*o.computeCriticalInstructionsHits(transaction) +
		COVERAGE_WEIGHT_MO*o.computeCoverage(transaction) +
		DISTANCE_WEIGHT_MO*o.computeDistance(transaction)
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

	if distanceSum != 0 {
		return float64(transaction.DeltaMinDistance) / float64(distanceSum)
	}

	return 0
}

func (o *distanceCoverageBasedOrderer) computeCriticalInstructionsHits(transaction *dto.TransactionDTO) float64 {
	if o.contract.TargetInstructionsFreq != 0 {
		return float64(transaction.CriticalInstructionsHits) / float64(o.contract.TargetInstructionsFreq)
	}

	return 0
}
