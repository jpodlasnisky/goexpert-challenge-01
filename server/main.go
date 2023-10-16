package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type USDBRL struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type Cotacao struct {
	USDBRL USDBRL
}

func main() {
	log.Println("Starting server...")
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", cotacaoHandler)

	http.ListenAndServe(":8080", mux)

}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Calling economia.awesomeapi.com.br")
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	client := http.Client{}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/last/USD-BRL", nil)
	if err != nil {
		log.Fatalf("%v", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("%v", err)
	}

	var cotacao Cotacao
	if err = json.Unmarshal(body, &cotacao); err != nil {
		log.Fatalf("%v", err)
	}

	db, err := NewDBConnection()
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cotacao(
		id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		code VARCHAR(1024),
		codein VARCHAR(1024) not null,
		name VARCHAR(1024) not null,
		high VARCHAR(1024) not null,
		low VARCHAR(1024) not null,
		varBid VARCHAR(1024) not null,
		pctChange VARCHAR(1024) not null,
		bid VARCHAR(1024) not null,
		ask VARCHAR(1024) not null,
		timestamp VARCHAR(1024) not null,
		create_date VARCHAR(1024) not null
	  )`)
	if err != nil {
		log.Fatalf("%v", err)
	}

	contextDB, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err = insertCotacao(db, cotacao, contextDB)
	if err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(w).Encode(cotacao.USDBRL.Bid)
	if err != nil {
		log.Fatal(err)
	}
}

func NewDBConnection() (*sql.DB, error) {
	log.Println("Connecting to database...")
	db, err := sql.Open("sqlite3", "./goexpert.db")
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func insertCotacao(db *sql.DB, cotacao Cotacao, context context.Context) error {
	log.Println("Inserting in the table cotacao...")
	stmt, err := db.PrepareContext(context, `INSERT INTO cotacao(code, codein, name, 
		high, low, varBid, pctChange, Bid, Ask, timestamp, create_date) 
		VALUES(?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(context, cotacao.USDBRL.Code, cotacao.USDBRL.Codein, cotacao.USDBRL.Name, cotacao.USDBRL.High, cotacao.USDBRL.Low,
		cotacao.USDBRL.VarBid, cotacao.USDBRL.PctChange, cotacao.USDBRL.Bid, cotacao.USDBRL.Ask, cotacao.USDBRL.Timestamp,
		cotacao.USDBRL.CreateDate)
	if err != nil {
		return err
	}
	return nil

}
