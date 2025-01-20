package usecase

import (
	"context"
	"errors"
	"fmt"

	// "fmt"
	// "math/rand"
	"CryptoParser/config"
	"CryptoParser/pkg/binanceParser"

	"CryptoParser/internal/entities"
	"CryptoParser/internal/repository/postgres"
	"time"

	// protos "CryptoParser/pkg/proto/gen/go"
	// jwt "github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type Usecase struct {
	cfg     *config.ConfigModel
	log     *zap.Logger
	Repo    *postgres.Repository
	BClient *binanceParser.Client
	ctx     context.Context
}

const (
	clientRole    = `CLIENT`
	adminRole     = `ADMIN`
	SuccessStatus = `Success`
	BadStatus     = `Fail`
)

// TODO: убрать статусы. Возравращается только ошибка.

func NewUsecase(logger *zap.Logger, Repo *postgres.Repository, cfg *config.ConfigModel, ctx context.Context) (*Usecase, error) {
	return &Usecase{
		cfg:     cfg,
		log:     logger,
		Repo:    Repo,
		BClient: binanceParser.NewClient(),
		ctx:     ctx,
	}, nil
}

func (uc *Usecase) OnStart(_ context.Context) error {
	go uc.startBinanceDeamon()
	return nil
}

func (uc *Usecase) OnStop(_ context.Context) error {
	return nil
}

func (uc *Usecase) CurrencyAdd(ctx context.Context, currencyPair string) error {
	exist, err := uc.Repo.IsCurrencyPairExists(ctx, currencyPair)
	if err != nil {
		uc.log.Error("fail to check currency pair exists", zap.Error(err))
		return err
	}
	if !exist {
		isValid, err := uc.BClient.IsTickerValid(ctx, currencyPair)
		if err != nil {
			uc.log.Error("can not valid currency pair", zap.Error(err))
			return err
		}
		if !isValid {
			uc.log.Error("currency pair unvalid", zap.Error(err))
			return err
		}
		if err := uc.Repo.AddCurrencyPair(ctx, currencyPair); err != nil {
			uc.log.Error("fail to add currency pair", zap.Error(err))
			return err
		}
	}
	return nil
}

func (uc *Usecase) CurrencyRemove(ctx context.Context, currencyPair string) error {
	exist, err := uc.Repo.IsCurrencyPairExists(ctx, currencyPair)
	if err != nil {
		uc.log.Error("fail to check currency pair exists", zap.Error(err))
		return err
	}
	if exist {
		if err := uc.Repo.RemoveCurrencyPair(ctx, currencyPair); err != nil {
			uc.log.Error("fail to remove currency pair", zap.Error(err))
			return err
		}
	}
	return nil
}

func (uc *Usecase) CurrencyPrice(ctx context.Context, currencyPair string, timestamp int64) (string, error) {
	exist, err := uc.Repo.IsCurrencyPairExistsInRecords(ctx, currencyPair)
	if err != nil {
		uc.log.Error("fail to check currency pair exists", zap.Error(err))
		return "", err
	}
	if !exist {
		uc.log.Error("currency pair does not exists", zap.Error(err))
		return "", errors.New("currency pair does not exists")
	}
	price, err := uc.Repo.GetCurrencyPrice(ctx, currencyPair, timestamp)
	if err != nil {
		uc.log.Error("can not get price", zap.Error(err))
		return "", err
	}
	return price, nil
}

func (uc *Usecase) startBinanceDeamon() {
	for now := range time.Tick(5 * time.Second) {
		tickers, err := uc.Repo.GetAllTickers(uc.ctx)
		if err != nil {
			uc.log.Error("Can not get all tickers", zap.Error(err))
		}
		for _, ticker := range tickers {
			record := &binanceParser.Record{
				Ticker: binanceParser.Ticker{
					ID:     ticker.ID,
					Ticker: ticker.Ticker,
				},
			}
			if err := uc.BClient.GetTickerPrice(uc.ctx, record); err != nil {
				uc.log.Error(fmt.Sprintf("Can not get price: %s", ticker.Ticker), zap.Error(err))
				continue
			}
			if err := uc.Repo.CreateRecord(
				uc.ctx,
				&entities.Record{
					TickerID:  record.Ticker.ID,
					Timestamp: record.Timestamp,
					Price:     record.Price,
				},
			); err != nil {
				uc.log.Error(fmt.Sprintf("Can not save price for %s at %d", ticker.Ticker, record.Timestamp), zap.Error(err))
			}
		}
		uc.log.Info(fmt.Sprintf("Get price by tickers at %d done", now.Unix()))
	}
}
