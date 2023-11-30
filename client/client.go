package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)


type CurrentCurrencyRate struct {
	Bid string `json:"bid"`
}

func NewCurrentCurrencyRate(bid string) *CurrentCurrencyRate {
	return &CurrentCurrencyRate{
		Bid: bid,
	}
}

const URL_API_CURRENCY_EXCHANGER_RATE_USD_BRL = "http://localhost:8080/cotacao"

func main() {
	currentBid, err := fetchBidFromApi()
	if err != nil {
		log.Fatalln("Error while fetching current bid from API", err)
		return
	}
	log.Println("Current bid fetched from API:", currentBid.Bid)
	err = saveBidInFile(currentBid.Bid)
	if err != nil {
		log.Fatalln("Error while saving bid in file", err)
		return
	}
}

func fetchBidFromApi() (*CurrentCurrencyRate, error) {
	// Cria um contexto com timeout de 300ms
	ctxApi := context.Background()
	ctxApi, cancel := context.WithTimeout(ctxApi, 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctxApi, http.MethodGet, URL_API_CURRENCY_EXCHANGER_RATE_USD_BRL, nil)
	if err != nil {
		return nil, err
	}

	// Realiza a requisição para a API de cotação
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var currentBid CurrentCurrencyRate
	err = json.NewDecoder(resp.Body).Decode(&currentBid)
	if err != nil {
		return nil, err
	}

	return &currentBid, nil
}

func saveBidInFile(bidValue string) error {
	// Cria o arquivo
	f, err := os.Create("cotacao.txt")
	if err != nil {
		return err
	}
	defer f.Close()

	// Escreve no arquivo
	_, err = f.WriteString("Dólar: " + bidValue)
	if err != nil {
		return err
	}

	log.Println("Bid saved in file")
	return nil
}
