package src

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"gorm.io/gorm"

	"github.com/gorilla/mux"
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

type FinanceValue struct {
	Raw float32 `json:"raw"`
	Fmt string  `json:"fmt"`
}

type StockPrice struct {
	RegularMarketPrice         FinanceValue `json:"regularMarketPrice"`
	RegularMarketChangePercent FinanceValue `json:"regularMarketChangePercent"`
}

type QuoteType struct {
	LongName string `json:"longName"`
	Symbol   string `json:"symbol"`
}

type StockPriceResponse struct {
	StockId   uint
	Price     StockPrice `json:"price"`
	QuoteType QuoteType  `json:"quoteType"`
}

func GetTrendingStocks(w http.ResponseWriter, r *http.Request) {
	url := "https://apidojo-yahoo-finance-v1.p.rapidapi.com/market/get-trending-tickers"
	client := http.Client{}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("x-rapidapi-key", os.Getenv("YAHOO_API_KEY"))
	res, _ := client.Do(req)
	var resData FinanceResponse

	json.NewDecoder(res.Body).Decode(&resData)
	json.NewEncoder(w).Encode(resData.Finance.Result[0].Quotes)
}

func GetWatchlist(db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(mux.Vars(r)["id"])
		var user User
		db.Find(&user, userId)

		var stocks []Stock
		db.Model(&user).Association("Stocks").Find(&stocks)

		stockPrices := make(chan StockPriceResponse, len(stocks))
		for i := 0; i < len(stocks); i++ {
			go getStockPrice(stocks[i].Symbol, stockPrices, stocks[i].ID)
		}

		var result []StockPriceResponse
		for i := 0; i < len(stocks); i++ {
			result = append(result, <-stockPrices)
		}

		json.NewEncoder(w).Encode(result)
	}
}

func getStockPrice(symbol string, ch chan StockPriceResponse, stockId uint) {
	url := "https://apidojo-yahoo-finance-v1.p.rapidapi.com/stock/v2/get-profile?symbol=" + symbol
	client := http.Client{}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("x-rapidapi-key", os.Getenv("YAHOO_API_KEY"))
	res, _ := client.Do(req)
	var resData StockPriceResponse

	json.NewDecoder(res.Body).Decode(&resData)
	resData.StockId = stockId
	ch <- resData
}
