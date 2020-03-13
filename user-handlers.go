package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func getMeHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		return
	}

	if session.Values["userID"] == nil {
		sendResponse(w, false, "Invalid request.")
		return
	}

	var user User
	gormDatabase.Model(&User{}).Where(&User{
		ID: session.Values["userID"].(string),
	}).First(&user)

	if user.ID == "" {
		sendResponse(w, false, "No such user exists.")
		return
	}

	q := "SELECT * FROM users WHERE id != ? AND id IN (SELECT to_user FROM friendships WHERE to_user = ? OR from_user = ?)"

	gormDatabase.Raw(q, user.ID, user.ID, user.ID).Scan(&user.Friends)

	log.Println("users -----", user.Friends)

	sendResponse(w, true, user)
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		return
	}

	if session.Values["userID"] == nil {
		sendResponse(w, false, "Invalid request.")
		return
	}

	users := make([]User, 0)

	values := r.URL.Query()

	sortBy := values.Get("SortBy")
	orderBy := values.Get("OrderBy")

	if orderBy == "name" {
		gormDatabase.Model(&User{}).Order("LOWER("+orderBy+") "+sortBy).Where("id != ?", session.Values["userID"].(string)).Find(&users)
	} else {
		gormDatabase.Model(&User{}).Order(orderBy+" "+sortBy).Where("id != ?", session.Values["userID"].(string)).Find(&users)
	}

	// userID := session.Values["userID"].(string)
	// q := "SELECT * FROM users WHERE id != ? AND id IN (SELECT to_user FROM friendships WHERE to_user = ? OR from_user = ?)"

	// var userss []User
	// gormDatabase.Raw(q, userID, userID, userID).Scan(&userss)
	// userID := session.Values["userID"].(string)
	// q := "SELECT * FROM users WHERE id != ? AND id IN (SELECT to_user FROM friendships WHERE to_user = ? OR from_user = ?)"

	// var userss []User
	// gormDatabase.Raw(q, userID, userID, userID).Scan(&userss)

	// log.Println(userss)
	// log.Println(userss)

	sendResponse(w, true, users)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendResponse(w, false, errServer)
		return
	}
	r.Body.Close()

	if user.Name == "" || user.Phone == "" || user.Email == "" {
		sendResponse(w, false, "All fields are required.")
		return
	}

	var existingUser User

	gormDatabase.Model(&User{}).Where(&User{
		ID: user.ID,
	}).Find(&existingUser)

	if existingUser.ID == "" || existingUser.ID != user.ID {
		sendResponse(w, false, "No such user exists.")
		return
	}

	var userX User
	gormDatabase.Model(&User{}).Where("email = ? OR phone = ?", user.Email, user.Phone).Find(&userX)

	if userX.ID != "" && userX.ID != user.ID {
		sendResponse(w, false, "User with this email or phone is already registered with us.")
		return
	}

	gormDatabase.Model(&User{}).Save(&user)

	sendResponse(w, true, "Details Updated.")
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["userID"] == "" {
		sendResponse(w, false, "Invalid Request.")
		return
	}

	var existingUser User
	gormDatabase.Where(&User{
		ID: vars["userID"],
	}).First(&existingUser)
	if existingUser.ID == "" {
		sendResponse(w, false, "No such user exists.")
		return
	}

	var friendships []Friendship
	gormDatabase.Model(&Friendship{}).
		Where("toUser = ? OR fromUser = ?", existingUser.ID, existingUser.ID).
		Delete(&friendships)

	gormDatabase.Model(&User{}).Delete(&existingUser)
}
