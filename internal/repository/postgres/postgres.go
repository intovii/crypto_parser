package postgres

import (
	"CryptoParser/config"
	"CryptoParser/internal/entities"
	"CryptoParser/pkg/binanceParser"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Repository struct {
	ctx context.Context
	log *zap.Logger
	cfg *config.ConfigModel
	DB  *pgxpool.Pool
}

func NewRepository(log *zap.Logger, cfg *config.ConfigModel, ctx context.Context) (*Repository, error) {
	return &Repository{
		ctx: ctx,
		log: log,
		cfg: cfg,
	}, nil
}

func (r *Repository) OnStart(_ context.Context) error {
	connectionUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		r.cfg.Postgres.Host,
		r.cfg.Postgres.Port,
		r.cfg.Postgres.User,
		r.cfg.Postgres.Password,
		r.cfg.Postgres.DBName,
		r.cfg.Postgres.SSLMode)
	pool, err := pgxpool.Connect(r.ctx, connectionUrl)
	if err != nil {
		return err
	}
	r.DB = pool
	return nil
}

func (r *Repository) OnStop(_ context.Context) error {
	r.DB.Close()
	return nil
}

const queryIsCurrencyPairExist = `
	SELECT EXISTS (SELECT id
               FROM tickers
               WHERE ticker = $1 
			   )
`

func (r *Repository) IsCurrencyPairExists(ctx context.Context, currencyPair string) (bool, error) {
	var res bool
	err := r.DB.QueryRow(ctx, queryIsCurrencyPairExist, currencyPair).Scan(&res)
	if err != nil {
		r.log.Error("fail to check currency pair exists", zap.Error(err))
		return false, err
	}
	return res, nil
}

const queryIsCurrencyPairExistsInRecords = `
	SELECT EXISTS (SELECT 
        FROM records
        JOIN tickers ON records.ticker_id = tickers.id
        WHERE tickers.ticker = $1)
`

func (r *Repository) IsCurrencyPairExistsInRecords(ctx context.Context, currencyPair string) (bool, error) {
	var res bool
	err := r.DB.QueryRow(ctx, queryIsCurrencyPairExistsInRecords, currencyPair).Scan(&res)
	if err != nil {
		r.log.Error("fail to check currency pair exists", zap.Error(err))
		return false, err
	}
	return res, nil
}

const queryGetAllTickers = `
	SELECT id, ticker
	FROM tickers 
`

func (r *Repository) GetAllTickers(ctx context.Context) ([]binanceParser.Ticker, error) {
	rows, err := r.DB.Query(ctx, queryGetAllTickers)
	if err != nil {
		r.log.Error("fail to get all tickers", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	tickers := make([]binanceParser.Ticker, 0, 8)
	for rows.Next() {
		var ticker binanceParser.Ticker
		er := rows.Scan(
			&ticker.ID,
			&ticker.Ticker,
		)
		if er != nil {
			return nil, er
		}
		tickers = append(tickers, ticker)
	}
	if err := rows.Err(); err != nil {
		r.log.Error("fail to scan all tickers", zap.Error(err))
		return nil, err
	}
	return tickers, nil
}

const queryAddCurrencyPair = `
	INSERT INTO tickers (ticker) 
	VALUES ($1)
`

func (r *Repository) AddCurrencyPair(ctx context.Context, currencyPair string) error {
	_, err := r.DB.Exec(ctx, queryAddCurrencyPair, currencyPair)
	if err != nil {
		r.log.Error("fail to add currency pair", zap.Error(err))
		return err
	}
	return nil
}

const queryRemoveCurrencyPair = `
	DELETE FROM tickers
	WHERE ticker = $1
`

func (r *Repository) RemoveCurrencyPair(ctx context.Context, currencyPair string) error {
	_, err := r.DB.Exec(ctx, queryRemoveCurrencyPair, currencyPair)
	if err != nil {
		r.log.Error("fail to remove currency pair", zap.Error(err))
		return err
	}
	return nil
}

const queryGetCurrencyPrice = `
		SELECT price
        FROM records
        JOIN tickers ON records.ticker_id = tickers.id
        WHERE tickers.ticker = $1
        ORDER BY ABS(records.timestamp - $2)
        LIMIT 1;
`

func (r *Repository) GetCurrencyPrice(ctx context.Context, currencyPair string, timestamp int64) (string, error) {
	var price string
	err := r.DB.QueryRow(ctx, queryGetCurrencyPrice, currencyPair, timestamp).Scan(&price)
	if err != nil {
		r.log.Error("fail to get price by currency pair and timestamp", zap.Error(err))
		return "", err
	}
	return price, nil
}

const queryCreateRecord = `
	INSERT INTO records 
	(ticker_id, timestamp, price)
	VALUES ($1, $2, $3)
`

func (r *Repository) CreateRecord(ctx context.Context, record *entities.Record) error {
	_, err := r.DB.Exec(ctx, queryCreateRecord, record.TickerID, record.Timestamp, record.Price)
	if err != nil {
		r.log.Error("fail to create record by tickerID, timestamp and price", zap.Error(err))
		return err
	}
	return nil
}
