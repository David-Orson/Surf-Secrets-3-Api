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
	router.HandleFunc("/match/{id}", auth(acceptMatch, "*")).Methods("POST")

	// PUT
	router.HandleFunc("/match/{id}", auth(reportMatch, "*")).Methods("PUT")
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
	id, err := getId(r, "id")
	if err != nil {
		respondMsg(w, "Error: {id} is not an integer", http.StatusBadRequest)
		return
	}

	var finderPost model.FinderPost
	var match model.Match
	var resultArray []int

	readBytes(r, &match)

	finderPost, err = s.Finder().GetPost(id)

	resultArray[0] = 0
	resultArray[1] = 0

	match.Team0 = finderPost.Team
	match.TeamSize = finderPost.TeamSize
	match.Time = finderPost.Time
	match.Maps = finderPost.Maps
	match.Result0 = resultArray
	match.Result1 = resultArray
	match.IsDisputed = false
	match.Result = 4

	err = s.Match().Create(&match)

	if err != nil {
		respondMsg(w, "Error: Could not accept match", http.StatusBadRequest)
		return
	}

	respond(w, match, http.StatusCreated)
}

func reportMatch(w http.ResponseWriter, r *http.Request, authModel model.Auth) {
	id, err := getId(r, "id")
	if err != nil {
		respondMsg(w, "Error: {id} is not an integer", http.StatusBadRequest)
		return
	}

	var report model.Report
	var match model.Match
	var isInTeam0 bool
	readBytes(r, &report)

	match, err = s.Match().Get(id)

	for _, account := range match.Team0 {
		if report.AccountId == account {
			match.Result0 = report.Score
		}
		isInTeam0 = true
	}

	if !isInTeam0 {
		for _, account := range match.Team1 {
			if report.AccountId == account {
				match.Result1 = report.Score
			}
		}
	}

	if (match.Result0[0] > 0 || match.Result0[1] > 0) &&
		(match.Result1[0] > 0 || match.Result1[1] > 0) {
		if match.Result0[0] == match.Result1[0] && match.Result0[1] == match.Result1[1] {
			if match.Result0[0] > match.Result0[1] {
				match.Result = 0
			} else {
				match.Result = 1
			}
		} else {
			match.IsDisputed = true
			match.Result = 3
		}
	}

	err = s.Match().Update(&match)

	if err != nil {
		respondMsg(w, "Error: Could not report match", http.StatusBadRequest)
		return
	}

	respond(w, match, http.StatusCreated)
}
