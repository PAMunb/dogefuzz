package fuzz

import (
	"math"
	"sort"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
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
	graphToTargetInstructions := o.generatePathsToTargetInstructions(targetInstructions)

	sort.SliceStable(transactions, func(i, j int) bool {
		graphExecuted := o.generateGraphExecuted(transactions[i])
		compareGraphExecutedAndGraphToTargetInstructions(graphExecuted, graphToTargetInstructions)
		computeDistance(graphExecuted, targetInstructions)
		return o.computeScore(transactions[i]) < o.computeScore(transactions[j])
	})
}

func compareGraphExecutedAndGraphToTargetInstructions(graphExecuted common.CFG, graphToTargetInstructions []map[string][]string) map[int]int {
	var countMatchKey = make(map[int]int)
	for i, path := range graphToTargetInstructions {
		for key := range path {
			_, ok := graphExecuted.Graph[key]
			if ok {
				countMatchKey[i]++
			}
		}
	}

	return countMatchKey
}

// comparar esse grafo com os paths ate instrucao critica
func (o *altDistanceBasedOrderer) generateGraphExecuted(transaction *dto.TransactionDTO) common.CFG {
	var graphExecuted common.CFG
	graphExecuted.Graph = make(map[string][]string)

	for key, value := range o.contract.CFG.Graph {
		for _, executedInstruction := range transaction.ExecutedInstructions {
			if key == executedInstruction {
				graphExecuted.Graph[key] = value
				break
			}
		}
	}

	return graphExecuted

	// preciso (?)
	// graphExecuted.Instructions = make(map[string]string)

	// graphExecuted.Blocks = make(map[string]common.CFGBlock)
	// graphExecuted.Instructions = o.contract.CFG.Instructions

	/*
		* o grafo representa as associacoes entre blocos (que contem um conjunto de instrucoes)
		* exemplo: bloco A -> B -> D
						   -> C
		* a instrucao (key) que identifica o bloco no grafo
		* esta contida no conjunto de intrucoes de um blobo
		* nesse caso, quando uma instrucao da match com a chave (key)
		* no grafo, significa que esse bloco foi executado
		*
		* Entao, basicamente, o algoritmo abaixo identifca os blocos que foram executados
		*
		* Falta:
		* 1) computar a distancia percorrida, os blocos executados que estao relacionados
		* caminho mais longo, seria  distancia percorrida
		*
		* 2) montar os caminhos ate instrucoes criticas
		*
		* 3) comparar a distancia percorrida pelo input com os caminhos ate instrucoes criticas
		* se der match fazer um calculo
		* se nao alcançar avaliar mal o input
		*
		*
		* comparar dois map res1 := reflect.DeepEqual(map_1, map_2)
	*/

}

// deve montar os caminhos ate instrucoes criticas
func (o *altDistanceBasedOrderer) generatePathsToTargetInstructions(targetInstructions []string) []map[string][]string {
	pathsToTargetInstructions := make([]map[string][]string, 0)

	// identifica os blocos que contem instrucoes criticas
	targetBlocks := findBlocksContainingTargetInstructions(o.contract.CFG, targetInstructions)
	reversedGraph := o.contract.CFG.GetReverseGraph()

	for _, block := range targetBlocks {
		for key, value := range reversedGraph {
			if key == block {
				path := make(map[string][]string)
				path[key] = value
				pathsToTargetInstructions = append(pathsToTargetInstructions, path)
				break
			}
		}
	}

	// targetBlock (?)
	// monta um caminhos ate as intrucoes criticas
	// i = 0; pathsToTargetInstructions contem as instrucoes criticas presentes no contrato
	// 		  e suas respectivas referencias
	for _, path := range pathsToTargetInstructions {
		n := 0
		for n < 1 {
			for k := range reversedGraph {
				// if "talvez" desnecessario
				if len(path[k]) > 0 {
					for _, v := range path[k] {
						path[v] = reversedGraph[v]
						pathsToTargetInstructions = append(pathsToTargetInstructions, path)

						// if len(path[v]) == 0 {
						// 	n++
						// }
						// condicao para sair do loop
						// significa que chegou a instrucao de entrada
					}
				} else {
					n++
				}
			}

			// for _, ref := range path {
			// 	if len(ref) == 0 {
			// 		n++
			// 	}
			// }
		}
	}

	return pathsToTargetInstructions
}

func computeTargetInstructionsFrequency(cfg common.CFG, targetInstructions []string) uint64 {
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

func computeDistance(cfg common.CFG, targetInstructions []string) uint64 {
	computeTargetInstructionsFrequency(cfg, targetInstructions)

	cfg.GetReverseGraph()

	return 0
}

func (o *altDistanceBasedOrderer) computeScore(transaction *dto.TransactionDTO) float64 {
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

// preciso (?)
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
