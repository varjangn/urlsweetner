package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/varjangn/urlsweetner/db"
	"github.com/varjangn/urlsweetner/models"
	"golang.org/x/crypto/bcrypt"
)

type registerForm struct {
	email     string
	password  string
	firstname string
	lastname  string
}

type loginForm struct {
	email    string
	password string
}

type loginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	cType := r.Header.Get("Content-Type")
	if cType != "application/x-www-form-urlencoded" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	formData := registerForm{
		email:     r.FormValue("email"),
		password:  r.FormValue("password"),
		firstname: r.FormValue("firstname"),
		lastname:  r.FormValue("lastname"),
	}
	if existingUser, _ := db.DbRepo.GetUser(formData.email); existingUser != nil {
		http.Error(w, "User Already Exists!", http.StatusBadRequest)
		return
	}
	user, err := models.NewUser(formData.email, formData.password, formData.firstname, formData.lastname)
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	err = db.DbRepo.AddUserToDB(user)
	if err != nil {
		http.Error(w, "Failed to Add User", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func Login(w http.ResponseWriter, r *http.Request) {
	cType := r.Header.Get("Content-Type")
	if cType != "application/x-www-form-urlencoded" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	formData := loginForm{
		email:    r.FormValue("email"),
		password: r.FormValue("password"),
	}
	if formData.email == "" || formData.password == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	user, err := db.DbRepo.GetUser(formData.email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(formData.password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Minute * 10).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	resp := loginResponse{
		Token: tokenString,
		User:  *user,
	}
	cookie := &http.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Expires:  time.Now().Add(20 * time.Minute),
		HttpOnly: false,
	}
	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
