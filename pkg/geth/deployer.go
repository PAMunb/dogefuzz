package geth

import (
	"context"
	"math/big"
	"strings"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

type Deployer interface {
	Deploy(ctx context.Context, contract *common.Contract, args ...interface{}) (gethcommon.Address, *types.Transaction, error)
}

type deployer struct {
	client *ethclient.Client
	wallet Wallet
	cfg    config.GethConfig
}

func NewDeployer(cfg config.GethConfig, wallet Wallet) (*deployer, error) {
	client, err := ethclient.Dial(cfg.NodeAddress)
	if err != nil {
		return nil, err
	}

	return &deployer{client, wallet, cfg}, nil
}

func (d *deployer) Deploy(ctx context.Context, contract *common.Contract, args ...interface{}) (gethcommon.Address, *types.Transaction, error) {
	parsedABI, err := abi.JSON(strings.NewReader(contract.AbiDefinition))
	if err != nil {
		return gethcommon.Address{}, nil, err
	}

	nonce, err := d.client.PendingNonceAt(ctx, d.wallet.GetAddress())
	if err != nil {
		return gethcommon.Address{}, nil, err
	}

	gasPrice, err := d.client.SuggestGasPrice(ctx)
	if err != nil {
		return gethcommon.Address{}, nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(d.wallet.GetPrivateKey(), big.NewInt(d.cfg.ChainID))
	if err != nil {
		return gethcommon.Address{}, nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	address, tx, _, err := bind.DeployContract(auth, parsedABI, gethcommon.FromHex(contract.CompiledCode), d.client, args...)
	if err != nil {
		return gethcommon.Address{}, nil, err
	}
	return address, tx, nil
}
