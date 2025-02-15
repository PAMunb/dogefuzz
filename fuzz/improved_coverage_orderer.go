package fuzz

import (
	"sort"

	"github.com/dogefuzz/dogefuzz/pkg/dto"
)

type improvedCoverageBasedOrderer struct{}

func newImprovedCoverageBasedOrderer() *improvedCoverageBasedOrderer {
	return &improvedCoverageBasedOrderer{}
}

func (o *improvedCoverageBasedOrderer) OrderTransactions(transactions []*dto.TransactionDTO) {
	assingEnergy(transactions)

	sort.SliceStable(transactions, func(i, j int) bool {
		return transactions[i].Energy > transactions[j].Energy
	})
}

func assingEnergy(transactions []*dto.TransactionDTO) {
	for _, transaction := range transactions {
		score := 0.0
		for _, t := range transactions {
			if len(transaction.ExecutedInstructions) != len(t.ExecutedInstructions) {
				continue
			}
			
			flag := true
			sort.Strings(t.ExecutedInstructions)
			sort.Strings(transaction.ExecutedInstructions)

			for i, v := range t.ExecutedInstructions {
				if v != transaction.ExecutedInstructions[i] {
					flag = false
					break
				}
			}

			if flag {
				score += 1
			}
		}

		transaction.Energy = ComputeEnergy(score)
	}
}

// func mapperExecutedInstructions(transactions []*dto.TransactionDTO) map[string]float64 {
// 	executedInstructions := make(map[string]float64)

// 	for _, t := range transactions {
// 		for _, instruction := range t.ExecutedInstructions {
// 			executedInstructions[instruction]++
// 		}
// 	}

// 	return executedInstructions
// }

// func scoreCoverage(hitMap map[string]float64, transaction *dto.TransactionDTO) float64 {
// 	var score float64 = ZERO

// 	for _, execInstruction := range transaction.ExecutedInstructions {
// 		score += hitMap[execInstruction]
// 	}

// 	if score == ZERO {
// 		return math.MaxFloat64
// 	}

// 	return score
// }
