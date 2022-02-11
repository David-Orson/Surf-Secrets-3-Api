package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/David-Orson/Surf-Secrets-3-Api/store"
	"github.com/David-Orson/Surf-Secrets-3-Api/store/model"
	"github.com/gorilla/mux"
)

type AuthHandler func(http.ResponseWriter, *http.Request, model.Auth)

var router *mux.Router
var s store.Store

func InitRouter(r *mux.Router, mainStore store.Store) {
	s = mainStore
	router = r

	authRoutes()
}

func getId(r *http.Request, idName string) (int, error) {
	idStr := mux.Vars(r)[idName]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("e0020: {id} is not an integer")
		log.Println(err)
		return 0, err
	}

	return id, nil
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

func respondMsg(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(model.MessageResponse{Message: message})
}

func auth(next AuthHandler, columns ...string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		authModel := authenticate(res, req, columns...)

		next(res, req, authModel)
	}
}

func authenticate(res http.ResponseWriter, req *http.Request, columns ...string) model.Auth {
	header := req.Header.Get("Authorization")
	tokens, err := s.Token().GetAll()
	if err != nil {
		failAuth(res)
		return model.Auth{}
	}
	foundToken := false
	accountId := 0
	for _, token := range tokens {
		if header == "Bearer "+strings.TrimSpace(token.Token) {
			foundToken = true
			accountId = token.AccountId
		}
	}

	if !foundToken {
		failAuth(res)
		return model.Auth{}
	}

	return model.Auth{AccountId: accountId}
}

func failAuth(res http.ResponseWriter) {
	res.WriteHeader(http.StatusUnauthorized)
	fmt.Fprint(res, "Error: Could not authenticate\n")
}
