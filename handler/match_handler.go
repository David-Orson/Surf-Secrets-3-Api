package handler

import (
	"net/http"

	"github.com/David-Orson/Surf-Secrets-3-Api/store/model"
)

func matchRoutes() {
	// GET
	router.HandleFunc("/match", auth(getAllMatches, "*")).Methods("GET")
	router.HandleFunc("/match/user", auth(getMatchesByAccount, "*")).Methods("GET")
	router.HandleFunc("/match/disputed", auth(getDisputesByAccount, "*")).Methods("GET")
	router.HandleFunc("/match/{id}", auth(getMatch, "*")).Methods("GET")

	// POST
	router.HandleFunc("/match", auth(acceptMatch, "*")).Methods("POST")

	// PUT
	router.HandleFunc("/match", auth(reportMatch, "*")).Methods("PUT")
}

func getMatch(w http.ResponseWriter, r *http.Request, authModel model.Auth) {
	id, err := getId(r, "id")
	if err != nil {
		respondMsg(w, "Error: {id} is not an integer", http.StatusBadRequest)
		return
	}

	var match model.Match

	match, err = s.Match().Get(id)

	if err != nil {
		respondMsg(w, "Error: Could not get match", http.StatusBadRequest)
		return
	}

	respond(w, match, http.StatusOK)
}

func getAllMatches(w http.ResponseWriter, r *http.Request, authModel model.Auth) {
	var matches []model.Match

	matches, err := s.Match().GetAll()

	if err != nil {
		respondMsg(w, "Error: Could not get all matches", http.StatusBadRequest)
		return
	}

	respond(w, matches, http.StatusOK)
}

func getMatchesByAccount(w http.ResponseWriter, r *http.Request, authModel model.Auth) {
	var matches []model.Match

	matches, err := s.Match().GetByAccount(authModel.AccountId)

	if err != nil {
		respondMsg(w, "Error: Could not get matches by account", http.StatusBadRequest)
		return
	}

	respond(w, matches, http.StatusOK)
}

func getDisputesByAccount(w http.ResponseWriter, r *http.Request, authModel model.Auth) {
	var matches []model.Match

	matches, err := s.Match().GetDisputesByAccount(authModel.AccountId)

	if err != nil {
		respondMsg(w, "Error: Could not get disputes by account", http.StatusBadRequest)
		return
	}

	respond(w, matches, http.StatusOK)
}

func acceptMatch(w http.ResponseWriter, r *http.Request, authModel model.Auth) {

}

func reportMatch(w http.ResponseWriter, r *http.Request, authModel model.Auth) {

}
