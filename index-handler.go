package main

import (
	"log"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Cache-Control", "no-store")

	if _, ok := session.Values["userID"]; ok {
		http.ServeFile(w, r, "./templates/index-x.html")
		return
	}

	http.ServeFile(w, r, "./templates/index.html")
}
