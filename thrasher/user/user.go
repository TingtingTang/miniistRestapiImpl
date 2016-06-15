package main

import (
	"./mgodb"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	valid "github.com/asaskevich/govalidator"
	"log"
	"net/http"
)

type UserRequestData struct {
	Action      string `json:"action"`
	Id          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Password    string `json:"password,omitempty"`
	Email       string `json:"email,omitempty"`
	Team        string `json:"team,omitempty"`
	TsoUser     string `json:"tso_user,omitempty"`
	TsoPassword string `json:"tso_password,omitempty"`
	NewPassword string `json:"new_password,omitempty"`
}

type ResponseData struct {
	Ret  int    `bson:"-" json:"ret"`
	Info string `bson:"-" json:"info"`
}

type UserResponseData struct {
	*mgodb.Person
}

//    if user name is not coincide with JWT token, tso password should not be returned.
func handleGetUser(w http.ResponseWriter, r *http.Request) *apiError {
	log.Printf("[user.handleGetUser] handleGetUser Entry %s\n", r.URL)

	username := r.URL.Query()["name"]

	var d UserRequestData
	var err error

	if len(username) <= 0 {
		// Read request JSON data
		dec := json.NewDecoder(r.Body)
		err = dec.Decode(&d)
		if err != nil {
			str := fmt.Sprintf("Failed to parse JSON %v", err)
			log.Println("[user.handleGetUser] ", str)
			w.WriteHeader(http.StatusBadRequest)
			return &apiError{Error: err, Ret: -1, Info: str}
		}
	} else { // else read parameter from URL
		d.Action = "get_user"
		d.Name = username[0]
	}

	// Validate input data
	if d.Action != "get_user" {
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: `action should be "update_user"`}
	} else if valid.Matches(d.Name, `(?i)^[a-z\d.]{5,}$`) == false {
		w.WriteHeader(http.StatusBadRequest)
		info := `user name should match the pattern "/^[a-z\d.]{5,}$/i"`
		return &apiError{Error: err, Ret: -2, Info: info}
	} else if len(d.Email) != 0 && valid.IsEmail(d.Email) == false {
		info := `email format is invalid`
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: info}
	}

	ae, _ := withAuthentication(w, r, d.Name)
	if ae != nil {
		return ae
	}

	var rd UserResponseData
	// Post-processing
	// Write back user data except for password
	rd.Person, err = mgodb.GetUserByName(d.Name)
	if len(d.Name) > 0 {
		rd.Person, err = mgodb.GetUserByName(d.Name)
	} else {
		rd.Person, err = mgodb.GetUserByEmail(d.Email)
	}
	if err != nil {
		log.Println("[user handleGetUser] Failed to get user data: ", d.Name)
		w.WriteHeader(http.StatusInternalServerError)
		return &apiError{Error: err, Ret: -4, Info: err.Error()}
	} else {
		rd.Ret = 0
		rd.Info = "OK"
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(rd); err != nil {
			log.Println("[user.handleGeteUser] Failed to encode JSON data")
			w.WriteHeader(http.StatusInternalServerError)
			return &apiError{Error: err, Ret: -5, Info: err.Error()}
		}
		log.Printf("[user.handleGetUser] %s get user data successfully\n", d.Name)
	}

	return nil
}

// TODO: [2016-05-09]  verify user with that in JWT token
//    if user name is not coincide with JWT token, tso password should not be returned.
func handleUpdateUser(w http.ResponseWriter, r *http.Request) *apiError {
	log.Printf("[user.handleUpdateUser] handleUpdateUser Entry %s\n", r.URL)

	// Read request JSON data
	var d UserRequestData
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&d)
	if err != nil {
		str := fmt.Sprintf("Failed to parse JSON %v", err)
		log.Println("[user.handleUpdateUser] ", str)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -1, Info: str}
	}

	// Validate input data
	if d.Action != "update_user" {
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: `action should be "update_user"`}
	} else if valid.Matches(d.Name, `(?i)^[a-z\d.]{5,}$`) == false {
		w.WriteHeader(http.StatusBadRequest)
		info := `user name should match the pattern "/^[a-z\d.]{5,}$/i"`
		return &apiError{Error: err, Ret: -2, Info: info}
	} else if len(d.Email) != 0 && valid.IsEmail(d.Email) == false {
		info := `email format is invalid`
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: info}
	}

	ae, _ := withAuthentication(w, r, d.Name)
	if ae != nil {
		return ae
	}

	p := &mgodb.Person{
		Name:        d.Name,
		Email:       d.Email,
		Team:        d.Team,
		TsoUser:     d.TsoUser,
		TsoPassword: d.TsoPassword,
	}
	err = mgodb.UpdateUser(p)
	if err != nil {
		log.Printf("[user handleUpdateUser] Failed to update user data: %s %s\n", d.Name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return &apiError{Error: err, Ret: -3, Info: err.Error()}
	}

	var rd UserResponseData
	// Post-processing
	// Write back user data except for password
	rd.Person, err = mgodb.GetUserByName(d.Name)
	if err != nil {
		log.Println("[user handleUpdateUser] Failed to get user data: ", d.Name)
		w.WriteHeader(http.StatusInternalServerError)
		return &apiError{Error: err, Ret: -4, Info: err.Error()}
	} else {
		rd.Ret = 0
		rd.Info = "OK"
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(rd); err != nil {
			log.Println("[user.handleUpdateUser] Failed to encode JSON data")
			w.WriteHeader(http.StatusInternalServerError)
			return &apiError{Error: err, Ret: -5, Info: err.Error()}
		}
		opshdlr.Write("User", d.Name, "update", "user", "")
		log.Printf("[user.handleUpdateUser] %s update user data successfully\n", d.Name)
	}

	return nil
}

