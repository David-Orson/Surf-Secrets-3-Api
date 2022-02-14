package model

import "time"

type FinderPost struct {
	Id       int       `json:"id"`
	Team     []int     `json:"team"`
	TeamSize int       `json:"teamSize"`
	Time     time.Time `json:"time"`
	Maps     []Map     `json:"maps"`
}
