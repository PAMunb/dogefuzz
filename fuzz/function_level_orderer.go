package fuzz

import (
	"sort"

	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/dominikbraun/graph"
)

type functionLevelBasedOrderer struct {
	contract *dto.ContractDTO
}

func newFunctionLevelBasedOrderer(contract *dto.ContractDTO) *functionLevelBasedOrderer {
	return &functionLevelBasedOrderer{contract}
}

func (o *functionLevelBasedOrderer) OrderTransactions(transactions []*dto.TransactionDTO) {
	o.assingEnergy(transactions)

	sort.SliceStable(transactions, func(i, j int) bool {
		return transactions[i].Energy > transactions[j].Energy
	})
}

func (o *functionLevelBasedOrderer) assingEnergy(transactions []*dto.TransactionDTO) {
	//g := o.generateOriginalGraph()
	//targetBlocks := FindBlocksContainingTargetInstructions(o.contract.CFG, targetInstructions)

	//for _, transaction := range transactions {
	//	distance := ComputeDistance(g, transaction, targetBlocks)
	//	meanDistance := MeanMinDistances(distance)
	//	transaction.Energy = ComputeEnergy(meanDistance)
	//}
}

func (o *functionLevelBasedOrderer) generateOriginalGraph() graph.Graph[string, string] {
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

// func (o *functionLevelBasedOrderer) generateGraphExecuted(transaction *dto.TransactionDTO) common.CFG {
// 	var graphExecuted common.CFG
// 	graphExecuted.Graph = make(map[string][]string)

// 	// create executed graph
// 	for key, value := range o.contract.CFG.Graph {
// 		for _, executedInstruction := range transaction.ExecutedInstructions {
// 			if key == executedInstruction {
// 				graphExecuted.Graph[key] = value
// 				break
// 			}
// 		}
// 	}

// 	// remove useless reference
// 	for key, values := range graphExecuted.Graph {
// 		var graphContains []string
// 		for _, v := range values {
// 			for block := range graphExecuted.Graph {
// 				if block == v {
// 					graphContains = append(graphContains, v)
// 				}
// 			}
// 		}

// 		graphExecuted.Graph[key] = graphContains
// 	}

// 	return graphExecuted
// }
