package authController

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"qaa/mail"
	"qaa/middlewares/authMiddleware"
	"qaa/services/usersService"
	"qaa/utils/authUtils"
	"qaa/utils/errorUtils"

	"golang.org/x/crypto/bcrypt"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	errorUtils.MethodNotAllowed_error(w, r)
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		log.Println("Error creating user: Email and password are required")
		http.Redirect(w, r, "/error?message=Email+and+password+are+required", http.StatusSeeOther)
		return
	}

	_, err := usersService.CreateUser(email, password)
	if err != nil {
		log.Println("Error creating user:", err)
		http.Redirect(w, r, "/error?message=Error+creating+user", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/?login=true", http.StatusSeeOther)
}

func Login(w http.ResponseWriter, r *http.Request) {
	errorUtils.MethodNotAllowed_error(w, r)

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		log.Println("Email and password are required")
		http.Redirect(w, r, "/error?message=Email+and+password+are+required", http.StatusSeeOther)
		return
	}

	user, err := usersService.GetUserByEmail(email)
	if err != nil {
		log.Println("Invalid email or password")
		http.Redirect(w, r, "/error?message=Invalid+email+or+password", http.StatusSeeOther)
		return
	}

	if !authUtils.CheckPassword(user, password) {
		log.Println("Invalid email or password")
		http.Redirect(w, r, "/error?message=Invalid+email+or+password", http.StatusSeeOther)
		return
	}

	// Generate token
	token, err := authUtils.GenerateToken(email)
	if err != nil {
		log.Println("Error generating token")
		http.Redirect(w, r, "/error?message=Error+generating+token", http.StatusSeeOther)
		return
	}

	// Set token in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,                    // Prevent JavaScript access
		Secure:   true,                    // Use only on HTTPS
		SameSite: http.SameSiteStrictMode, // prevent CSRF
		MaxAge:   3600 * 24,               // Token expires in 1 day
	})

	// Redirect to user dashboard or home page after successful login
	http.Redirect(w, r, "/questions", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the authentication token by setting an expired cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Setting MaxAge to -1 deletes the cookie
	})

	// Optionally, you can also invalidate the session or token on the server-side

	// Redirect to the homepage or login page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	errorUtils.MethodNotAllowed_error(w, r)
	email := r.FormValue("email")

	if email == "" {
		log.Println("Email is required")
		http.Redirect(w, r, "/error?message=Email+is+required", http.StatusSeeOther)
		return
	}

	user, err := usersService.GetUserByEmail(email)
	if err != nil || user == nil {
		log.Println("User does not exist with that email address")
		http.Redirect(w, r, "/error?message=User+does+not+exist+with+that+email+address", http.StatusSeeOther)
		return
	}

	// Generate a random token
	tokenBytes := make([]byte, 32)
	_, err = rand.Read(tokenBytes)
	if err != nil {
		log.Println("Error generating token")
		http.Redirect(w, r, "/error?message=Error+generating+token", http.StatusSeeOther)
		return
	}
	token := base64.URLEncoding.EncodeToString(tokenBytes)

	// Set the expiration time (e.g., 24 hours from now)
	expiration := time.Now().Add(24 * time.Hour)

	// Store the token with email and expiration time
	authUtils.ResetTokens[token] = authUtils.ResetTokenData{
		Email:      email,
		Expiration: expiration,
	}

	// Send the reset link via email
	resetLink := fmt.Sprintf(os.Getenv("URL")+"?token=%s", token)
	go mail.SendEmail(email, "Trading Alerts: Password Reset", fmt.Sprintf(
		"Click the link below to reset your password: %s", resetLink,
	))

	http.Redirect(w, r, "/reset-password-sent", http.StatusSeeOther)

	w.Write([]byte("Password reset email sent."))
}

func SetPassword(w http.ResponseWriter, r *http.Request) {
	errorUtils.MethodNotAllowed_error(w, r)
	token := r.FormValue("token")
	password := r.FormValue("password")

	tokenData, exists := authUtils.ResetTokens[token]
	if !exists {
		http.Redirect(w, r, "/token-expired", http.StatusSeeOther)
		return
	}

	// Check if the token has expired
	if time.Now().After(tokenData.Expiration) {
		log.Print("token has expired")
		delete(authUtils.ResetTokens, token)
		http.Redirect(w, r, "/token-expired", http.StatusSeeOther)
		return
	}
	log.Print("token is valid", tokenData.Expiration, time.Now())

	user, err := usersService.GetUserByEmail(tokenData.Email)
	if err != nil || user == nil {
		log.Println("User does not exist with that email address")
		http.Redirect(w, r, "/error?message=User+does+not+exist+with+that+email+address", http.StatusSeeOther)
		return
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password")
		http.Redirect(w, r, "/error?message=Error+hashing+password", http.StatusSeeOther)
		return
	}

	// Save the new password in your database (pseudo-code)
	err = usersService.UpdatePassword(tokenData.Email, string(hashedPassword))
	if err != nil {
		log.Println("Error saving new password")
		http.Redirect(w, r, "/error?message=Error+saving+new+password", http.StatusSeeOther)
		return
	}

	// Invalidate the token
	delete(authUtils.ResetTokens, token)
	http.Redirect(w, r, "/reset-password-success", http.StatusSeeOther)

	w.Write([]byte("Password has been reset successfully."))
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	errorUtils.MethodNotAllowed_error(w, r)
	email := r.Context().Value(authMiddleware.UserEmailKey).(string)

	err := usersService.DeleteAccount(email)

	if err != nil {
		http.Redirect(w, r, "/error?message=Error+deleting+account", http.StatusSeeOther)
		return
	}

	Logout(w, r)

	w.Write([]byte("Account successfully deleted"))
}

