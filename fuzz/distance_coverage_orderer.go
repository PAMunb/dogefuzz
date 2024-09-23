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
	mapped := mapperExecutedInstructions(transactions)

	sort.SliceStable(transactions, func(i, j int) bool {
		return evaluateSeedCoverage(mapped, transactions[i]) < evaluateSeedCoverage(mapped, transactions[j])
	})
}

func mapperExecutedInstructions(transactions []*dto.TransactionDTO) map[string]float64 {
	executedInstructions := make(map[string]float64)

	for _, t := range transactions {
		for _, instruction := range t.ExecutedInstructions {
			executedInstructions[instruction] += 1
		}
	}

	return executedInstructions
}

// compute coverage
func evaluateSeedCoverage(hitMap map[string]float64, transaction *dto.TransactionDTO) float64 {
	var score float64 = 0.0
	var a float64 = 5

	for _, execInstruction := range transaction.ExecutedInstructions {
		score += (1 / hitMap[execInstruction])
	}

	return math.Pow(score, a)
}
