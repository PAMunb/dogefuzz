package distance

import (
	"math"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dominikbraun/graph"
)

func ComputeDistanceMap(cfg common.CFG, targetInstructions []string) common.DistanceMap {
	targetBlocks := findBlocksContainingTargetInstructions(cfg, targetInstructions)
	reversedCFG := cfg.GetReverseGraph()
	distanceMap := make(common.DistanceMap)
	for cfgBlockPC := range cfg.Graph {
		distanceMap[cfgBlockPC] = make(map[string]uint32)
		for _, targetBlock := range targetBlocks {
			distanceMap[cfgBlockPC][targetBlock] = math.MaxUint32
		}
	}

	for _, targetBlock := range targetBlocks {
		distanceMap[targetBlock][targetBlock] = 0
		queue := []string{targetBlock}
		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]
			for _, edge := range reversedCFG[current] {
				if distanceMap[edge][targetBlock] == math.MaxUint32 {
					distanceMap[edge][targetBlock] = distanceMap[current][targetBlock] + 1
					queue = append(queue, edge)
				}
			}
		}
	}

	return distanceMap
}

// breath first search over cfg graph searching for target instructions
func ComputeTargetInstructionsFrequency(cfg common.CFG, targetInstructions []string) uint64 {
	var count uint64 = 0
	for _, block := range cfg.Blocks {
		for _, instr := range block.Instructions {
			if common.Contains(targetInstructions, instr) {
				count++
			}
		}
	}
	return count
}

func ComputeDistance(cfg common.CFG, instructionsExecutedInTransaction []string, targetInstructions []string) float64 {
	g := generateOriginalGraph(cfg)
	targetBlocks := findBlocksContainingTargetInstructions(cfg, targetInstructions)

	return computeDistanceToInstruction(g, instructionsExecutedInTransaction, targetBlocks)
}

func computeDistanceToInstruction(g graph.Graph[string, string], instructionsExecutedInTransaction []string, targetBlocks []string) float64 {
	var distance = 0.0
	var minDistance = math.MaxFloat64

	for _, target := range targetBlocks {
		for _, source := range instructionsExecutedInTransaction {
			k, err := graph.ShortestPath(g, source, target)
			if err != nil {
				continue
			}

			if float64(len(k)) < minDistance {
				minDistance = float64(len(k))
			}
		}

		distance += minDistance
	}

	if distance == 0.0 {
		return math.MaxFloat64
	}

	return distance / float64(len(targetBlocks))
}

func generateOriginalGraph(cfg common.CFG) graph.Graph[string, string] {
	g := graph.New(graph.StringHash, graph.Directed())

	for key := range cfg.Graph {
		_ = g.AddVertex(key)
	}

	for key, value := range cfg.Graph {
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
