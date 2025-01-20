package binanceParser

type TickerPriceResponse struct{
	Symbol string `json:"symbol"`
	Price string `json:"price"`
}

type Ticker struct {
	ID 			int64
	Ticker 		string
}

type Record struct {
	ID 			uint64
	Ticker		Ticker
	Timestamp 	int64
	Price 		string
}