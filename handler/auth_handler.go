package handler

import (
	"log"
	"net/http"

	"github.com/David-Orson/Surf-Secrets-3-Api/store/model"
)

func authRoutes() {
	router.HandleFunc("/auth/login", login).Methods("POST")
	router.HandleFunc("/auth/signup", signup).Methods("POST")
}

func login(w http.ResponseWriter, r *http.Request) {
	var account model.Account
	readBytes(r, &account)

	token, err := s.Auth().Login(&account)
	if err != nil {
		log.Println("e0024: login error, ***validation needed")
		return
	}

	respond(w, token, http.StatusOK)
}

func signup(w http.ResponseWriter, r *http.Request) {
	var account model.Account
	readBytes(r, &account)

	err := s.Account().Create(&account)
	if err != nil {
		respondMsg(w, "Error: Could not create account", http.StatusBadRequest)
		return
	}

	respond(w, account, http.StatusCreated)
}
