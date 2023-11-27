package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const DEFAULT_BR_DATE_FORMAT = "02/01/2006 15:04:05"
const DEFAULT_TIMESTAMP_MASK = "2006-01-02T15:04:05Z"

func createDbConnection() *sql.DB {
	db, err := sql.Open("sqlite3", "./cotacao.db")
	if err != nil {
		panic(err)
	}
	return db
}

func initDatabaseTable() {
	db := createDbConnection()
	defer db.Close()

	sqlCreateTableStmt := `
	CREATE TABLE IF NOT EXISTS cotacao (
		id INTEGER not null primary key, 
		code TEXT, 
		codein TEXT, 
		name TEXT, 
		high TEXT, 
		low TEXT, 
		var_bid TEXT, 
		pct_change TEXT, 
		bid TEXT,
		ask TEXT, 
		timestamp_api TEXT, 
		create_date_api DATE,
		date_inserted TIMESTAMP DEFAULT CURRENT_TIMESTAMP);
	`
	_, err := db.Exec(sqlCreateTableStmt)
	if err != nil {
		panic(err)
	}
}

func insertExchangeRate(ctx context.Context, db *sql.DB, c *CotacaoUsdToBrl) error {
	// Utilizar uma goroutine para executar a query de inserção
	ch := make(chan error, 1)

	go func() {
		// Usar o contexto para limitar o tempo de execução
		select {
		case <-ctx.Done():
			ch <- fmt.Errorf("timeout reached")
		default:
			// Executar a query de inserção
			stmt, err := db.Prepare("INSERT INTO cotacao(code, codein, name, high, low, var_bid, pct_change, bid, ask, timestamp_api, create_date_api) VALUES(?,?,?,?,?,?,?,?,?,?,?)")
			if err != nil {
				ch <- err
			}
			defer stmt.Close()

			_, err = stmt.Exec(
				c.Code,
				c.Codein,
				c.Name,
				c.High,
				c.Low,
				c.VarBid,
				c.PctChange,
				c.Bid,
				c.Ask,
				c.Timestamp,
				c.CreateDate,
			)
			ch <- err
		}
	}()

	// Aguardar o resultado da goroutine
	select {
	case <-ctx.Done():
		return fmt.Errorf("cotacao persistance cancelled. timeout reached")
	case err := <-ch:
		return err
	}
}

func getAllCurrencyExchangerRates(db *sql.DB) ([]*CotacaoUsdToBrl, error) {
	rows, err := db.Query("SELECT * FROM cotacao")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allCurrencyRates []*CotacaoUsdToBrl
	for rows.Next() {
		var currencyRate CotacaoUsdToBrl
		err := rows.Scan(
			&currencyRate.Id,
			&currencyRate.Code,
			&currencyRate.Codein,
			&currencyRate.Name,
			&currencyRate.High,
			&currencyRate.Low,
			&currencyRate.VarBid,
			&currencyRate.PctChange,
			&currencyRate.Bid,
			&currencyRate.Ask,
			&currencyRate.Timestamp,
			&currencyRate.CreateDate,
			&currencyRate.DateInsertedFormatted,
		)
		if err != nil {
			return nil, err
		}

		// Fazer o parse do timestamp para o formato desejado
		t, err := time.Parse(DEFAULT_TIMESTAMP_MASK, currencyRate.DateInsertedFormatted)
		if err != nil {
			fmt.Println("Erro ao fazer o parse do timestamp:", err)
		}

		// Formatar a data no formato desejado
		currencyRate.DateInsertedFormatted = t.Format(DEFAULT_BR_DATE_FORMAT)
		allCurrencyRates = append(allCurrencyRates, &currencyRate)
	}
	return allCurrencyRates, nil
}
