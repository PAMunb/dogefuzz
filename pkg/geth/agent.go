package geth

import (
	"context"
	"math/big"
	"strings"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Agent interface {
	Send(ctx context.Context, contract *common.Contract, functionName string, args ...interface{}) (string, error)
}

type agent struct {
	client *ethclient.Client
	wallet Wallet
	cfg    config.GethConfig
}

func NewAgent(cfg config.GethConfig) (*agent, error) {
	wallet, err := NewWalletFromPrivateKeyHex(cfg.AgentPrivateKeyHex)
	if err != nil {
		return nil, err
	}

	client, err := ethclient.Dial(cfg.NodeAddress)
	if err != nil {
		return nil, err
	}

	return &agent{client, wallet, cfg}, nil
}

func (d *agent) Send(ctx context.Context, contract *common.Contract, functionName string, args ...interface{}) (string, error) {
	parsedABI, err := abi.JSON(strings.NewReader(contract.AbiDefinition))
	if err != nil {
		return "", err
	}

	boundContract := bind.NewBoundContract(gethcommon.HexToAddress(contract.Address), parsedABI, d.client, d.client, d.client)

	nonce, err := d.client.PendingNonceAt(ctx, d.wallet.GetAddress())
	if err != nil {
		return "", err
	}

	gasPrice, err := d.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(d.wallet.GetPrivateKey(), big.NewInt(d.cfg.ChainID))
	if err != nil {
		return "", err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(999999)
	auth.GasPrice = gasPrice

	tx, err := boundContract.Transact(auth, functionName, args...)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}
