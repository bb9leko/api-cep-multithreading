package main

import (
	"context"
	"fmt"
	"os"
	"time"

	cep "github.com/bb9leko/api-cep-multithreading/service"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <CEP>")
		return
	}
	cepInput := os.Args[1]

	ch := make(chan cep.Address, 2)
	errCh := make(chan error, 2)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	go cep.FetchFromBrasilAPI(ctx, cepInput, ch, errCh)
	go cep.FetchFromViaCEP(ctx, cepInput, ch, errCh)

	select {
	case addr := <-ch:
		fmt.Printf("Resposta recebida da %s:\n", addr.Api)
		fmt.Printf("CEP: %s\nLogradouro: %s\nBairro: %s\nCidade: %s\nUF: %s\n", addr.Cep, addr.Logradouro, addr.Bairro, addr.Localidade, addr.Uf)
	case err := <-errCh:
		select {
		case addr := <-ch:
			fmt.Printf("Resposta recebida da %s:\n", addr.Api)
			fmt.Printf("CEP: %s\nLogradouro: %s\nBairro: %s\nCidade: %s\nUF: %s\n", addr.Cep, addr.Logradouro, addr.Bairro, addr.Localidade, addr.Uf)
		case <-ctx.Done():
			fmt.Println("Erro: Timeout de 1 segundo atingido.")
		}
		fmt.Println("Erro:", err)
	case <-ctx.Done():
		fmt.Println("Erro: Timeout de 1 segundo atingido.")
	}
}
