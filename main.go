package main

import (
	"log"
	"net/http"
	"os"

	ghandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// api headers
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

	// start server
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
