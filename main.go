package main

import (
	"flag"
	"net/http"

	"log"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler).Methods("GET")

	router.HandleFunc("/api/v1/sign-up", signUpHandler).Methods("POST")
	router.HandleFunc("/api/v1/sign-in", signInHandler).Methods("POST")
	router.HandleFunc("/api/v1/sign-out", signOutHandler).Methods("GET")
	router.HandleFunc("/api/v1/validate", validateHandler).Methods("GET")

	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public")))).Methods("GET")

	router.NotFoundHandler = http.HandlerFunc(indexHandler)

	server := &http.Server{
		Addr:    ":" + flag.Lookup("port").Value.String(),
		Handler: cors.Default().Handler(gziphandler.GzipHandler(noCacheMW(router))),
	}

	flag.VisitAll(func(flag *flag.Flag) {
		log.Println(flag.Name, "->", flag.Value)
	})

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func noCacheMW(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")

		h.ServeHTTP(w, r)
	})
}
