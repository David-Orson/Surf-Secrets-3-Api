package model

type Account struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Win      int    `json:"win"`
	Loss     int    `json:"loss"`
	Disputes int    `json:"disputes"`
	SteamId  int    `json:"steam_id"`
}
