package main

import (
	"flag"
	"net/http"

	"log"

	"github.com/NYTimes/gziphandler"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler).Methods("GET")

	router.HandleFunc("/api/v1/sign-up", signUpHandler).Methods("POST")
	router.HandleFunc("/api/v1/sign-in", signInHandler).Methods("POST")
	router.HandleFunc("/api/v1/sign-out", signOutHandler).Methods("GET")
	router.HandleFunc("/api/v1/get-me", getMeHandler).Methods("GET")

	router.HandleFunc("/api/v1/users/get", getUsersHandler).Methods("GET")

	router.HandleFunc("/api/v1/user/update", updateUserHandler).Methods("PUT")
	router.HandleFunc("/api/v1/user/delete", deleteUserHandler).Methods("DELETE")
	router.HandleFunc("/api/v1/user/{userID}/get", getUserHandler).Methods("GET")

	router.HandleFunc("/api/v1/send-friend-request/{userID}", sendFriendRequestHandler).Methods("POST")
	router.HandleFunc("/api/v1/send-unfriend-request/{userID}", sendUnfriendRequestHandler).Methods("POST")

	router.HandleFunc("/api/v1/recieve-friend-request/{friendshipID}", recieveFriendRequestHandler).Methods("GET")
	router.HandleFunc("/api/v1/recieve-unfriend-request/{friendshipID}", recieveUnfriendRequestHandler).Methods("GET")

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
