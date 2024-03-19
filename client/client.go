package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {		
		panic(err)
		
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {		
		panic(err)	
	}
	defer resp.Body.Close()
	

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)			
	}
	
	if resp.StatusCode != http.StatusOK {
		println(string(body))
		log.Fatalf("Erro no servidor: %s", resp.Status)		
	} else {
		println(string(body))
	}

	err = os.WriteFile("cotacao.txt", body, 0644)
	if err != nil {
		log.Fatalf("Ocorreu um erro ao tentar salvar a contacao no arquivo: %v", err)
	}
	
}