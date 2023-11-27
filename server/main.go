package main

import "net/http"

func main() {
	initDatabaseTable()
	http.HandleFunc("/cotacao", CheckUsdToBrlAndReturn)
	http.HandleFunc("/list", ListAllCurrencyExchangerRates)
	http.ListenAndServe(":8080", nil)
}