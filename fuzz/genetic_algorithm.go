package fuzz

import (
	"strings"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

const MUTATION_CHANCE float64 = 0.1

// consts to define what interval of the seeds slice will selected
const FIRST_INTERVAL float64 = 0.6
const SECOND_INTERVAL float64 = 0.8

// consts to define what range interval will used to prioritize some seeds
const FIRST_RANGE float64 = 0.4
const SECOND_RANGE float64 = 0.7

type geneticAlgorithmFuzzer struct {
	geneticAlgorithmPowerSchedule interfaces.GeneticAlgorithmPowerSchedule
	solidityService               interfaces.SolidityService
	functionService               interfaces.FunctionService
	contractService               interfaces.ContractService
}

func NewGeneticAlgorithmFuzzer(e env) *geneticAlgorithmFuzzer {
	return &geneticAlgorithmFuzzer{
		geneticAlgorithmPowerSchedule: e.GeneticAlgorithmPowerSchedule(),
		solidityService:               e.SolidityService(),
		functionService:               e.FunctionService(),
		contractService:               e.ContractService(),
	}
}

func (f *geneticAlgorithmFuzzer) GenerateInput(functionId string) ([]interface{}, error) {
	function, err := f.functionService.Get(functionId)
	if err != nil {
		return nil, err
	}

	contract, err := f.contractService.Get(function.ContractId)
	if err != nil {
		return nil, err
	}

	abiDefinition, err := abi.JSON(strings.NewReader(contract.AbiDefinition))
	if err != nil {
		return nil, err
	}
	method := abiDefinition.Methods[function.Name]

	// evaluate seeds - order a list by strategy
	seedsList, err := f.geneticAlgorithmPowerSchedule.RequestSeeds(functionId, common.DISTANCE_COVERAGE_BASED_STRATEGY)
	if err != nil {
		return nil, err
	}

	chosenSeeds := common.RandomChoice(seedsList)

	inputs := make([]interface{}, len(method.Inputs))
	for inputsIdx, inputDefinition := range method.Inputs {
		rnd := common.RandomFloat64()

		handler, err := f.solidityService.GetTypeHandlerWithContext(inputDefinition.Type)
		if err != nil {
			return nil, err
		}

		handler.SetValue(chosenSeeds[inputsIdx])

		if rnd >= 0 && rnd < MUTATION_CHANCE {
			mutationFunction := common.RandomChoice(handler.GetMutators())
			mutationFunction()
		}

		inputs[inputsIdx] = handler.GetValue()
	}

	return inputs, nil
}
