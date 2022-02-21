package handler

import (
	"net/http"

	"github.com/David-Orson/Surf-Secrets-3-Api/store/model"
)

func matchRoutes() {
	// GET
	router.HandleFunc("/matches", getAllMatches).Methods("GET")
	router.HandleFunc("/matches/user/{id}", getMatchesByAccount).Methods("GET")
	router.HandleFunc("/match/disputed", auth(getDisputesByAccount, "*")).Methods("GET")
	router.HandleFunc("/match/{id}", getMatch).Methods("GET")

	// POST
	router.HandleFunc("/match/{id}", auth(acceptMatch, "*")).Methods("POST")

	// PUT
	router.HandleFunc("/match/{id}", auth(reportMatch, "*")).Methods("PUT")
}

func getMatch(w http.ResponseWriter, r *http.Request) {
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

func getAllMatches(w http.ResponseWriter, r *http.Request) {
	var matches []model.Match

	matches, err := s.Match().GetAll()

	if err != nil {
		respondMsg(w, "Error: Could not get all matches", http.StatusBadRequest)
		return
	}

	respond(w, matches, http.StatusOK)
}

func getMatchesByAccount(w http.ResponseWriter, r *http.Request) {
	var matches []model.Match

	id, err := getId(r, "id")
	if err != nil {
		respondMsg(w, "Error: {id} is not an integer", http.StatusBadRequest)
		return
	}

	matches, err = s.Match().GetByAccount(id)
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
	var resultArray [2]int
	var team1 [1]int

	finderPost, err = s.Finder().GetPost(id)
	if err != nil {
		respondMsg(w, "Error: Could not get match finder post", http.StatusBadRequest)
		return
	}

	resultArray[0] = 0
	resultArray[1] = 0

	team1[0] = authModel.AccountId

	match.Team0 = finderPost.Team
	match.Team1 = team1[:]
	match.TeamSize = len(finderPost.Team)
	match.Time = finderPost.Time
	match.Maps = finderPost.Maps
	match.Result0 = resultArray[:]
	match.Result1 = resultArray[:]
	match.IsDisputed = false
	match.Result = 4

	err = s.Match().Create(&match)
	if err != nil {
		respondMsg(w, "Error: Could not accept match", http.StatusBadRequest)
		return
	}

	err = s.Finder().SetAccepted(id)
	if err != nil {
		respondMsg(w, "Error: Match finder post was not removed from the finder", http.StatusBadRequest)
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
	isInTeam0 := false

	readBytes(r, &report)
	report.AccountId = authModel.AccountId

	match, err = s.Match().Get(id)
	if err != nil {
		respondMsg(w, "Error: Could not get match", http.StatusBadRequest)
		return
	}

	for _, account := range match.Team0 {
		if report.AccountId == account {
			match.Result0 = report.Score
			isInTeam0 = true
		}
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
				s.Account().IncrementWin(match.Team0[0])
				s.Account().IncrementLoss(match.Team1[0])
			} else {
				match.Result = 1
				s.Account().IncrementWin(match.Team1[0])
				s.Account().IncrementLoss(match.Team0[0])
			}
		} else {
			match.IsDisputed = true
			match.Result = 3
			s.Account().IncrementDispute(match.Team0[0])
			s.Account().IncrementDispute(match.Team1[0])

		}
	}

	err = s.Match().Update(&match)
	if err != nil {
		respondMsg(w, "Error: Could not report match", http.StatusBadRequest)
		return
	}

	respond(w, match, http.StatusCreated)
}
