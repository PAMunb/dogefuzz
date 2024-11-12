package fuzz

import (
	"math"
	"sort"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
	"github.com/dominikbraun/graph"
)

// this is need to be constant
var targetInstructions = []string{"CALL", "SELFDESTRUCT", "CALLCODE", "DELEGATECALL"}

type altDistanceBasedOrderer struct {
	contract *dto.ContractDTO
}

func newAltDistanceBasedOrderer(contract *dto.ContractDTO) *altDistanceBasedOrderer {
	return &altDistanceBasedOrderer{contract}
}

func (o *altDistanceBasedOrderer) OrderTransactions(transactions []*dto.TransactionDTO) {
	g := o.generateOriginalGraph()

	sort.SliceStable(transactions, func(i, j int) bool {
		return o.computeScore(g, transactions[i]) < o.computeScore(g, transactions[j])
	})
}

func (o *altDistanceBasedOrderer) computeScore(g graph.Graph[string, string], transaction *dto.TransactionDTO) float64 {
	targetBlocks := findBlocksContainingTargetInstructions(o.contract.CFG, targetInstructions)
	var score float64 = 0
	var mapDistanceToTargetInstruction = make(map[string]float64)

	for _, target := range targetBlocks {
		mapDistanceToTargetInstruction[target] = computeShortestPathToTarget(g, transaction, target)
	}

	for _, distance := range mapDistanceToTargetInstruction {
		score += 1 / distance
	}

	return score
}

func computeShortestPathToTarget(g graph.Graph[string, string], transaction *dto.TransactionDTO, target string) float64 {
	var distance = math.MaxFloat64

	for _, source := range transaction.ExecutedInstructions {
		k, err := graph.ShortestPath(g, source, target)
		if err != nil {
			continue
		}

		distance = math.Min(distance, float64(len(k)))
	}

	return distance
}

func (o *altDistanceBasedOrderer) generateOriginalGraph() graph.Graph[string, string] {
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

func findBlocksContainingTargetInstructions(cfg common.CFG, targetInstructions []string) []string {
	targetBlocks := make([]string, 0)
	for blockPC, block := range cfg.Blocks {
		for _, instr := range block.Instructions {
			if common.Contains(targetInstructions, instr) {
				targetBlocks = append(targetBlocks, blockPC)
				break
			}
		}
	}

	return targetBlocks
}
