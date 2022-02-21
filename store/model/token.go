package model

type Token struct {
	Token     string `json:"token"`
	AccountId int    `json:"accountId"`
	Username  string `json:"username"`
}
