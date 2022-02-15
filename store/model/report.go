package model

type Report struct {
	AccountId int   `json:"account"`
	Score     []int `json:"score"`
}
