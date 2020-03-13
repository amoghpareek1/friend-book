package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/dgrijalva/jwt-go"
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

	if user.Email == "" || user.Password == "" || user.FirstName == "" || user.LastName == "" {
		sendResponse(w, false, "All fields marked (*) are required.")
		return
	}

	user.Email = strings.ToLower(user.Email)

	if err := checkmail.ValidateFormat(user.Email); err != nil {
		log.Println(err)
		sendResponse(w, false, "Email address is not valid.")
		return
	}

	var existingUser User
	gormDatabase.Model(&User{}).Where(&User{
		Email: user.Email,
	}).Find(&existingUser)
	if existingUser.ID != "" {
		sendResponse(w, false, "This email address is already registered with us.")
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

	expirationTime := time.Now().Add(10 * time.Minute)

	claims := &Claims{
		Username: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    accessToken,
		Value:   tokenString,
		Expires: expirationTime,
	})

	session.Options.HttpOnly = true

	session.Save(r, w)

	sendResponse(w, true, "Sign in successful.")
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			sendResponse(w, false, "Unauthorised Request.")
			return
		}
		sendResponse(w, false, "Bad Request.")
		return
	}

	tknStr := c.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			sendResponse(w, false, "Unauthorised Request.")
			return
		}
		sendResponse(w, false, "Bad Request.")
		return
	}
	if !tkn.Valid {
		sendResponse(w, false, "Unauthorised Request.")
		return
	}

	sendResponse(w, true, "Request is valid.")
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
