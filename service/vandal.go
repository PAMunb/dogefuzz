package service

import (
	"context"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/dogefuzz/dogefuzz/pkg/vandal"
)

type vandalService struct {
	cfg    config.VandalConfig
	client interfaces.VandalClient
}

func NewVandalService(e Env) *vandalService {
	return &vandalService{
		cfg: e.Config().VandalConfig,
	}
}

func (s *vandalService) GetCFG(ctx context.Context, contract *common.Contract) (*common.CFG, error) {
	client := s.getClient()

	blocks, _, err := client.Decompile(ctx, contract.RuntimeBytecode, contract.Name)
	if err != nil {
		return nil, err
	}

	cfg := new(common.CFG)
	cfg.Graph = make(map[string][]string)
	cfg.Blocks = make(map[string]common.CFGBlock)
	cfg.Instructions = make(map[string]string)
	for _, block := range blocks {
		instructions := make(map[string]string)
		pcs := make([]string, 0)
		for pc, instruction := range block.Instructions {
			instructions[pc] = instruction.Op
			cfg.Instructions[pc] = instruction.Op
			pcs = append(pcs, pc)
		}

		cfg.Graph[block.PC] = block.Successors
		cfg.Blocks[block.PC] = common.CFGBlock{
			InitialPC:       block.PC,
			Instructions:    instructions,
			InstructionsPCs: pcs,
		}
	}

	return cfg, nil
}

func (s *vandalService) getClient() interfaces.VandalClient {
	if s.client == nil {
		s.client = vandal.NewVandalClient(s.cfg.Endpoint, s.cfg.CfgTimeout)
	}
	return s.client
}
