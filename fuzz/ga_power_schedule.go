package fuzz

import (
	"math"
	"strings"

	"github.com/dogefuzz/dogefuzz/config"
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
	// it will take all seeds
	for idx := 0; idx < len(transactions); idx++ {
		seeds = append(seeds, transactions[idx].Inputs)
	}

	// select seeds to crossover
	var selectedSeeds [][]string
	for range seeds {
		selectedSeeds = append(selectedSeeds, rouletteWheelSelection(seeds))
	}

	// do crossover
	var crossoverSeeds [][]string
	if len(selectedSeeds) > 2 {
		for i := 0; i < s.cfg.FuzzerConfig.SeedsSize*2; i++ {
			seed1, seed2 := crossover(chooseSeedsToCrossover(selectedSeeds))
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

func rouletteWheelSelection(seedsList [][]string) []string {
	rnd := common.RandomFloat64()
	lengthList := len(seedsList)
	var seedsSlice [][]string

	if rnd >= 0 && rnd < FIRST_INTERVAL {
		seedsSlice = seedsList[0:maxIndex(lengthList, FIRST_RANGE)]
	} else if rnd >= FIRST_INTERVAL && rnd < SECOND_INTERVAL {
		seedsSlice = seedsList[minIndex(lengthList, FIRST_RANGE):maxIndex(lengthList, SECOND_RANGE)]
	} else {
		// [minIndex:len(seedsList)]
		seedsSlice = seedsList[minIndex(lengthList, SECOND_RANGE):]
	}

	return common.RandomChoice(seedsSlice)
}

func chooseSeedsToCrossover(seedsList [][]string) ([]string, []string) {
	return common.RandomChoice(seedsList), common.RandomChoice(seedsList)
}

// metodos a serem testados podem ter mais de um parametro de entrada
// uma seed representa uma lista de paramentros de entrada
func crossover(seed1, seed2 []string) ([]string, []string) {
	var crossoverSeed1 []string
	var crossoverSeed2 []string

	if len(seed1) == len(seed2) {
		for i := range seed1 {
			smallest := smallestSize(seed1[i], seed2[i])
			// analisar a necessidade desta verificacao
			if smallest == 0 {
				crossoverSeed1 = seed1
				crossoverSeed2 = seed2
				//	funcionando como esperado
			} else if smallest == 1 {
				var tempSeed1 string
				var tempSeed2 string
				for j := 0; j < smallest; j++ {
					tempSeed1 = seed1[i][j:j+1] + seed2[i][j:j+1]
					tempSeed2 = seed2[i][j:j+1] + seed1[i][j:j+1]
				}
				crossoverSeed1 = append(crossoverSeed1, tempSeed1)
				crossoverSeed2 = append(crossoverSeed2, tempSeed2)
			} else {
				var tempSeed1 string
				var tempSeed2 string
				for j := 0; j < smallest; j++ {

					if j%2 == 0 {
						tempSeed1 += seed2[i][j : j+1]
						tempSeed2 += seed1[i][j : j+1]
					} else {
						tempSeed1 += seed1[i][j : j+1]
						tempSeed2 += seed2[i][j : j+1]
					}
				}
				crossoverSeed1 = append(crossoverSeed1, tempSeed1)
				crossoverSeed2 = append(crossoverSeed2, tempSeed2)
			}
		}
	}

	return crossoverSeed1, crossoverSeed2
}

func smallestSize(seed1, seed2 string) int {
	if len(seed1) >= len(seed2) {
		return len(seed2)
	}

	return len(seed1)
}

func maxIndex(lengthList int, factor float64) int {
	return int(math.Ceil(float64(lengthList) * factor))
}

func minIndex(lengthList int, factor float64) int {
	return int(math.Floor(float64(lengthList) * factor))
}
