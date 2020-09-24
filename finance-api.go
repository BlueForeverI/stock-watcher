package main

import (
	"encoding/json"
	"net/http"
	"os"
)

type StockQuote struct {
	Symbol                     string  `json:"symbol"`
	LongName                   string  `json:"longName"`
	RegularMarketChangePercent float32 `json:"regularMarketChangePercent"`
}

type FinanceResult struct {
	Count  int          `json:"count"`
	Quotes []StockQuote `jon:"quotes"`
}

type FinanceResponseData struct {
	Result []FinanceResult `json:"result"`
}

type FinanceResponse struct {
	Finance FinanceResponseData `json:"finance"`
}

func getTrendingStocks(w http.ResponseWriter, r *http.Request) {
	url := "https://apidojo-yahoo-finance-v1.p.rapidapi.com/market/get-trending-tickers"
	spaceClient := http.Client{}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("x-rapidapi-key", os.Getenv("YAHOO_API_KEY"))
	res, _ := spaceClient.Do(req)
	var resData FinanceResponse

	json.NewDecoder(res.Body).Decode(&resData)
	json.NewEncoder(w).Encode(resData.Finance.Result[0].Quotes)
}