func handleDeleteUser(w http.ResponseWriter, r *http.Request) *apiError {
	log.Printf("[user.handleDeleteUser] handleDeleteUser Entry %s\n", r.URL)

	// Read request JSON data
	var d UserRequestData
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&d)
	if err != nil {
		str := fmt.Sprintf("failed to parse JSON %v", err)
		log.Println("[user.handleJoin] ", str)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -1, Info: str}
	}

	// Validate data format
	if d.Action != "delete_user" {
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: `action should be "delete_user"`}
	} else if valid.Matches(d.Name, `(?i)^[a-z\d.]{5,}$`) == false {
		info := `user name should match the pattern "/^[a-z\d.]{5,}$/i"`
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: info}
	} else if valid.Matches(d.Password, `^.{6,}$`) == false {
		info := `password are case sensitive and must contain a minimum of 6 English characters`
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: info}
	}

	ae, uuid := withAuthentication(w, r, d.Name)
	if ae != nil {
		return ae
	}

	hash := sha1.Sum([]byte(d.Password))
	err = mgodb.DeleteUser(d.Name, hash[:])
	if err != nil {
		str := fmt.Sprintf("Failed to delete user: %s", d.Name)
		log.Println("[user.handleDeleteUser]", str)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -3, Info: str}
	}

	// Write header before write JSON data
	w.WriteHeader(http.StatusOK)

	var rd ResponseData
	rd.Ret = 0
	rd.Info = "OK"
	if err = json.NewEncoder(w).Encode(rd); err != nil {
		log.Println("[user.handleDeleteUser] Failed to encode JSON data")
		w.WriteHeader(http.StatusInternalServerError)
		return &apiError{Error: err, Ret: -4, Info: err.Error()}
	}

	if err = mgodb.DeleteSession(uuid); err != nil {
		str := fmt.Sprintf("Failed to delete session id %s", uuid)
		log.Println("[user.handleLogin] ", str)
		w.WriteHeader(http.StatusInternalServerError)
		return &apiError{Error: err, Ret: -3, Info: str}
	}

	opshdlr.Write("User", d.Name, "delete", "user", "")
	log.Printf("[user.handleDeleteUser] Delete user %s successfully\n", d.Name)

	return nil
}

func handleChangePassword(w http.ResponseWriter, r *http.Request) *apiError {
	log.Printf("[user.handleChangePassword] handleChangePassword Entry %s\n", r.URL)

	// Read request JSON data
	var d UserRequestData
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&d)
	if err != nil {
		str := fmt.Sprintf("failed to parse JSON %v", err)
		log.Println("[user.handleChangePassword] ", str)
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -1, Info: str}
	}

	// Validate data format
	if d.Action != "change_password" {
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: `action should be "delete_user"`}
	} else if valid.Matches(d.Name, `(?i)^[a-z\d.]{5,}$`) == false {
		info := `user name should match the pattern "/^[a-z\d.]{5,}$/i"`
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: info}
	} else if valid.Matches(d.Password, `^.{6,}$`) == false {
		info := `password are case sensitive and must contain a minimum of 6 English characters`
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: info}
	} else if valid.Matches(d.NewPassword, `^.{6,}$`) == false {
		info := `new password is case sensitive and must contain a minimum of 6 English characters`
		w.WriteHeader(http.StatusBadRequest)
		return &apiError{Error: err, Ret: -2, Info: info}
	}

	ae, _ := withAuthentication(w, r, d.Name)
	if ae != nil {
		return ae
	}

	hash := sha1.Sum([]byte(d.Password))
	hash2 := sha1.Sum([]byte(d.NewPassword))
	err = mgodb.ChangePassword(d.Name, hash[:], hash2[:])
	if err != nil {
		str := fmt.Sprintf("Failed to change password for user %s: %s", d.Name, err)
		log.Println("[user.handleChangePassword]", str)
		w.WriteHeader(http.StatusUnauthorized)
		return &apiError{Error: err, Ret: -3, Info: str}
	}

	// Write header before write JSON data
	w.WriteHeader(http.StatusOK)

	var rd ResponseData
	rd.Ret = 0
	rd.Info = "OK"
	if err = json.NewEncoder(w).Encode(rd); err != nil {
		log.Println("[user.handleChangePassword] Failed to encode JSON data")
		w.WriteHeader(http.StatusInternalServerError)
		return &apiError{Error: err, Ret: -4, Info: err.Error()}
	}

	opshdlr.Write("User", d.Name, "change", "password", "")
	log.Printf("[user.handleChangePassword] Change password for user %s successfully\n", d.Name)

	return nil
}
