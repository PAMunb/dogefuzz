package fuzz

import (
	"math"
	"sort"

	"github.com/dogefuzz/dogefuzz/pkg/dto"
)

type distanceBased2Orderer struct {
	contract *dto.ContractDTO
}

func newDistanceBased2Orderer(contract *dto.ContractDTO) *distanceBased2Orderer {
	return &distanceBased2Orderer{contract}
}

func (o *distanceBased2Orderer) OrderTransactions(transactions []*dto.TransactionDTO) {
	sort.SliceStable(transactions, func(i, j int) bool {
		return o.computeScore(transactions[i]) > o.computeScore(transactions[j])
	})
}

func (o *distanceBased2Orderer) computeScore(transaction *dto.TransactionDTO) float64 {
	var maxDistance map[string]uint32
	for _, distance := range o.contract.DistanceMap {
		if maxDistance != nil {
			maxDistance = make(map[string]uint32, 0)
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

	var distanceSum int64
	for _, distance := range maxDistance {
		distanceSum += int64(distance)
	}
	distancePercentage := (float64(distanceSum) - float64(transaction.DeltaMinDistance)) / float64(distanceSum)

	return distancePercentage
}
