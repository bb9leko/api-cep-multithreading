package cep

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Address struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"` // ViaCEP
	Bairro     string `json:"bairro"`     // ViaCEP
	Localidade string `json:"localidade"` // ViaCEP
	Uf         string `json:"uf"`         // ViaCEP
	// Campos BrasilAPI:
	Street       string `json:"street"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	Api          string
}

func FetchFromBrasilAPI(ctx context.Context, cep string, ch chan<- Address, errCh chan<- error) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errCh <- fmt.Errorf("BrasilAPI: %w", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		errCh <- fmt.Errorf("BrasilAPI: status %d", resp.StatusCode)
		return
	}
	var addr Address
	if err := json.Unmarshal(body, &addr); err != nil {
		errCh <- fmt.Errorf("BrasilAPI: %w", err)
		return
	}
	// Ajuste para preencher os campos comuns
	addr.Logradouro = addr.Street
	addr.Bairro = addr.Neighborhood
	addr.Localidade = addr.City
	addr.Uf = addr.State
	addr.Api = "BrasilAPI"
	ch <- addr
}

func FetchFromViaCEP(ctx context.Context, cep string, ch chan<- Address, errCh chan<- error) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errCh <- fmt.Errorf("ViaCEP: %w", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		errCh <- fmt.Errorf("ViaCEP: status %d", resp.StatusCode)
		return
	}
	var addr Address
	if err := json.Unmarshal(body, &addr); err != nil {
		errCh <- fmt.Errorf("ViaCEP: %w", err)
		return
	}
	addr.Api = "ViaCEP"
	ch <- addr
}
