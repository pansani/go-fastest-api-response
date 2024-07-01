package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	cep          = "01153000"
	brazilAPIURL = "https://brasilapi.com.br/api/cep/v1/"
	viaCepURL    = "http://viacep.com.br/ws/"
	timeout      = 1 * time.Second
)

type Address struct {
	Street       string `json:"logradouro"`
	Neighborhood string `json:"bairro"`
	City         string `json:"localidade"`
	State        string `json:"uf"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ch := make(chan result)

	go fetch(ctx, ch, brazilAPIURL+cep, "BrasilAPI")
	go fetch(ctx, ch, viaCepURL+cep+"/json/", "ViaCEP")

	select {
	case res := <-ch:
		if res.err != nil {
			fmt.Println("Error:", res.err)
		} else {
			fmt.Printf("Response from %s:\n%+v\n", res.api, res.address)
		}
	case <-ctx.Done():
		fmt.Println("Error: timeout")
	}
}

type result struct {
	api     string
	address Address
	err     error
}

func fetch(ctx context.Context, ch chan<- result, url string, api string) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		ch <- result{api: api, err: err}
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ch <- result{api: api, err: err}
		return
	}
	defer resp.Body.Close()

	var address Address
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		ch <- result{api: api, err: err}
		return
	}

	ch <- result{api: api, address: address, err: nil}
}
