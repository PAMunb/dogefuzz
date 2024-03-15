package fuzz

import (
	"math/rand"
	"strings"
	"time"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

// consts to define what interval of the seeds slice will selected
const FIRST_INTERVAL float64 = 0.6
const SECOND_INTERVAL float64 = 0.8

// consts to define what range interval will used to prioritize some seeds
const FIRST_RANGE float64 = 0.4
const SECOND_RANGE float64 = 0.7

type geneticAlgorithmFuzzer struct {
	powerSchedule   interfaces.PowerSchedule
	solidityService interfaces.SolidityService
	functionService interfaces.FunctionService
	contractService interfaces.ContractService
}

func NewGeneticAlgorithmFuzzer(e env) *geneticAlgorithmFuzzer {
	return &geneticAlgorithmFuzzer{
		powerSchedule:   e.PowerSchedule(),
		solidityService: e.SolidityService(),
		functionService: e.FunctionService(),
		contractService: e.ContractService(),
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

	seedsList, err := f.powerSchedule.RequestSeeds(functionId, common.DISTANCE_COVERAGE_BASED_STRATEGY)
	if err != nil {
		return nil, err
	}

	var newSeedsList [][]interface{}
	for range seedsList {
		newSeedsList = append(newSeedsList, rouletteWheelSelection(seedsList))
	}

	chosenSeeds := common.RandomChoice(newSeedsList)

	inputs := make([]interface{}, len(method.Inputs))
	for inputsIdx, inputDefinition := range method.Inputs {
		handler, err := f.solidityService.GetTypeHandlerWithContext(inputDefinition.Type)
		if err != nil {
			return nil, err
		}
		handler.SetValue(chosenSeeds[inputsIdx])
		mutationFunction := common.RandomChoice(handler.GetMutators())
		mutationFunction()
		inputs[inputsIdx] = handler.GetValue()
	}

	return inputs, nil
}

func rouletteWheelSelection(seedsList [][]interface{}) []interface{} {
	rand.Seed(time.Now().UnixNano())
	rnd := rand.Float64()

	if rnd >= 0 || rnd < FIRST_INTERVAL {
		slice := seedsList[0:int(float64(len(seedsList))*FIRST_RANGE)]

		return common.RandomChoice(slice)
	} else if rnd >= FIRST_INTERVAL || rnd < SECOND_INTERVAL {
		slice := seedsList[int(float64(len(seedsList))*FIRST_RANGE):int(float64(len(seedsList))*SECOND_RANGE)]

		return common.RandomChoice(slice)
	} else {
		slice := seedsList[int(float64(len(seedsList))*SECOND_RANGE):]

		return common.RandomChoice(slice)
	}

}
