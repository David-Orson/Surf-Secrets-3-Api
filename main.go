package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/David-Orson/Surf-Secrets-3-Api/config"
	"github.com/David-Orson/Surf-Secrets-3-Api/store"
	"github.com/David-Orson/Surf-Secrets-3-Api/store/psqlstore"

	ghandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	var s store.Store
	var err error

	s, err = psqlstore.Open(config.LoadFile("./db.json"))
	if err != nil {
		log.Fatal("0001: Can't start the server without a store.")
		log.Println(err)
		return
	}

	router := mux.NewRouter()

	originsOk := ghandlers.AllowedOrigins([]string{
		os.Getenv("ORIGIN_ALLOWED"),
		"*",
	})
	headersOk := ghandlers.AllowedHeaders([]string{
		"Access-Control-Allow-Origin",
		"X-Requested-With",
		"Content-Type",
		"Authorization",
		"range",
	})
	methodsOk := ghandlers.AllowedMethods([]string{
		"GET",
		"POST",
		"PUT",
		"DELETE",
		"HEAD",
		"OPTIONS",
	})

	// Don't worry golang, I was going to use it later...
	// GOLANG: NO FUCK YOU, VARIABLE NOT USED!
	fmt.Print(s)

	log.Println("Listening on :8085")
	log.Println(
		http.ListenAndServe(
			":8085",
			ghandlers.CORS(
				originsOk,
				headersOk,
				methodsOk,
			)(router),
		),
	)
}
