package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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

	var friendships []Friendship
	gormDatabase.Model(&Friendship{}).
		Select("to_user, from_user").
		Where("to_user = ? OR from_user = ?", user.ID, user.ID).
		Find(&friendships)

	mapOfUniqueUserIDs := make(map[string]string, 0)
	sliceOfUserIDs := make([]string, 0)

	for _, v := range friendships {
		if _, ok := mapOfUniqueUserIDs[v.FromUser]; !ok {
			mapOfUniqueUserIDs[v.FromUser] = v.FromUser
			sliceOfUserIDs = append(sliceOfUserIDs, v.FromUser)
		}

		if _, ok1 := mapOfUniqueUserIDs[v.ToUser]; !ok1 {
			mapOfUniqueUserIDs[v.ToUser] = v.ToUser
			sliceOfUserIDs = append(sliceOfUserIDs, v.ToUser)
		}
	}

	gormDatabase.Model(&User{}).
		Where("id != ? AND id IN (?)", user.ID, sliceOfUserIDs).
		Find(&user.Friends)

	sendResponse(w, true, user)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	type ResponseData struct {
		User
		Status string
	}

	vars := mux.Vars(r)

	userID := vars["userID"]

	if userID == "" {
		sendResponse(w, false, "Invalid userID recieved.")
		return
	}

	session, err := sessionStore.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		return
	}

	if session.Values["userID"] == nil {
		sendResponse(w, false, "Invalid request.")
		return
	}

	userID2 := session.Values["userID"].(string)

	var user User
	gormDatabase.Model(&User{}).Where(&User{
		ID: userID,
	}).Find(&user)

	if user.ID == "" {
		sendResponse(w, false, "No such user exists.")
		return
	}

	var responseData ResponseData

	responseData.User = user

	var friendship Friendship
	gormDatabase.Model(&Friendship{}).
		Where("(from_user = ? AND to_user = ?) OR (to_user = ? AND from_user = ?)", userID2, userID, userID, userID2).
		First(&friendship)

	if friendship.ID != "" {
		responseData.Status = friendship.Status
	}

	sendResponse(w, true, responseData)
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	type ResponseData struct {
		Users      []User
		TotalCount int
	}

	session, err := sessionStore.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		return
	}

	if session.Values["userID"] == nil {
		sendResponse(w, false, "Invalid request.")
		return
	}

	userID := session.Values["userID"].(string)
	users := make([]User, 0)
	values := r.URL.Query()

	sortBy := values.Get("SortBy")
	orderBy := values.Get("OrderBy")
	searchBy := values.Get("searchBy")
	limit, _ := strconv.Atoi(values.Get("limit"))
	offset, _ := strconv.Atoi(values.Get("offset"))

	if offset < 0 || limit == 0 {
		sendResponse(w, false, "Invalid request.")
		return
	}

	orderByClauseValue := ""

	if orderBy == "name" {
		orderByClauseValue = "LOWER(" + orderBy + ") " + sortBy
	} else {
		orderByClauseValue = orderBy + " " + sortBy
	}

	if searchBy != "" {
		gormDatabase.Model(&User{}).
			Order(orderByClauseValue).
			Where("id != ? AND (LOWER(name) LIKE ? OR LOWER(email) LIKE ?)", userID, searchBy+"%", searchBy+"%").
			Limit(limit).
			Offset(offset).
			Find(&users)
	} else {
		gormDatabase.Model(&User{}).
			Order(orderByClauseValue).
			Where("id != ?", userID).
			Limit(limit).
			Offset(offset).
			Find(&users)
	}

	var responseData ResponseData
	responseData.Users = users

	var totalCount int
	gormDatabase.Model(&User{}).Where("id != ?", userID).Count(&totalCount)
	responseData.TotalCount = totalCount

	sendResponse(w, true, responseData)
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

	existingUser.Name = user.Name
	existingUser.Phone = user.Phone
	existingUser.Email = user.Email

	gormDatabase.Model(&User{}).Save(&existingUser)

	sendResponse(w, true, "Details Updated.")
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		return
	}

	if session.Values["userID"] == nil {
		sendResponse(w, false, "Invalid request.")
		return
	}

	userID := session.Values["userID"].(string)

	var existingUser User
	gormDatabase.Where(&User{
		ID: userID,
	}).First(&existingUser)
	if existingUser.ID == "" {
		sendResponse(w, false, "No such user exists.")
		return
	}

	var friendships []Friendship
	gormDatabase.Model(&Friendship{}).
		Where("to_user = ? OR from_user = ?", existingUser.ID, existingUser.ID).
		Delete(&friendships)

	gormDatabase.Model(&User{}).Delete(&existingUser)

	sendResponse(w, true, "Profile deleted.")
}
