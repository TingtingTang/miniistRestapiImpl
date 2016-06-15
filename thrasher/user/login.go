package main

import (
	"./mgodb"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	valid "github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	// "github.com/thinkhy/logops"
	"log"
	"net/http"
	"time"
)

var AuthorizationHeader string = "Authorization"

type LoginRequestData struct {
	Action   string `json:action`
	Name     string `json:name,omitempty`
	Email    string `json:email,omitempty`
	Password string `json:password`
}

type LoginResponseData struct {
	*mgodb.Person
	Token string `json:token`
}

type LogoutRequestData struct {
	Action string `json:action`
	Name   string `json:user`
	// JWT token is stored in http header
}

type LogoutResponseData struct {
	mgodb.Person
	Token string `json:token`
}

func handleLogin(w http.ResponseWriter, r *http.Request) *apiError {
	log.Printf("[user.handleLogin] handleLogin Entry %s\n", r.URL)

	// Read request JSON data
	var d LoginRequestData
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&d)
	if err != nil {
		str := fmt.Sprintf("Failed to parse JSON %v", err)
		log.Println("[user.handleLogin] ", str)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -1, Info: str}
	}

	// Validate input data
	if d.Action != "login" {
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: `action should be "login"`}
	} else if len(d.Name) > 0 && valid.Matches(d.Name, `(?i)^[a-z\d.]{5,}$`) == false {
		w.WriteHeader(http.StatusBadRequest)
		info := `user name should match the pattern "/^[a-z\d.]{5,}$/i"`
		return &apiError{Error: err, Ret: -2, Info: info}
	} else if len(d.Email) > 0 && valid.IsEmail(d.Email) == false {
		w.WriteHeader(http.StatusBadRequest)
		info := `email format is invalid`
		return &apiError{Error: err, Ret: -2, Info: info}
	} else if len(d.Name) == 0 && len(d.Email) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		info := `both name and email are empty`
		return &apiError{Error: err, Ret: -2, Info: info}
	}

	var rd LoginResponseData
	// Post-processing
	// Write back user data except for password
	if len(d.Name) > 0 {
		rd.Person, err = mgodb.GetUserByName(d.Name)
	} else {
		rd.Person, err = mgodb.GetUserByEmail(d.Email)
	}
	if err != nil {
		log.Println("[user handleLogin] Failed to get user data: ", d.Name, d.Email)
		w.WriteHeader(http.StatusInternalServerError)
		return &apiError{Error: err, Ret: -4, Info: "Failed to get user data"}
	}

	// Authenticate user
	hash := sha1.Sum([]byte(d.Password))
	if mgodb.IsUserOrEmailExisted(d.Name, d.Email, hash[:]) == false {
		str := fmt.Sprintf("Failed to login with user %s, email: %s", d.Name, d.Email)
		log.Println("[user.handleLogin] ", str)
		w.WriteHeader(http.StatusUnauthorized)
		return &apiError{Error: err, Ret: -3, Info: str}
	}

	// Create JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	uuid := Uuidv4()
	token.Claims["id"] = uuid
	token.Claims["iat"] = time.Now().Unix()
	token.Claims["exp"] = time.Now().Add(Config.Expiration).Unix()
	token.Claims["user"] = rd.Person.Name
	jwtString, err := token.SignedString([]byte(Config.Secret))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return &apiError{Error: err, Ret: -3, Info: "Failed to create JWT token"}
	}

	// [2016-05-10] added TTL index for session collection
	if err = mgodb.InsertSession(uuid); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return &apiError{Error: err, Ret: -3, Info: "Failed to insert session ID"}
	}
	rd.Token = jwtString
	rd.Ret = 0
	rd.Info = "OK"
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(rd); err != nil {
		log.Println("[user.handleLogin] Failed to encode JSON data")
		w.WriteHeader(http.StatusInternalServerError)
		return &apiError{Error: err, Ret: -5, Info: err.Error()}
	}
	opshdlr.Write("User", rd.Person.Name, "login", "", "")
	log.Printf("[user.handleLogin] %s login successfully\n", rd.Person.Name)

	return nil
}

func handleLogout(w http.ResponseWriter, r *http.Request) *apiError {
	log.Printf("[user.handleLogout] handleLogout Entry %s\n", r.URL)

	// Read request JSON data
	var d LogoutRequestData
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&d)
	if err != nil {
		str := fmt.Sprintf("Failed to parse JSON %v", err)
		log.Println("[user.handleLogout] ", str)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -1, Info: str}
	}

	// Validate input data
	if d.Action != "logout" {
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: `action should be "logout"`}
	} else if valid.Matches(d.Name, `(?i)^[a-z\d.]{5,}$`) == false {
		w.WriteHeader(http.StatusBadRequest)
		info := `user name should match the pattern "/^[a-z\d.]{5,}$/i"`
		return &apiError{Error: err, Ret: -2, Info: info}
	}

	// Authenticate user
	jwtString := r.Header.Get("Authorization")
	key := []byte(Config.Secret)
	log.Println("[user.handleLogout] Verifying JWT token ", jwtString)
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what we expect
		// Comparing agaist *jwt.SigningMethodHSA if using an HS256 token
		// Refer to: https://github.com/dgrijalva/jwt-go/issues/123
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return key, errors.New("Unexpected signing method")
		}
		if mgodb.IsSessionIdExisted(token.Claims["id"].(string)) == false {
			return key, errors.New("session ID is invalid")
		}
		return key, nil
	})

	if err != nil || token.Valid == false {
		str := fmt.Sprintf("Failed to authenticate user %s: %v", d.Name, err)
		log.Println("[user.handleLogin] ", str)
		w.WriteHeader(http.StatusUnauthorized)
		return &apiError{Error: err, Ret: -3, Info: str}
	}

	uuid := token.Claims["id"].(string)
	if err = mgodb.DeleteSession(uuid); err != nil {
		str := fmt.Sprintf("Failed to delete session id %s", uuid)
		log.Println("[user.handleLogin] ", str)
		w.WriteHeader(http.StatusInternalServerError)
		return &apiError{Error: err, Ret: -3, Info: str}
	}

	// Post-processing
	var rd LogoutResponseData
	rd.Token = ""
	rd.Ret = 0
	rd.Info = "OK"
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(rd); err != nil {
		log.Println("[user.handleLogout] Failed to encode JSON data")
		w.WriteHeader(http.StatusInternalServerError)
		return &apiError{Error: err, Ret: -5, Info: err.Error()}
	}
	opshdlr.Write("User", d.Name, "logout", "", "")
	log.Printf("[user.handleLogout] %s logout successfully\n", d.Name)

	return nil
}
