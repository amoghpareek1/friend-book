package main

import (
	"log"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

func sendFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		sendResponse(w, false, errServer)
		return
	}

	if session.Values["userID"] == nil {
		sendResponse(w, false, "Invalid request.")
		return
	}

	fromUser := session.Values["userID"].(string)

	vars := mux.Vars(r)
	toUser := vars["userID"]

	// validating both users
	if toUser == "" || fromUser == "" {
		sendResponse(w, false, "invalid request.")
		return
	}

	var friendship Friendship

	friendship.ToUser = toUser
	friendship.FromUser = fromUser

	var user1 User
	gormDatabase.Model(&User{}).Where(&User{
		ID: friendship.FromUser,
	}).First(&user1)

	if user1.ID == "" {
		sendResponse(w, false, "invalid request.")
		return
	}

	var user2 User
	gormDatabase.Model(&User{}).Where(&User{
		ID: friendship.ToUser,
	}).First(&user2)

	if user2.ID == "" {
		sendResponse(w, false, "The user to which to you want to send request is not registered with us.")
		return
	}

	friendship.Status = requestInProcess

	uuid, err := uuid.NewV4()
	if err != nil {
		sendResponse(w, false, errServer)
		return
	}

	friendship.ID = uuid.String()

	gormDatabase.Model(&Friendship{}).Create(&friendship)

	subject := "Friend Request Sent"
	message := "<p>Hi " + user2.Name + "</p>" +
		"<p>Please click on the following link to accept friend request from " + user1.Name + ".</p>" +
		"<p><a href='" + Config().GetString("serverURL") + "/api/v1/recieve-friend-request/" + friendship.ID + "'>Accept Request</a></p>" +
		"<p>Thank you.</p>"

	sendEmailFrom := Config().GetString("sesEmail")

	// send email
	if _, err := sesConfig.SendEmailHTML(sendEmailFrom, user2.Email, subject, message, message); err != nil {
		log.Println(err)
		sendResponse(w, false, "Problem sending request email.")
		return
	}

	friendship.Status = friendRequestSent

	gormDatabase.Model(&Friendship{}).Save(&friendship)
}

func sendUnfriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		sendResponse(w, false, errServer)
		return
	}

	if session.Values["userID"] == nil {
		sendResponse(w, false, "Invalid request.")
		return
	}

	fromUser := session.Values["userID"].(string)

	vars := mux.Vars(r)
	toUser := vars["userID"]

	// validating both users
	if toUser == "" || fromUser == "" {
		sendResponse(w, false, "invalid request.")
		return
	}

	var friendship Friendship

	gormDatabase.Model(&Friendship{}).
		Where("(to_user = ? AND from_user = ?) OR (from_user = ? AND to_user = ?)", toUser, fromUser, fromUser, toUser).
		First(&friendship)

	if friendship.ID == "" {
		sendResponse(w, false, "Invalid request.")
		return
	}

	var user1 User
	gormDatabase.Model(&User{}).Where(&User{
		ID: friendship.FromUser,
	}).First(&user1)

	if user1.ID == "" {
		sendResponse(w, false, "invalid request.")
		return
	}

	var user2 User
	gormDatabase.Model(&User{}).Where(&User{
		ID: friendship.ToUser,
	}).First(&user2)

	if user2.ID == "" {
		sendResponse(w, false, "The user to which to you want to send request is not registered with us.")
		return
	}

	subject := "Delete Friend Request"
	message := "<p>Hi " + user2.Name + "</p>" +
		"<p>Please click on the following link to remove friend request from " + user1.Name + ".</p>" +
		"<p><a href='" + Config().GetString("serviceURL") + "/api/v1/recieve-unfriend-request/" + friendship.ID + "'>Remove Friend</a></p>" +
		"<p>Thank you.</p>"

	sendEmailFrom := Config().GetString("sesEmail")

	if _, err := sesConfig.SendEmailHTML(sendEmailFrom, user2.Email, subject, message, message); err != nil {
		log.Println(err)
		sendResponse(w, false, "Problem sending request email.")
		return
	}

	friendship.Status = friendRequestSent

	gormDatabase.Model(&Friendship{}).Save(&friendship)
}

func recieveFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		sendResponse(w, false, errServer)
		return
	}

	vars := mux.Vars(r)
	friendshipID := vars["friendshipID"]

	if friendshipID == "" {
		sendResponse(w, false, "Invalid friendshipID recieved.")
		return
	}

	var friendship Friendship
	gormDatabase.Model(&Friendship{}).Where(&Friendship{
		ID: friendshipID,
	}).First(&friendship)

	if friendship.ID == "" {
		sendResponse(w, false, "Friendship request is not valid.")
		return
	}

	gormDatabase.Model(&Friendship{}).Select("status").Updates(map[string]interface{}{
		"status": friendRequestApproved,
	})

	session.Values["userID"] = friendship.ToUser

	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func recieveUnfriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		sendResponse(w, false, errServer)
		return
	}

	vars := mux.Vars(r)
	friendshipID := vars["friendshipID"]

	if friendshipID == "" {
		sendResponse(w, false, "Invalid friendshipID recieved.")
		return
	}

	var friendship Friendship
	gormDatabase.Model(&Friendship{}).Where(&Friendship{
		ID: friendshipID,
	}).Delete(&friendship)

	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
