package handler

import (
	"net/http"

	"github.com/David-Orson/Surf-Secrets-3-Api/store/model"
)

func finderRoutes() {
	// GET
	router.HandleFunc("/finder-posts", getAllFinderPosts).Methods("GET")

	// POST
	router.HandleFunc("/finder-post", createFinderPost).Methods("POST")

}

func getAllFinderPosts(w http.ResponseWriter, r *http.Request) {
	var finderPosts []model.FinderPost

	finderPosts, err := s.Finder().GetAllPosts()

	if err != nil {
		respondMsg(w, "Error: Could not get match finder posts", http.StatusBadRequest)
		return
	}

	respond(w, finderPosts, http.StatusOK)
}

func createFinderPost(w http.ResponseWriter, r *http.Request) {
	var finderPost model.FinderPost
	var err error

	readBytes(r, &finderPost)

	err = s.Finder().CreatePost(&finderPost)

	if err != nil {
		respondMsg(w, "Error: Could not create match finder post", http.StatusBadRequest)
		return
	}

	respond(w, finderPost, http.StatusCreated)
}
