package handler

import (
	"net/http"

	"github.com/David-Orson/Surf-Secrets-3-Api/store/model"
)

func accountRoutes() {
	// GET
	router.HandleFunc("/account/{id}", auth(getAccount, "*")).Methods("GET")

	// POST
	router.HandleFunc("/account", auth(createAccount, "canCreateAccount")).Methods("POST")

	// PUT
	router.HandleFunc("/account", auth(updateAccount, "canUpdateAccount")).Methods("PUT")

	// DELETE
	router.HandleFunc("/account/{id}", auth(deleteAccount, "canDeleteAccount")).Methods("DELETE")
}

func getAccount(w http.ResponseWriter, r *http.Request, authModel model.Auth) {
	id, err := getId(r, "id")
	if err != nil {
		respondMsg(w, "Error: {id} is not an integer", http.StatusBadRequest)
		return
	}

	var account model.Account

	account, err = s.Account().Get(id)

	if err != nil {
		respondMsg(w, "Error: Could not get account", http.StatusBadRequest)
		return
	}

	respond(w, account, http.StatusOK)
}

func createAccount(w http.ResponseWriter, r *http.Request, authModel model.Auth) {
	var account model.Account
	var err error

	readBytes(r, &account)

	err = s.Account().Create(&account)

	if err != nil {
		respondMsg(w, "Error: Could not create account", http.StatusBadRequest)
		return
	}

	respond(w, account, http.StatusCreated)
}

func updateAccount(w http.ResponseWriter, r *http.Request, authModel model.Auth) {
	var account model.Account
	var err error

	readBytes(r, &account)

	err = s.Account().Update(&account)

	if err != nil {
		respondMsg(w, "Error: Could not update account", http.StatusBadRequest)
		return
	}

	respond(w, account, http.StatusCreated)
}

func deleteAccount(w http.ResponseWriter, r *http.Request, authModel model.Auth) {
	id, err := getId(r, "id")
	if err != nil {
		respondMsg(w, "Error: {id} is not an integer", http.StatusBadRequest)
		return
	}

	err = s.Account().Delete(id)
	if err != nil {
		respondMsg(w, "Error: Could not delete account", http.StatusBadRequest)
		return
	}

	respondMsg(w, "Successfully deleted account", http.StatusOK)
}
