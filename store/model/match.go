package model

import "time"

type Match struct {
	Id         int       `json:"id"`
	Team0      []int     `json:"team0"`
	Team1      []int     `json:"team1"`
	TeamSize   int       `json:"teamSize"`
	Time       time.Time `json:"time"`
	Maps       []Map     `json:"maps"`
	Result0    []int     `json:"result0"`
	Result1    []int     `json:"result1"`
	IsDisputed bool      `json:"isDisputed"`
	Result     int       `json:"result"`
}
