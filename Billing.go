package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

var GoLoAdsToken = TOKEN

type MoneyRequest struct {
	Token       string  `json:"token"`
	AccountID   int     `json:"account_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

func sendMoneyToUser(user_id int, money_am float64) (*http.Response, error) {
	var moneyRequest = MoneyRequest{
		Token:       GoLoAdsToken,
		AccountID:   user_id,
		Amount:      money_am,
		Description: "Вывод средств со счёта GoloAds на счет пользователя",
	}

	postBody, _ := json.Marshal(moneyRequest)
	responseBody := bytes.NewBuffer(postBody)
	response, err := http.Post("https://bank.goto.msk.ru/api/send", "application/json", responseBody)

	return response, err
}
