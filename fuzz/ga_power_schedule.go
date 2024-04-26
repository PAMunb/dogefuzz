package fuzz

import (
	"strings"

	"github.com/dogefuzz/dogefuzz/config"
	ga "github.com/dogefuzz/dogefuzz/fuzz/genetic_operators"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type geneticAlgorithmPowerSchedule struct {
	cfg                *config.Config
	transactionService interfaces.TransactionService
	solidityService    interfaces.SolidityService
	functionService    interfaces.FunctionService
	contractService    interfaces.ContractService
}

func NewGeneticAlgorithmPowerSchedule(e env) *geneticAlgorithmPowerSchedule {
	return &geneticAlgorithmPowerSchedule{
		cfg:                e.Config(),
		transactionService: e.TransactionService(),
		solidityService:    e.SolidityService(),
		functionService:    e.FunctionService(),
		contractService:    e.ContractService(),
	}
}

func (s *geneticAlgorithmPowerSchedule) RequestSeeds(functionId string, strategy common.PowerScheduleStrategy) ([][]interface{}, error) {
	function, err := s.functionService.Get(functionId)
	if err != nil {
		return nil, err
	}

	contract, err := s.contractService.Get(function.ContractId)
	if err != nil {
		return nil, err
	}

	abiDefinition, err := abi.JSON(strings.NewReader(contract.AbiDefinition))
	if err != nil {
		return nil, err
	}
	method := abiDefinition.Methods[function.Name]

	transactions, err := s.transactionService.FindDoneTransactionsByFunctionIdAndOrderByTimestamp(functionId, int64(s.cfg.FuzzerConfig.SeedsSize)*2)
	if err != nil {
		return nil, err
	}

	orderer := buildOrderer(strategy, contract)
	orderer.OrderTransactions(transactions)

	seeds := make([][]string, 0)
	// it will takes all seeds
	for idx := 0; idx < len(transactions); idx++ {
		seeds = append(seeds, transactions[idx].Inputs)
	}

	// select seeds to crossover
	selectedSeeds := make([][]string, 0)
	for range seeds {
		selectedSeeds = append(selectedSeeds, ga.Selection(seeds))
	}

	// do crossover
	crossoverSeeds := make([][]string, 0)
	if len(selectedSeeds) > 2 {
		for i := 0; i < len(selectedSeeds); i++ {
			seed1, seed2 := ga.Crossover(chooseSeedsToCrossover(selectedSeeds))
			crossoverSeeds = append(crossoverSeeds, seed1)
			crossoverSeeds = append(crossoverSeeds, seed2)
		}
	}

	deserializedSeeds, err := s.deserializeSeedsList(method, crossoverSeeds)
	if err != nil {
		return nil, err
	}

	if len(crossoverSeeds) < s.cfg.FuzzerConfig.SeedsSize {
		deserializedSeeds, err = s.completeSeedsWithPreConfiguredSeeds(method, deserializedSeeds, s.cfg.FuzzerConfig.SeedsSize-len(crossoverSeeds))
		if err != nil {
			return nil, err
		}
	}

	return deserializedSeeds, nil
}

func (s *geneticAlgorithmPowerSchedule) completeSeedsWithPreConfiguredSeeds(method abi.Method, seeds [][]interface{}, seedsAmountToBeAdded int) ([][]interface{}, error) {
	result := make([][]interface{}, len(seeds)+seedsAmountToBeAdded)
	copy(result, seeds)
	for icr := 0; icr < int(seedsAmountToBeAdded); icr++ {
		functionSeeds := make([]interface{}, len(method.Inputs))
		for inputsIdx, input := range method.Inputs {
			handler, err := s.solidityService.GetTypeHandlerWithContext(input.Type)
			if err != nil {
				return nil, err
			}

			err = handler.LoadSeedsAndChooseOneRandomly(s.cfg.FuzzerConfig.Seeds)
			if err != nil {
				return nil, err
			}

			functionSeeds[inputsIdx] = handler.GetValue()
		}
		result[icr+len(seeds)] = functionSeeds
	}
	return result, nil
}

func (s *geneticAlgorithmPowerSchedule) deserializeSeedsList(method abi.Method, seedsList [][]string) ([][]interface{}, error) {
	result := make([][]interface{}, len(seedsList))
	for seedsListIdx, seeds := range seedsList {
		deserializedSeeds := make([]interface{}, len(seeds))
		for inputsIdx, inputDefinition := range method.Inputs {
			if len(seeds) <= inputsIdx {
				return nil, ErrSeedsListInvalid
			}

			handler, err := s.solidityService.GetTypeHandlerWithContext(inputDefinition.Type)
			if err != nil {
				return nil, err
			}

			err = handler.Deserialize(seeds[inputsIdx])
			if err != nil {
				return nil, err
			}
			deserializedSeeds[inputsIdx] = handler.GetValue()
		}
		result[seedsListIdx] = deserializedSeeds
	}
	return result, nil
}

func chooseSeedsToCrossover(seedsList [][]string) ([]string, []string) {
	return common.RandomChoice(seedsList), common.RandomChoice(seedsList)
}
