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

func main() {
	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
		defer cancel()

		
		req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
		if err != nil {
			http.Error(w, "Ocorreu um erro ao criar a requisição", http.StatusInternalServerError)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "Ocorreu um ao fazer a requisição", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Ocorreu um erro ao ler a resposta", http.StatusInternalServerError)
			return
		}

		var apiResp map[string]struct {
			Bid string `json:"bid"`
		}

		if err := json.Unmarshal(body, &apiResp); err != nil {
			http.Error(w, "Ocorreu um ao fazer o parse", http.StatusInternalServerError)
			return
		}

		
		bid := apiResp["USDBRL"].Bid

		
		dbCtx, dbCancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer dbCancel()

		db, err := sql.Open("sqlite3", "./quotes.sqlite")
		if err != nil {
			http.Error(w, "Ocorreu um erro ao tentar abrir conexao", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		_, err = db.ExecContext(dbCtx, "CREATE TABLE IF NOT EXISTS quotes (id INTEGER PRIMARY KEY, bid TEXT)")
		if err != nil {
			http.Error(w, "Ocorreu um erro ao tentar criar a tabela", http.StatusInternalServerError)
			return
		}

		_, err = db.ExecContext(dbCtx, "INSERT INTO quotes (bid) VALUES (?)", bid)
		if err != nil {
			http.Error(w, "Ocorreu um erro ao tentar inserir a cotação", http.StatusInternalServerError)
			return
		}

		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"bid": bid})
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}