package test

import (
	"encoding/json"
	"log"

	"github.com/dogefuzz/dogefuzz/mapper"
	"github.com/dogefuzz/dogefuzz/pkg/geth"
	"github.com/dogefuzz/dogefuzz/repo"
	"go.uber.org/zap"
)

type TestEnv struct {
	logger            *zap.Logger
	contractMapper    mapper.ContractMapper
	transactionMapper mapper.TransactionMapper
	taskMapper        mapper.TaskMapper
	taskRepo          repo.TaskRepo
	contractRepo      repo.ContractRepo
	transactionRepo   repo.TransactionRepo
	deployer          geth.Deployer
}

func NewTestEnv(
	contractMapper mapper.ContractMapper,
	transactionMapper mapper.TransactionMapper,
	taskMapper mapper.TaskMapper,
	taskRepo repo.TaskRepo,
	contractRepo repo.ContractRepo,
	transactionRepo repo.TransactionRepo,
	deployer geth.Deployer,
) *TestEnv {
	return &TestEnv{
		contractMapper: contractMapper,
		contractRepo:   contractRepo,
		deployer:       deployer,
	}
}

func (e *TestEnv) ContractMapper() mapper.ContractMapper {
	return e.contractMapper
}

func (e *TestEnv) TransactionMapper() mapper.TransactionMapper {
	return e.transactionMapper
}

func (e *TestEnv) TaskMapper() mapper.TaskMapper {
	return e.taskMapper
}

func (e *TestEnv) TaskRepo() repo.TaskRepo {
	return e.taskRepo
}

func (e *TestEnv) TransactionRepo() repo.TransactionRepo {
	return e.transactionRepo
}

func (e *TestEnv) ContractRepo() repo.ContractRepo {
	return e.contractRepo
}

func (e *TestEnv) Logger() *zap.Logger {
	if e.logger == nil {
		logger, err := initLogger()
		if err != nil {
			log.Panicf("Error while loading zap logger: %s", err)
			return nil
		}

		e.logger = logger
	}
	return e.logger
}

func (e *TestEnv) Deployer() geth.Deployer {
	return e.deployer
}

func initLogger() (*zap.Logger, error) {
	rawJSON := []byte(`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": ["stdout", "/tmp/logs"],
		"errorOutputPaths": ["stderr"],
		"encoderConfig": {
			"messageKey": "message",
			"levelKey": "level",
			"levelEncoder": "lowercase"
		}
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		return nil, err
	}
	l, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return l, nil
}