package wallet

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alan1-666/eth-wallet/common/tasks"
	"github.com/alan1-666/eth-wallet/config"
	"github.com/alan1-666/eth-wallet/database"
	"github.com/alan1-666/eth-wallet/wallet/node"
	"github.com/ethereum/go-ethereum/log"
)

type Withdraw struct {
	db             *database.DB
	chainConfig    *config.ChainConfig
	client         node.EthClient
	resourceCtx    context.Context
	resourceCancel context.CancelFunc
	tasks          tasks.Group
}

func NewWithdraw(cfg *config.Config, db *database.DB, client node.EthClient, shutdown context.CancelCauseFunc) (*Withdraw, error) {
	resCtx, resCancel := context.WithCancel(context.Background())
	return &Withdraw{
		db:             db,
		chainConfig:    &cfg.Chain,
		client:         client,
		resourceCtx:    resCtx,
		resourceCancel: resCancel,
		tasks: tasks.Group{HandleCrit: func(err error) {
			shutdown(fmt.Errorf("critical error in deposit: %w", err))
		}},
	}, nil
}

func (w *Withdraw) Close() error {
	var result error
	w.resourceCancel()
	if err := w.tasks.Wait(); err != nil {
		result = errors.Join(result, fmt.Errorf("failed to await deposit %w"), err)
	}
	return nil
}

func (w *Withdraw) Start() error {
	log.Info("start withdraw......")
	tickerWithdrawWorker := time.NewTicker(time.Second * 5)
	w.tasks.Go(func() error {
		for range tickerWithdrawWorker.C {
			log.Info("withdraw work task go")
		}
		return nil
	})
	return nil
}
