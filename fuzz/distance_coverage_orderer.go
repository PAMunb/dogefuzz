package fuzz

import (
	"math"
	"sort"

	"github.com/dogefuzz/dogefuzz/pkg/dto"
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
	return math.Max(o.computeCriticalInstructionsHits(transaction), math.Max(o.computeCoverage(transaction), o.computeDistance(transaction)))
}

func (o *distanceCoverageBasedOrderer) computeCoverage(transaction *dto.TransactionDTO) float64 {
	var totalInstructions = len(o.contract.DistanceMap)
	var coveragePercentage float64 = 0

	if totalInstructions != 0 {
		coveragePercentage = float64(transaction.Coverage) / float64(totalInstructions)
	}

	return coveragePercentage
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
	var hitsPercentage float64 = 0

	if o.contract.TargetInstructionsFreq != 0 {
		hitsPercentage = float64(transaction.CriticalInstructionsHits) / float64(o.contract.TargetInstructionsFreq)
	}

	return hitsPercentage
}
