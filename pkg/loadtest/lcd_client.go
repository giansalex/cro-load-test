package loadtest

import (
	"encoding/json"
	"net/http"
)

type LcdClient struct {
	client  *http.Client
	baseUrl string
}

func NewLcdClient(client *http.Client, baseUrl string) *LcdClient {
	return &LcdClient{client, baseUrl}
}

func (lcd *LcdClient) Account(address string) (*AccountResponse, error) {
	resp, err := lcd.client.Get(lcd.baseUrl + "auth/accounts/" + address)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var data AccountResponse
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (lcd *LcdClient) Balances(address string) (*BalancesResponse, error) {
	resp, err := lcd.client.Get(lcd.baseUrl + "bank/balances/" + address)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var data BalancesResponse
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

type AccountResponse struct {
	Height string `json:"height"`
	Result struct {
		Type  string `json:"type"`
		Value struct {
			Address   string `json:"address"`
			PublicKey struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"public_key"`
			AccountNumber string `json:"account_number"`
			Sequence      string `json:"sequence"`
		} `json:"value"`
	} `json:"result"`
}

type BalancesResponse struct {
	Height string `json:"height"`
	Result []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"result"`
}
