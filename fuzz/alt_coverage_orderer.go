package fuzz

import (
	"sort"

	"github.com/dogefuzz/dogefuzz/pkg/dto"
)

type altCoverageBasedOrderer struct{}

func newAltCoverageBasedOrderer() *altCoverageBasedOrderer {
	return &altCoverageBasedOrderer{}
}

func (o *altCoverageBasedOrderer) OrderTransactions(transactions []*dto.TransactionDTO) {
	mapped := mapperExecutedInstructions(transactions)

	sort.SliceStable(transactions, func(i, j int) bool {
		return computeScore(mapped, transactions[i]) < computeScore(mapped, transactions[j])
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

func computeScore(hitMap map[string]float64, transaction *dto.TransactionDTO) float64 {
	var score float64 = 0.0

	for _, execInstruction := range transaction.ExecutedInstructions {
		score += (1 / hitMap[execInstruction])
	}

	return score
}
