// package fuzz

// import (
// 	"math"
// 	"sort"

// 	"github.com/dogefuzz/dogefuzz/pkg/dto"
// )

// type improvedFunctionLevelBasedOrderer struct {
// 	contract *dto.ContractDTO
// }

// func newImprovedFunctionLevelBasedOrderer(contract *dto.ContractDTO) *improvedFunctionLevelBasedOrderer {
// 	return &improvedFunctionLevelBasedOrderer{contract}
// }

// func (o *improvedFunctionLevelBasedOrderer) OrderTransactions(transactions []*dto.TransactionDTO) {
// 	o.assingEnergy(transactions)
// 	sort.SliceStable(transactions, func(i, j int) bool {
// 		return transactions[i].Energy > transactions[i].Energy
// 	})
// }

// func (o *improvedFunctionLevelBasedOrderer) assingEnergy(transactions []*dto.TransactionDTO) {
// 	var minDist = math.MaxFloat64
// 	var maxDist float64 = 0
// 	var BIG_NUMBER float64 = 999999999

// 	for _, transaction := range transactions {
// 		var deltaMin = float64(transaction.DeltaMinDistance)

// 		if deltaMin != BIG_NUMBER {
// 			if deltaMin < minDist {
// 				minDist = deltaMin
// 			}
// 			if deltaMin > maxDist {
// 				maxDist = deltaMin
// 			}
// 		}
// 	}

// 	for _, transaction := range transactions {
// 		var deltaMin = float64(transaction.DeltaMinDistance)

// 		if deltaMin == minDist {
// 			if minDist == maxDist {
// 				transaction.Energy = 1
// 			} else {
// 				transaction.Energy = maxDist - minDist
// 			}
// 		} else {
// 			transaction.Energy = (maxDist - minDist) / (deltaMin - minDist)
// 		}
// 	}
// }

package fuzz

import (
	"math"
	"sort"

	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/dominikbraun/graph"
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

func (o *improvedFunctionLevelBasedOrderer) generateOriginalGraph() graph.Graph[string, string] {
	g := graph.New(graph.StringHash, graph.Directed())

	for key := range o.contract.CFG.Graph {
		_ = g.AddVertex(key)
	}

	for key, value := range o.contract.CFG.Graph {
		for _, v := range value {
			_ = g.AddEdge(key, v)
		}
	}

	return g
}

func (o *improvedFunctionLevelBasedOrderer) assingEnergy(transactions []*dto.TransactionDTO) {
	var minDist = math.MaxFloat64
	var maxDist float64 = 0
	var BIG_NUMBER float64 = 999999999

	g := o.generateOriginalGraph()
	targetBlocks := FindBlocksContainingTargetInstructions(o.contract.CFG, targetInstructions)

	for _, transaction := range transactions {
		var distanceMap = make(map[string]float64)
		var qDistance float64 = 1

		for _, source := range transaction.ExecutedInstructions {
			distanceMap[source] = ComputeDistance(g, source, targetBlocks)
		}

		for _, distance := range distanceMap {
			if distance < BIG_NUMBER {
				if distance < minDist {
					minDist = distance
				}
				if distance > maxDist {
					maxDist = distance
				}
				transaction.Distance += distance
				qDistance += 1
			}
		}

		if transaction.Distance == 0.0 {
			transaction.Distance = BIG_NUMBER
		}

		transaction.Distance = float64(int(transaction.Distance / qDistance))
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

func ComputeDistance(g graph.Graph[string, string], source string, targetBlocks []string) float64 {
	var distance = ZERO

	for _, target := range targetBlocks {
		k, err := graph.ShortestPath(g, source, target)
		if err != nil {
			continue
		}

		distance += float64(len(k))
	}

	if distance == ZERO {
		return math.MaxFloat64
	}

	return distance / float64(len(targetBlocks))
}
