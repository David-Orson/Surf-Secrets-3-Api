package handler

import (
	"log"
	"net/http"

	"github.com/David-Orson/Surf-Secrets-3-Api/store/model"
)

func authRoutes() {
	router.HandleFunc("/auth", login).Methods("POST")
}

func login(w http.ResponseWriter, r *http.Request) {
	var account model.Account
	readBytes(r, &account)

	token, err := s.Auth().Login(&account)
	if err != nil {
		log.Println("login error, ***validation needed")
		return
	}

	respond(w, token, http.StatusOK)
}
