package controller

import (
	// base core imports
	"fmt"
	"log"
	"net/http"

	// local imports
	"github.com/voidhofer/bingo-site/app/model"
	"github.com/voidhofer/bingo-site/app/shared/passhash"
	"github.com/voidhofer/bingo-site/app/shared/recaptcha"
	"github.com/voidhofer/bingo-site/app/shared/session"
	"github.com/voidhofer/bingo-site/app/shared/view"

	// external imports
	"github.com/gorilla/sessions"
	"github.com/josephspurrier/csrfbanana"
)

const (
	// Name of the session variable that tracks login attempts
	sessLoginAttempt = "login_attempt"
)

// loginAttempt increments the number of login attempts in sessions variable
func loginAttempt(sess *sessions.Session) {
	// Log the attempt
	if sess.Values[sessLoginAttempt] == nil {
		sess.Values[sessLoginAttempt] = 1
	} else {
		sess.Values[sessLoginAttempt] = sess.Values[sessLoginAttempt].(int) + 1
	}
}

// LoginGET displays the login page
func LoginGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Display the view
	v := view.New(r)
	v.Name = "login/login"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	// Refill any form fields
	view.Repopulate([]string{"email"}, r.Form, v.Vars)
	v.Render(w)
}

// LoginPOST handles the login form submission
func LoginPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Prevent brute force login attempts by not hitting MySQL and pretending like it was invalid :-)
	if sess.Values[sessLoginAttempt] != nil && sess.Values[sessLoginAttempt].(int) >= 50 {
		log.Println("Brute force login prevented")
		sess.AddFlash(view.Flash{"Sorry, no brute force :-)", view.FlashNotice})
		sess.Save(r, w)
		LoginGET(w, r)
		return
	}

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"email", "password"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		LoginGET(w, r)
		return
	}

	// Form values
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Get database result
	result, err := model.ClientByEmail(email)
	// Determine if client exists
	if err == model.ErrNoResult {
		loginAttempt(sess)
		sess.AddFlash(view.Flash{"Password is incorrect - Attempt: " + fmt.Sprintf("%v", sess.Values[sessLoginAttempt]), view.FlashWarning})
		sess.Save(r, w)
	} else if err != nil {
		// Display error message
		log.Println(err)
		sess.AddFlash(view.Flash{"There was an error. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else if passhash.MatchString(result.Password, password) {
		if result.Status != 1 {
			// Client inactive and display inactive message
			sess.AddFlash(view.Flash{"Account is inactive so login is disabled.", view.FlashNotice})
			sess.Save(r, w)
		} else {
			// Login successfully
			session.Empty(sess)
			sess.AddFlash(view.Flash{"Login successful!", view.FlashSuccess})
			sess.Values["id"] = result.ClientID()
			sess.Values["email"] = email
			sess.Values["first_name"] = result.FirstName
			sess.Save(r, w)
			http.Redirect(w, r, "/environment/data", http.StatusFound)
			return
		}
	} else {
		loginAttempt(sess)
		sess.Save(r, w)
	}

	// Show the login page again
	LoginGET(w, r)
}

// LogoutGET clears the session and logs the client out
func LogoutGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// If client is authenticated
	if sess.Values["id"] != nil {
		session.Empty(sess)
		sess.AddFlash(view.Flash{"Goodbye!", view.FlashNotice})
		sess.Save(r, w)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// RegisterPOST handles the registration form submission
func RegisterPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	// Prevent brute force login attempts by not hitting MySQL and pretending like it was invalid :-)
	if sess.Values["register_attempt"] != nil && sess.Values["register_attempt"].(int) >= 5 {
		log.Println("Brute force register prevented")
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}
	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"first_name", "last_name", "email", "password"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		LoginGET(w, r)
		return
	}
	// Validate with Google reCAPTCHA
	if !recaptcha.Verified(r) {
		sess.AddFlash(view.Flash{"reCAPTCHA invalid!", view.FlashError})
		sess.Save(r, w)
		LoginGET(w, r)
		return
	}
	// Get form values
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	password, errp := passhash.HashString(r.FormValue("password"))

	// If password hashing failed
	if errp != nil {
		log.Println(errp)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	fmt.Println("REGISTERING.....")
	// Get database result
	_, err := model.ClientByEmail(email)
	fmt.Println("REGISTERED")
	if err == model.ErrNoResult || err == model.ErrUnavailable { // If success (no user exists with that email)
		ex := model.ClientCreate(firstName, lastName, email, password)
		// Will only error if there is a problem with the query
		if ex != nil {
			log.Println(ex)
			sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
			sess.Save(r, w)
		} else {
			sess.AddFlash(view.Flash{"Account created successfully for: " + email, view.FlashSuccess})
			sess.Save(r, w)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
	} else if err != nil { // Catch all other errors
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else { // Else the user already exists
		sess.AddFlash(view.Flash{"Account already exists for: " + email, view.FlashError})
		sess.Save(r, w)
	}
	// Display the page
	LoginGET(w, r)
}
