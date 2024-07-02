package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	cep          = "01153000"
	brazilAPIURL = "https://brasilapi.com.br/api/cep/v1/"
	viaCepURL    = "http://viacep.com.br/ws/"
	timeout      = 1 * time.Second
)

type BrazilAPIAddress struct {
	Street       string `json:"street"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
}

type ViaCepAddress struct {
	Street       string `json:"logradouro"`
	Neighborhood string `json:"bairro"`
	City         string `json:"localidade"`
	State        string `json:"uf"`
}

type Address struct {
	Street       string
	Neighborhood string
	City         string
	State        string
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ch := make(chan result)

	go fetchBrazilAPI(ctx, ch, brazilAPIURL+cep, "BrasilAPI")
	go fetchViaCep(ctx, ch, viaCepURL+cep+"/json/", "ViaCEP")

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

func fetchBrazilAPI(ctx context.Context, ch chan<- result, url string, api string) {
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- result{api: api, err: err}
		return
	}

	var brazilAPIAddress BrazilAPIAddress
	if err := json.Unmarshal(body, &brazilAPIAddress); err != nil {
		ch <- result{api: api, err: err}
		return
	}

	address := Address{
		Street:       brazilAPIAddress.Street,
		Neighborhood: brazilAPIAddress.Neighborhood,
		City:         brazilAPIAddress.City,
		State:        brazilAPIAddress.State,
	}

	ch <- result{api: api, address: address, err: nil}
}

func fetchViaCep(ctx context.Context, ch chan<- result, url string, api string) {
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- result{api: api, err: err}
		return
	}

	var viaCepAddress ViaCepAddress
	if err := json.Unmarshal(body, &viaCepAddress); err != nil {
		ch <- result{api: api, err: err}
		return
	}

	address := Address{
		Street:       viaCepAddress.Street,
		Neighborhood: viaCepAddress.Neighborhood,
		City:         viaCepAddress.City,
		State:        viaCepAddress.State,
	}

	ch <- result{api: api, address: address, err: nil}
}
