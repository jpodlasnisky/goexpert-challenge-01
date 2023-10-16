package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatalf("%v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer res.Body.Close()
	log.Printf("%v\n", res.Status)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("%v", err)
	}

	var bid string
	err = json.Unmarshal(body, &bid)
	if err != nil {
		log.Fatalf("%v", err)
	}

	writeFile(bid)
}

func writeFile(bid string) {
	testaArquivo()
	file, err := os.OpenFile("./cotacao.txt", os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("%v", err)
	}

	_, err = file.WriteString("Dólar: " + bid + "\n")
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func testaArquivo() {
	_, err := os.OpenFile("./cotacao.txt", os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		os.Create("./cotacao.txt")
	}
}
