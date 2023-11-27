package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const URL_CURRENCY_EXCHANGER_RATE_USD_BRL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

type CotacaoUsdToBrl struct {
	Id                    string `json:"id"`
	Code                  string `json:"code"`
	Codein                string `json:"codein"`
	Name                  string `json:"name"`
	High                  string `json:"high"`
	Low                   string `json:"low"`
	VarBid                string `json:"varBid"`
	PctChange             string `json:"pctChange"`
	Bid                   string `json:"bid"`
	Ask                   string `json:"ask"`
	Timestamp             string `json:"timestamp"`
	CreateDate            string `json:"create_date"`
	DateInsertedFormatted string `json:"date_inserted_formatted"`
}

type Cotacao struct {
	CotacaoUsdToBrl CotacaoUsdToBrl `json:"USDBRL"`
}

type CurrentCurrencyRate struct {
	Bid string `json:"bid"`
}

func NewCurrentCurrencyRate(bid string) *CurrentCurrencyRate {
	return &CurrentCurrencyRate{
		Bid: bid,
	}
}

func CheckUsdToBrlAndReturn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode("Method not allowed")
		return
	}

	currencyRate, err := getCurrentCurrencyRateFromApi()
	if err != nil {
		writeError(w, "Error when fetching API data", err)
		return
	}
	err = persistData(currencyRate)
	if err != nil {
		writeError(w, "Error when persisting data", err)
		return
	}

	currentCurrencyRateToReturn := NewCurrentCurrencyRate(currencyRate.Bid)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(currentCurrencyRateToReturn)
}

func ListAllCurrencyExchangerRates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode("Method not allowed")
		return
	}

	db := createDbConnection()
	defer db.Close()
	currencyRates, err := getAllCurrencyExchangerRates(db)
	if err != nil {
		writeError(w, "Error on fetching cotacoes", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(currencyRates)
}

func getCurrentCurrencyRateFromApi() (*CotacaoUsdToBrl, error) {
	// Cria um contexto com timeout de 200ms
	ctxApi := context.Background()
	ctxApi, cancel := context.WithTimeout(ctxApi, 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctxApi, http.MethodGet, URL_CURRENCY_EXCHANGER_RATE_USD_BRL, nil)
	if err != nil {
		return nil, err
	}

	// Realiza a requisição para a API de cotação
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var currencyRate Cotacao
	err = json.NewDecoder(resp.Body).Decode(&currencyRate)
	if err != nil {
		return nil, err
	}

	return &currencyRate.CotacaoUsdToBrl, nil
}

func persistData(c *CotacaoUsdToBrl) error {
	// Cria um contexto com timeout de 10ms
	ctxDb, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// Persiste a cotação no banco de dados usando o contexto com timeout
	db := createDbConnection()
	defer db.Close()

	err := insertExchangeRate(ctxDb, db, c)
	if err != nil {
		return err
	}

	return nil
}

func writeError(w http.ResponseWriter, prefixMsg string, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	msg := fmt.Sprintf("%s - %s", prefixMsg, err.Error())
	log.Fatalln(msg)
	json.NewEncoder(w).Encode("Internal server error: " + msg)
}
