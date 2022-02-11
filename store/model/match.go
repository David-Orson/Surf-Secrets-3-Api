package model

import "time"

type Match struct {
	Id         int       `json:"id"`
	Team0      []int     `json:"teamA"`
	Team1      []int     `json:"teamB"`
	TeamSize   int       `json:"teamSize"`
	Time       time.Time `json:"time"`
	Maps       []int     `json:"maps"`
	Result0    []int     `json:"resultA"`
	Result1    []int     `json:"resultB"`
	IsDisputed bool      `json:"isDisputed"`
	Result     int       `json:"result"`
}
