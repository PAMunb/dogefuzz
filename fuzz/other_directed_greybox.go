package fuzz

import (
	"strings"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

// This is temporally, it is like 'directed_greybox' but with another strategy to order seeds
type otherDirectedGreyboxFuzzer struct {
	powerSchedule interfaces.PowerSchedule

	solidityService interfaces.SolidityService
	functionService interfaces.FunctionService
	contractService interfaces.ContractService
}

func NewOtherDirectedGreyboxFuzzer(e env) *otherDirectedGreyboxFuzzer {
	return &otherDirectedGreyboxFuzzer{
		powerSchedule:   e.PowerSchedule(),
		solidityService: e.SolidityService(),
		functionService: e.FunctionService(),
		contractService: e.ContractService(),
	}
}

func (f *otherDirectedGreyboxFuzzer) GenerateInput(functionId string) ([]interface{}, error) {
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

	chosenSeeds := common.RandomChoice(seedsList)

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
