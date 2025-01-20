package binanceParser

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	client *http.Client
}

func NewClient() *Client {
	return &Client{
		client: http.DefaultClient,
	}
}

const domainURL = `https://api.binance.com/`
const endpointTickerPrice = `api/v3/ticker/price?symbol=`

func (c *Client) GetTickerPrice(ctx context.Context, record *Record) error {
	record.Timestamp = time.Now().Unix()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s%s", domainURL, endpointTickerPrice, record.Ticker.Ticker), nil)
	if err != nil {
		return err
	}

	request.Header.Set("Accept", "application/vnd.yclients.v2+json")
	request.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("client does not registered")
	}

	var tickerPriceResp TickerPriceResponse
	err = json.NewDecoder(resp.Body).Decode(&tickerPriceResp)
	if err != nil {
		return err
	}
	record.Price = tickerPriceResp.Price
	return nil
}

func (c *Client) IsTickerValid(ctx context.Context, currecnyPair string) (bool, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s%s", domainURL, endpointTickerPrice, currecnyPair), nil)
	if err != nil {
		return false, err
	}

	request.Header.Set("Accept", "application/vnd.yclients.v2+json")
	request.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(request)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != 200 {
		return false, errors.New("currency pair is unvalid")
	}
	return true, nil
}
