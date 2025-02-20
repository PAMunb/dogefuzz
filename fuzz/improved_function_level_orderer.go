package fuzz

import (
	"math"
	"sort"

	"github.com/dogefuzz/dogefuzz/pkg/dto"
)

type improvedFunctionLevelBasedOrderer struct {
	contract *dto.ContractDTO
}

func newImprovedFunctionLevelBasedOrderer(contract *dto.ContractDTO) *improvedFunctionLevelBasedOrderer {
	return &improvedFunctionLevelBasedOrderer{contract}
}

func (o *improvedFunctionLevelBasedOrderer) OrderTransactions(transactions []*dto.TransactionDTO) {
	o.assingEnergy(transactions)
	sort.SliceStable(transactions, func(i, j int) bool {
		return transactions[i].Energy > transactions[i].Energy
	})
}

func (o *improvedFunctionLevelBasedOrderer) assingEnergy(transactions []*dto.TransactionDTO) {
	var minDist = math.MaxFloat64
	var maxDist float64 = 0
	var BIG_NUMBER float64 = 999999999

	for _, transaction := range transactions {
		if transaction.Distance < BIG_NUMBER {
			if transaction.Distance < minDist {
				minDist = transaction.Distance
			}
			if transaction.Distance > maxDist {
				maxDist = transaction.Distance
			}
		}
	}

	for _, transaction := range transactions {
		if transaction.Distance == minDist {
			if minDist == maxDist {
				transaction.Energy = 1
			} else {
				transaction.Energy = maxDist - minDist
			}
		} else {
			transaction.Energy = (maxDist - minDist) / (transaction.Distance - minDist)
		}
	}
}
