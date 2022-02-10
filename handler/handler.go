package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/David-Orson/Surf-Secrets-3-Api/store"
	"github.com/gorilla/mux"
)

var router *mux.Router
var s store.Store

func InitRouter(r *mux.Router, mainStore store.Store) {
	s = mainStore
	router = r

	authRoutes()
}

func readBytes(r *http.Request, model interface{}) string {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		log.Println("e0009: Could not read bytes")

	}
	jsonString := string(bodyBytes)
	json.Unmarshal([]byte(jsonString), model)

	return jsonString
}

func respond(w http.ResponseWriter, model interface{}, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(model)
}
