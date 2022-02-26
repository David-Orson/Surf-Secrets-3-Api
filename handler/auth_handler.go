package handler

import (
	"net/http"

	"github.com/David-Orson/Surf-Secrets-3-Api/store/model"
	"github.com/David-Orson/Surf-Secrets-3-Api/validation"
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
		v := validation.Validation{}
		v.AddError("Invalid email or password")
		respond(w, v, http.StatusBadRequest)
		return
	}

	respond(w, token, http.StatusOK)
}

func signup(w http.ResponseWriter, r *http.Request) {
	var account model.Account
	readBytes(r, &account)

	v, err := validation.ValidateAccount(&account)
	if err != nil {
		respondMsg(w, "Error: Could not create account", http.StatusBadRequest)
		return
	}
	if !v.IsValid() {
		respond(w, v, http.StatusBadRequest)
		return
	}

	err = s.Account().Create(&account)
	if err != nil {
		respondMsg(w, "Error: Could not create account", http.StatusBadRequest)
		return
	}

	createdAccount, err := s.Account().Get(account.Username)
	if err != nil {
		respondMsg(w, "Error: Could not get newly created account", http.StatusBadRequest)
		return
	}

	respond(w, createdAccount, http.StatusCreated)
}
