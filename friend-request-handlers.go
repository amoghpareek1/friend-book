package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

func sendFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	var friendship Friendship
	if err := json.NewDecoder(r.Body).Decode(&friendship); err != nil {
		sendResponse(w, false, errServer)
		return
	}
	r.Body.Close()

	if friendship.ToUser == "" || friendship.FromUser == "" {
		sendResponse(w, false, "invalid request.")
		return
	}

	var user User
	gormDatabase.Model(&User{}).Where(&User{
		ID: friendship.FromUser,
	}).First(&user)

	if user.ID == "" {
		sendResponse(w, false, "invalid request.")
		return
	}

	var toUser User
	gormDatabase.Model(&User{}).Where(&User{
		ID: friendship.ToUser,
	}).First(&toUser)

	if toUser.ID == "" {
		sendResponse(w, false, "The user to which to you want to send request is not registered with us.")
		return
	}

	friendship.Status = "initiated"

	uuid, err := uuid.NewV4()
	if err != nil {
		sendResponse(w, false, errServer)
		return
	}

	friendship.ID = uuid.String()

	gormDatabase.Model(&Friendship{}).Create(&friendship)

	subject := "Friend Request Sent"
	message := "<p>Hi " + user.Name + "</p>" +
		"<p>Please click on the following link to reset your password.</p>" +
		"<p><a href='" + Config().GetString("serviceURL") + "/api/v1/recieve-friend-request" + friendship.ID + "'>Accept Request</a></p>" +
		"<p>Thank you.</p>"

	sendEmailFrom := Config().GetString("SesEmail")

	if _, err := sesConfig.SendEmailHTML(sendEmailFrom, toUser.Email, subject, message, message); err != nil {
		log.Println(err)
		sendResponse(w, false, "Problem sending password assistance email.")
		return
	}

	friendship.Status = "request sent"

	gormDatabase.Model(&Friendship{}).Save(&friendship)
}

func recieveFriendRequest(w http.ResponseWriter, r *http.Request) {
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
		"status": requestApproved,
	})

	session.Values["userID"] = friendship.ToUser

	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
