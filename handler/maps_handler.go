package handler

import (
	"net/http"

	"github.com/David-Orson/Surf-Secrets-3-Api/store/model"
)

func mapRoutes() {
	// GET
	router.HandleFunc("/map", auth(getAllMaps, "*")).Methods("GET")
}

func getAllMaps(w http.ResponseWriter, r *http.Request, authModel model.Auth) {
	var maps []model.Map

	maps, err := s.Map().GetAll()

	if err != nil {
		respondMsg(w, "Error: Could not get all maps", http.StatusBadRequest)
		return
	}

	respond(w, maps, http.StatusOK)
}
