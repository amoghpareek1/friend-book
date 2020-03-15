package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/gofrs/uuid"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println(err)
		sendResponse(w, false, errServer)
		return
	}
	r.Body.Close()

	if user.Email == "" || user.Password == "" || user.Name == "" || user.Phone == "" {
		sendResponse(w, false, "All fields marked (*) are required.")
		return
	}

	user.Email = strings.ToLower(user.Email)

	if err := checkmail.ValidateFormat(user.Email); err != nil {
		log.Println(err)
		sendResponse(w, false, "Email address is not valid.")
		return
	}

	if len(user.Phone) != 10 {
		sendResponse(w, false, "Phone number is not valid.")
		return
	}

	var existingUser User
	gormDatabase.Model(&User{}).Where("email = ? OR phone = ?", user.Email, user.Phone).Find(&existingUser)
	if existingUser.ID != "" {
		sendResponse(w, false, "User with this email address or phone number is already registered with us.")
		return
	}

	b, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		sendResponse(w, false, errServer)
		return
	}
	user.Password = string(b)

	uuid, err := uuid.NewV4()
	if err != nil {
		sendResponse(w, false, errServer)
		return
	}

	user.ID = uuid.String()

	gormDatabase.Model(&User{}).Create(&user)

	sendResponse(w, true, "Sign up successful. Please login with your credentials.")
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		sendResponse(w, false, errServer)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendResponse(w, false, errServer)
		return
	}
	r.Body.Close()

	if user.Email == "" || user.Password == "" {
		sendResponse(w, false, "All fields are required.")
		return
	}

	var userX User
	gormDatabase.Where(&User{
		Email: user.Email,
	}).First(&userX)

	if userX.ID == "" {
		sendResponse(w, false, "This email address is not registered.")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userX.Password), []byte(user.Password)); err != nil {
		if err != bcrypt.ErrMismatchedHashAndPassword {
			log.Println(err)
		}
		sendResponse(w, false, "Password is not valid.")
		return
	}

	session.Values["userID"] = userX.ID

	session.Options.HttpOnly = true

	session.Save(r, w)

	sendResponse(w, true, "Sign in successful.")
}

func signOutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		return
	}

	session.Options = &sessions.Options{
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	}

	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
